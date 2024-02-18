package svc

import (
	"null-links/rpc_service/user/internal/config"
	"null-links/rpc_service/user/internal/model"
	"null-links/rpc_service/user/internal/nl_redis"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config                   config.Config
	UserModel                model.TUserModel
	RedisClient              *redis.Redis
	ValidationKqPusherClient *kq.Pusher
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:                   c,
		UserModel:                model.NewTUserModel(sqlx.NewMysql(c.DataSource)),
		RedisClient:              nl_redis.NewRedisClient(c.Redis.RedisConf),
		ValidationKqPusherClient: kq.NewPusher(c.ValidationKqPusherConf.Brokers, c.ValidationKqPusherConf.Topic),
	}
}
