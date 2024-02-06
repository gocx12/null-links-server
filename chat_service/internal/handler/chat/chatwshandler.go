package chat

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"null-links/chat_service/internal/logic/chat"
	"null-links/chat_service/internal/svc"
	"null-links/chat_service/internal/types"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func ChatWsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ChatWsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// start a websocket connection
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		hub := chat.NewHub(req.WebsetID)

		client := &chat.Client{Hub: hub, Conn: conn, Send: make(chan []byte, 256)}
		client.Hub.Register <- client

		// Allow collection of memory referenced by the caller by doing all work in
		// new goroutines.
		go client.WritePump()
		go client.ReadPump()
	}
}
