package svc

import (
	"null-links/rpc_service/user/pb/user"
	"null-links/rpc_service/user/userservice"
	"null-links/rpc_service/webset/internal/config"
	"null-links/rpc_service/webset/internal/model"
	"null-links/rpc_service/webset/internal/nl_redis"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-queue/kq"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config             config.Config
	FavoriteModel      model.TFavoriteModel
	LikeModel          model.TLikeModel
	WeblinkModel       model.TWeblinkModel
	WebsetModel        model.TWebsetModel
	RedisClient        *redis.Client
	WlCoverKqConsumser *kq.Pusher
	UserRpc            userservice.UserService
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:        c,
		FavoriteModel: model.NewTFavoriteModel(sqlx.NewMysql(c.DataSource)),
		LikeModel:     model.NewTLikeModel(sqlx.NewMysql(c.DataSource)),
		WeblinkModel:  model.NewTWeblinkModel(sqlx.NewMysql(c.DataSource)),
		WebsetModel:   model.NewTWebsetModel(sqlx.NewMysql(c.DataSource)),
		RedisClient: nl_redis.NewClient(&redis.Options{
			Addr:     c.RedisConf.Host,
			Password: c.RedisConf.Password,
			DB:       c.RedisConf.DB,
		}),
		WlCoverKqConsumser: kq.NewPusher(c.WlCoverKqConsumser.Brokers, c.WlCoverKqConsumser.Topic),
		UserRpc:            user.NewUserServiceClient(zrpc.MustNewClient(c.UserRpcConf).Conn()),
	}
}
