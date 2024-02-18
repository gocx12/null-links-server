package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	DataSource             string
	ValidationKqPusherConf struct {
		Brokers []string
		Topic   string
	}
}
