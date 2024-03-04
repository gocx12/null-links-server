package svc

import (
	"null-links/chat_service/internal/config"
	"null-links/chat_service/internal/model"
	"null-links/chat_service/internal/nl_redis"
	"null-links/rpc_service/user/pb/user"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config      config.Config
	ChatModel   model.TChatModel
	RedisClient *redis.Client
	UserRpc     user.UserServiceClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		ChatModel: model.NewTChatModel(sqlx.NewMysql(c.DataSource)),
		RedisClient: nl_redis.NewClient(&redis.Options{
			Addr:     c.RedisConf.Host,
			Password: c.RedisConf.Password,
			DB:       c.RedisConf.DB,
		}),
		UserRpc: user.NewUserServiceClient(zrpc.MustNewClient(c.UserRpcConf).Conn()),
	}
}
