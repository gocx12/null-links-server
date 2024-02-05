package webset

import (
	"context"

	"nulltv/http_service/internal/svc"
	"nulltv/http_service/internal/types"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

type ChatWebSocketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatWebSocketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatWebSocketLogic {
	return &ChatWebSocketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChatWebSocketLogic) ChatWebSocket(req *types.ChatWebSocketReq) (resp *types.ChatWebSocketResp, err error) {
	// todo: add your logic here and delete this line

	return
}

// Define our message object
type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan Message)           // broadcast channel

func (l *ChatWebSocketLogic) HandleConnection(conn *websocket.Conn) {
	clients[conn] = true
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			logx.Debug("Error during message reading:", err)
			break
		}
		logx.Debugf("Received: %s", message)
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			logx.Debug("Error during message writing:", err)
			break
		}
	}
}
