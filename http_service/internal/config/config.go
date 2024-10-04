package config

import (
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	DataSource string
	Auth       struct {
		AccessSecret string
		AccessExpire int64
	}
	// UserRpcConf   zrpc.RpcClientConf
	// WebsetRpcConf zrpc.RpcClientConf

	RedisConf struct {
		Host     string
		Password string
		DB       int
	}
	MinIO struct {
		Endpoint        string
		AccessKeyID     string
		SecretAccessKey string
		UseSSL          bool
		DownloadHost    string
	}
	VdEmailMqPusherConf struct {
		Brokers []string
		Topic   string
	}
}
