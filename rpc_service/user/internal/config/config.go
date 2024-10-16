package config

import (
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	DataSource string
	RedisConf  struct {
		Host     string
		Password string
		DB       int
	}
	VdEmailKqPusherConf struct {
		Brokers []string
		Topic   string
	}
}
