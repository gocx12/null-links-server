package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
	DataSource string
	RedisConf  struct {
		Host     string
		Password string
		DB       int
	}
	UserRpcConf zrpc.RpcClientConf
}
