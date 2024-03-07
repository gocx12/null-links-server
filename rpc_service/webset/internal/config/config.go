package config

import (
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	DataSource  string
	UserRpcConf zrpc.RpcClientConf
	RedisConf   struct {
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
	WlCoverKqConsumser struct {
		Brokers []string
		Topic   string
	}
}
