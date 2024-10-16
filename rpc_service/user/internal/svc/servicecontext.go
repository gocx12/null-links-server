package svc

import (
	"null-links/rpc_service/user/internal/config"
	"null-links/rpc_service/user/internal/model"
	"null-links/rpc_service/user/internal/nl_redis"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config          config.Config
	UserModel       model.TUserModel
	RedisClient     *redis.Client
	VdEmailKqPusher *kq.Pusher
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		UserModel: model.NewTUserModel(sqlx.NewMysql(c.DataSource)),
		RedisClient: nl_redis.NewClient(&redis.Options{
			Addr:     c.RedisConf.Host,
			Password: c.RedisConf.Password,
			DB:       c.RedisConf.DB,
		}),
		VdEmailKqPusher: kq.NewPusher(c.VdEmailKqPusherConf.Brokers, c.VdEmailKqPusherConf.Topic),
	}
}
