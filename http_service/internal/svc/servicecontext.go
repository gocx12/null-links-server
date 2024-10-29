package svc

import (
	"null-links/http_service/internal/config"
	"null-links/http_service/internal/infrastructure/cache"
	"null-links/http_service/internal/infrastructure/model"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config
	// UserRpc   user.UserServiceClient
	// WebsetRpc webset.WebsetServiceClient

	FavoriteModel   model.TFavoriteModel
	ChatModel       model.TChatModel
	LikeModel       model.TLikeModel
	WeblinkModel    model.TWeblinkModel
	WebsetModel     model.TWebsetModel
	UserModel       model.TUserModel
	TopicModel      model.TTopicModel
	AdviceModel     model.TAdviceModel
	BalanceModel    model.TBalanceModel
	RedisClient     *redis.Client
	WlCoverKqPusher *kq.Pusher
	VdEmailKqPusher *kq.Pusher
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		// UserRpc:       user.NewUserServiceClient(zrpc.MustNewClient(c.UserRpcConf).Conn()),
		// WebsetRpc:     webset.NewWebsetServiceClient(zrpc.MustNewClient(c.WebsetRpcConf).Conn()),
		ChatModel:     model.NewTChatModel(sqlx.NewMysql(c.DataSource)),
		UserModel:     model.NewTUserModel(sqlx.NewMysql(c.DataSource)),
		FavoriteModel: model.NewTFavoriteModel(sqlx.NewMysql(c.DataSource)),
		LikeModel:     model.NewTLikeModel(sqlx.NewMysql(c.DataSource)),
		WeblinkModel:  model.NewTWeblinkModel(sqlx.NewMysql(c.DataSource)),
		WebsetModel:   model.NewTWebsetModel(sqlx.NewMysql(c.DataSource)),
		TopicModel:    model.NewTTopicModel(sqlx.NewMysql(c.DataSource)),
		AdviceModel:   model.NewTAdviceModel(sqlx.NewMysql(c.DataSource)),
		BalanceModel:  model.NewTBalanceModel(sqlx.NewMysql(c.DataSource)),

		RedisClient: cache.NewClient(&redis.Options{
			Addr:     c.RedisConf.Host,
			Password: c.RedisConf.Password,
			DB:       c.RedisConf.DB,
		}),

		WlCoverKqPusher: kq.NewPusher(c.WlCoverKqPusherConf.Brokers, c.WlCoverKqPusherConf.Topic),
		VdEmailKqPusher: kq.NewPusher(c.VdEmailKqPusherConf.Brokers, c.VdEmailKqPusherConf.Topic),
	}
}
