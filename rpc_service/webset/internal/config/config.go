package config

import (
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	DataSource  string
	UserRpcConf zrpc.RpcClientConf
	RedisConf   redis.Options
}
