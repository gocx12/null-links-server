// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chat

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"github.com/demdxx/gocast"
	"github.com/zeromicro/go-zero/core/logx"

	"null-links/chat_service/internal/model"
	"null-links/chat_service/internal/svc"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type ChatWriteMsg struct {
	UserId  int64  `json:"user_id"`
	Token   string `json:"token"`
	Content string `json:"content"`
}

type ChatSendMsg struct {
	UserId  int64  `json:"user_id"`
	Content string `json:"content"`
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	Hub *Hub

	// The websocket connection.
	Conn *websocket.Conn

	// Buffered channel of outbound messages.
	Send chan []byte

	WebsetId int64
	Ctx      context.Context
	SvcCtx   *svc.ServiceContext
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logx.Error("error: ", err)
			}
			break
		}
		logx.Debug("recv message: ", string(message))
		chatWriteMsg := ChatWriteMsg{}
		json.Unmarshal(message, &chatWriteMsg)
		logx.Debug("chat msg: ", chatWriteMsg)

		// generate chat msg id
		chatMsgId := c.genChatMsgId(chatWriteMsg.UserId, c.WebsetId)

		// save to db
		c.SvcCtx.ChatModel.Insert(c.Ctx, &model.TChat{
			ChatId:   chatMsgId,
			UserId:   chatWriteMsg.UserId,
			WebsetId: c.WebsetId,
			Content:  chatWriteMsg.Content,
			Status:   1, // 1 for online
		})

		// broadcast to all clients
		ChatSendMsg := ChatSendMsg{
			UserId:  chatWriteMsg.UserId,
			Content: chatWriteMsg.Content,
		}
		chatSendMsgByte, err := json.Marshal(ChatSendMsg)
		if err != nil {
			logx.Error("json marshal the chatSendMsg failed, err: ", err)
			continue
		}
		message = bytes.TrimSpace(bytes.Replace(chatSendMsgByte, newline, space, -1))
		c.Hub.broadcast <- message
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			logx.Debug("send message: ", string(message))

			// Add queued chat messages to the current websocket message.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) genChatMsgId(userId, websetId int64) string {
	// current time
	curTimeStr := time.Now().Format("20060102150405")
	userIdStr := gocast.ToString(userId)
	websetIdStr := gocast.ToString(websetId)
	chatMsgId := curTimeStr + userIdStr + websetIdStr
	return chatMsgId
}
