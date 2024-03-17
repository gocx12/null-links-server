package main

import (
	"flag"
	"fmt"
	"net/http"

	"null-links/rpc_service/user/internal/config"
	"null-links/rpc_service/user/internal/server"
	"null-links/rpc_service/user/internal/svc"
	"null-links/rpc_service/user/pb/user"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "rpc_service/user/etc/user.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		user.RegisterUserServiceServer(grpcServer, server.NewUserServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	// defer s.Stop()

	svcGroup := service.NewServiceGroup()
	defer svcGroup.Stop()
	svcGroup.Add(s)
	if c.Mode == service.DevMode || c.Mode == service.TestMode {
		svcGroup.Add(pprofServer{})
	}

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	// s.Start()
	svcGroup.Start()
}

type pprofServer struct{}

func (pprofServer) Start() {
	addr := "0.0.0.0:6061"
	fmt.Printf("Start pprof server, listen addr %s\n", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		logx.Error(err)
	}
}

func (pprofServer) Stop() {
	fmt.Printf("Stop pprof server\n")
}
