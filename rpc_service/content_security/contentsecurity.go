package main

import (
	"flag"
	"fmt"

	"null-links/rpc_service/content_security/internal/config"
	"null-links/rpc_service/content_security/internal/server"
	"null-links/rpc_service/content_security/internal/svc"
	"null-links/rpc_service/content_security/pb/content_security"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/contentsecurity.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		content_security.RegisterContentSecurityServiceServer(grpcServer, server.NewContentSecurityServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
