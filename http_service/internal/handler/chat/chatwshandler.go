package chat

import (
	"context"
	"net/http"

	"null-links/http_service/internal/logic/chat"
	"null-links/http_service/internal/svc"
	"null-links/http_service/internal/types"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ChatWsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ChatWsReq
		if err := httpx.Parse(r, &req); err != nil {
			logx.Error("error params, error: ", err)
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		// start a websocket connection
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logx.Error("fail to upgrade a http to websocket, error: ", err)
			return
		}

		hub := chat.NewHub(req.WebsetID)

		// userInfoRespRpc, err := svcCtx.UserRpc.UserInfo(context.Background(), &user.UserInfoReq{
		// 	UserId: req.UserID,
		// })
		userInfoDb, err := svcCtx.UserModel.FindOne(context.Background(), req.UserID)
		if err != nil {
			logx.Error("get user name and avatar url from user rpc error: ", err)
			return
		}

		client := &chat.Client{Hub: hub, Conn: conn, Send: make(chan []byte, 256),
			WebsetId: req.WebsetID, UserId: req.UserID,
			UserName:  userInfoDb.Username,
			AvatarUrl: userInfoDb.AvatarUrl,
			Ctx:       r.Context(), SvcCtx: svcCtx,
		}
		client.Hub.Register <- client

		// Allow collection of memory referenced by the caller by doing all work in
		// new goroutines.
		go client.WritePump()
		go client.ReadPump()

		// l := chat.NewChatWsLogic(r.Context(), svcCtx)
		// resp, err := l.ChatWs(&req)
		// if err != nil {
		// 	httpx.ErrorCtx(r.Context(), w, err)
		// } else {
		// 	httpx.OkJsonCtx(r.Context(), w, resp)
		// }
	}
}
