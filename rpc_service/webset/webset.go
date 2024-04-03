package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"null-links/rpc_service/webset/internal/config"
	"null-links/rpc_service/webset/internal/server"
	"null-links/rpc_service/webset/internal/svc"
	"null-links/rpc_service/webset/pb/webset"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "rpc_service/webset/etc/webset.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		webset.RegisterWebsetServiceServer(grpcServer, server.NewWebsetServiceServer(ctx))

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
	addr := "0.0.0.0:6062"
	fmt.Printf("Start pprof server, listen addr %s\n", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func (pprofServer) Stop() {
	fmt.Printf("Stop pprof server\n")
}
