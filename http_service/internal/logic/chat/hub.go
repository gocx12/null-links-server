package chat

import (
	"encoding/json"
	"sync"

	"github.com/zeromicro/go-zero/core/logx"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	ViewingCnt uint32
}

var onceHubMap sync.Once
var hubMap map[int64]*Hub
var hubMapMutex sync.Mutex

func NewHub(websetId int64) *Hub {
	hubMapMutex.Lock()
	defer hubMapMutex.Unlock()

	onceHubMap.Do(func() {
		hubMap = make(map[int64]*Hub)
	})

	hub, ok := hubMap[websetId]
	if !ok {
		hub = &Hub{
			broadcast:  make(chan []byte),
			Register:   make(chan *Client),
			unregister: make(chan *Client),
			clients:    make(map[*Client]bool),
		}
		hubMap[websetId] = hub
		go hub.run()
	}
	return hub
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.Register:
			logx.Debug("register client, webset_id=", client.WebsetId, ", user_id=", client.UserId)
			h.clients[client] = true
			h.ViewingCnt++
			message := genViewCntChangeMsg(h.ViewingCnt)
			h.broadcastMsg(message)
		case client := <-h.unregister:
			logx.Debug("unregister client, webset_id=", client.WebsetId, ", user_id=", client.UserId)
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
				h.ViewingCnt--
				message := genViewCntChangeMsg(h.ViewingCnt)
				h.broadcastMsg(message)
			}
		case message := <-h.broadcast:
			h.broadcastMsg(message)
		}
	}
}

func (h *Hub) broadcastMsg(message []byte) {
	for client := range h.clients {
		select {
		case client.Send <- message:
		default:
			close(client.Send)
			delete(h.clients, client)
			h.ViewingCnt--
		}
	}
}

func genViewCntChangeMsg(viewingCnt uint32) []byte {
	viewCntChangeMsg := InitSendMsg{
		ViewingCnt: viewingCnt,
	}
	viewCntChangeMsgByte, err := json.Marshal(viewCntChangeMsg)
	if err != nil {
		logx.Error("json marshal the chatSendMsg failed, err: ", err, " viewCntChangeMsg: ", viewCntChangeMsg)
		return nil
	}
	return viewCntChangeMsgByte
}
