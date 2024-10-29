// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.2

package handler

import (
	"net/http"

	chat "null-links/http_service/internal/handler/chat"
	common "null-links/http_service/internal/handler/common"
	pay "null-links/http_service/internal/handler/pay"
	social "null-links/http_service/internal/handler/social"
	user "null-links/http_service/internal/handler/user"
	webset "null-links/http_service/internal/handler/webset"
	"null-links/http_service/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/ping",
				Handler: PingHandler(serverCtx),
			},
		},
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/get_all_topic",
				Handler: chat.GetAllTopicHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/history",
				Handler: chat.HistoryHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/ws",
				Handler: chat.ChatWsHandler(serverCtx),
			},
		},
		rest.WithPrefix("/chat"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/like",
				Handler: chat.LikeHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/link",
				Handler: chat.LinkHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
		rest.WithPrefix("/chat"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/advice",
				Handler: common.AdviceHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/report",
				Handler: common.ReportHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/upload/file",
				Handler: common.UploadFileHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/upload/pic",
				Handler: common.UploadPicHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
		rest.WithPrefix("/common"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/info",
				Handler: pay.InfoHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/withdraw",
				Handler: pay.WithdrawHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
		rest.WithPrefix("/pay"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/relation/action",
				Handler: social.RelationActionHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/relation/follow/list",
				Handler: social.RelationFollowListHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/relation/follower/list",
				Handler: social.RelationFollowerListHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/relation/friend/list",
				Handler: social.RelationFriendListHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
		rest.WithPrefix("/social"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/check_username",
				Handler: user.CheckUsernameHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/get_validation_code",
				Handler: user.GetValidationCodeHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/login",
				Handler: user.LoginHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/register",
				Handler: user.RegisterHandler(serverCtx),
			},
		},
		rest.WithPrefix("/user"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/info",
				Handler: user.UserInfoHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/modify",
				Handler: user.ModifyHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
		rest.WithPrefix("/user"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/feed",
				Handler: webset.FeedHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/info",
				Handler: webset.WebsetInfoHandler(serverCtx),
			},
		},
		rest.WithPrefix("/webset"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/favorite/action",
				Handler: webset.FavoriteActionHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/favorite/list",
				Handler: webset.FavoriteListHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/like/action",
				Handler: webset.LikeActionHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/publish/action",
				Handler: webset.PublishActionHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/publish/list",
				Handler: webset.PublishListHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
		rest.WithPrefix("/webset"),
	)
}
