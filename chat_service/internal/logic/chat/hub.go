package chat

import (
	"sync"
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
}

var onceHubMap sync.Once
var hubMap map[int64]*Hub
var mutex sync.Mutex

func NewHub(websetId int64) *Hub {
	onceHubMap.Do(func() {
		hubMap = make(map[int64]*Hub)
	})

	mutex.Lock()
	defer mutex.Unlock()
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
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client)
				}
			}
		}
	}
}
