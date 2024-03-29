// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	chat "null-links/chat_service/internal/handler/chat"
	"null-links/chat_service/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/ws",
				Handler: chat.ChatWsHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/history",
				Handler: chat.HistoryHandler(serverCtx),
			},
		},
		rest.WithPrefix("/chat"),
	)
}
