package svc

import (
	"null-links/http_service/internal/config"
	"null-links/http_service/internal/infrastructure/cache"
	"null-links/http_service/internal/infrastructure/model"
	"null-links/rpc_service/user/pb/user"

	"null-links/rpc_service/webset/pb/webset"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config    config.Config
	UserRpc   user.UserServiceClient
	WebsetRpc webset.WebsetServiceClient

	FavoriteModel      model.TFavoriteModel
	LikeModel          model.TLikeModel
	WeblinkModel       model.TWeblinkModel
	WebsetModel        model.TWebsetModel
	UserModel          model.TUserModel
	RedisClient        *redis.Client
	WlCoverKqConsumser *kq.Pusher

	VdEmailMqPusher *kq.Pusher
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:        c,
		UserRpc:       user.NewUserServiceClient(zrpc.MustNewClient(c.UserRpcConf).Conn()),
		WebsetRpc:     webset.NewWebsetServiceClient(zrpc.MustNewClient(c.WebsetRpcConf).Conn()),
		UserModel:     model.NewTUserModel(sqlx.NewMysql(c.DataSource)),
		FavoriteModel: model.NewTFavoriteModel(sqlx.NewMysql(c.DataSource)),
		LikeModel:     model.NewTLikeModel(sqlx.NewMysql(c.DataSource)),
		WeblinkModel:  model.NewTWeblinkModel(sqlx.NewMysql(c.DataSource)),
		WebsetModel:   model.NewTWebsetModel(sqlx.NewMysql(c.DataSource)),

		RedisClient: cache.NewClient(&redis.Options{
			Addr:     c.RedisConf.Host,
			Password: c.RedisConf.Password,
			DB:       c.RedisConf.DB,
		}),

		VdEmailMqPusher: kq.NewPusher(c.VdEmailMqPusherConf.Brokers, c.VdEmailMqPusherConf.Topic),
	}
}
