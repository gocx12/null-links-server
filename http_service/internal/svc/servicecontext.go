package svc

import (
	"null-links/http_service/internal/config"
	"null-links/rpc_service/user/pb/user"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config  config.Config
	UserRpc user.UserServiceClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		UserRpc: user.NewUserServiceClient(zrpc.MustNewClient(c.UserRpcConf).Conn()),
	}
}
