package webset

import (
	"log"
	"net/http"

	"nulltv/http_service/internal/logic/webset"
	"nulltv/http_service/internal/svc"
	"nulltv/http_service/internal/types"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var upgrader = websocket.Upgrader{} // use default options

func ChatWebSocketHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ChatWebSocketReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("Error during connection upgradation:", err)
			return
		}
		defer conn.Close()

		l := webset.NewChatWebSocketLogic(r.Context(), svcCtx)
		resp, err := l.ChatWebSocket(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
