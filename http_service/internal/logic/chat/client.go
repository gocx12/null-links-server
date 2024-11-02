// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/demdxx/gocast"
	"github.com/zeromicro/go-zero/core/logx"

	"null-links/http_service/internal/infrastructure/model"
	"null-links/http_service/internal/svc"
	"null-links/internal"

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
	UserId       int64  `json:"user_id"`
	Token        string `json:"token"`
	Content      string `json:"content"`
	QuatedChatId int64  `json:"quated_chat_id"`
}

type ChatSendMsg struct {
	UserId    int64  `json:"user_id"`
	UserName  string `json:"user_name"`
	AvatarUrl string `json:"avatar_url"`
	Content   string `json:"content"`
	ChatId    int64  `json:"chat_id"`
	CreatedAt string `json:"created_at"`
	TopicId   int64  `json:"topic_id"`
}

type InitSendMsg struct {
	ViewingCnt uint32 `json:"viewing_cnt"`
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	Hub *Hub

	// The websocket connection.
	Conn *websocket.Conn

	// Buffered channel of outbound messages.
	Send chan []byte

	WebsetId  int64
	UserId    int64
	UserName  string
	AvatarUrl string
	InTime    time.Time

	Ctx    context.Context
	SvcCtx *svc.ServiceContext
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

		// generate chat msg id
		chatMsgId, err := c.genChatMsgId(c.UserId, c.WebsetId)
		if err != nil {
			logx.Error("get chat id key from redis error: ", err)
			// TODO(chancy): 1.如何通知用户发送失败了 2.告警
			continue
		}

		// get topic id
		// topicId, err := c.getTopicId(chatWriteMsg.QuatedChatId)
		// if err != nil {
		// 	logx.Error("get topic id error: ", err)
		// 	continue
		// }
		topicId := int64(-1)

		// save to db
		_, err = c.SvcCtx.ChatModel.Insert(context.Background(), &model.TChat{
			ChatId:   chatMsgId,
			UserId:   c.UserId,
			WebsetId: c.WebsetId,
			Content:  chatWriteMsg.Content,
			TopicId:  topicId,
			Status:   internal.ChatValid.Code(),
		})
		if err != nil {
			logx.Error("insert chat msg to db failed, err: ", err, " chatWriteMsg: ", chatWriteMsg)
			continue
		}

		// TODO(chancyGao): 机审 过关键词库

		// broadcast to all clients
		chatSendMsg := ChatSendMsg{
			UserId:    c.UserId,
			UserName:  c.UserName,
			AvatarUrl: c.AvatarUrl,
			ChatId:    chatMsgId,
			Content:   chatWriteMsg.Content,
			CreatedAt: time.Now().Format("2006-01-02 15:04"),
			TopicId:   c.WebsetId,
		}
		chatSendMsgByte, err := json.Marshal(chatSendMsg)
		if err != nil {
			logx.Error("json marshal the chatSendMsg failed, err: ", err, " chatSendMsg: ", chatSendMsg)
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

var (
	RdsKeyChatIdKeyPrefix = "ChatIdCnt_"
)

func (c *Client) genChatMsgId(userId, websetId int64) (int64, error) {
	// id格式 6位日期 + 5位秒数 + 4位计数
	dateStr := time.Now().Format("060102")

	// 获取当前时间到当天 0 点的秒数（一天有86,400秒）
	now := time.Now()
	zeroTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	duration := now.Sub(zeroTime)
	secondsStr := fmt.Sprintf("%05d", int64(duration.Seconds()))

	chatMsgId := dateStr + secondsStr
	key := RdsKeyChatIdKeyPrefix + gocast.ToString(websetId) + "_" + chatMsgId
	ctx := context.Background()
	count, err := c.SvcCtx.RedisClient.Incr(ctx, key).Result()
	if err != nil {
		logx.Error("get chat id count from redis error:", err)
	}
	logx.Debug("count:", count)
	if count == 1 {
		// 每一秒的第一个，设置一下key过期时间
		c.SvcCtx.RedisClient.Expire(ctx, key, 1*time.Second)
	} else if count >= 10000 {
		return -1, fmt.Errorf("more than 9999 msg during 1 second")
	}
	countStr := fmt.Sprintf("%05d", count)

	chatMsgId += countStr
	logx.Debug("gen chat id:", chatMsgId)

	return gocast.ToInt64(chatMsgId), nil
}

func (c *Client) getTopicId(quotedChatId int64) (int64, error) {
	// 查询引用的消息是否已在某个topic中
	quatedChatDb, err := c.SvcCtx.ChatModel.FindOne(context.Background(), quotedChatId)
	if err != nil {
		return -1, err
	}
	topicId := quatedChatDb.TopicId

	if quatedChatDb.TopicId == -1 {
		quatedChatContent := quatedChatDb.Content[:internal.Min(len(quatedChatDb.Content), 10)]
		resTopic, err := c.SvcCtx.TopicModel.Insert(context.Background(), &model.TTopic{
			Title:  fmt.Sprintf("新话题 %s", quatedChatContent),
			Status: internal.TopicValid.Code(),
		})
		if err != nil {
			return -1, err
		}
		topicId, err = resTopic.LastInsertId()
		if err != nil {
			return -1, err
		}
	}

	return topicId, nil
}
