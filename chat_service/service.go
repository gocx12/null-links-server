package main

import (
	"flag"
	"fmt"
	"net/http"

	"null-links/chat_service/internal/config"
	"null-links/chat_service/internal/handler"
	"null-links/chat_service/internal/middleware"
	"null-links/chat_service/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "./chat_service/etc/service.yaml", "the config file")

// var configFile = flag.String("f", "./etc/service.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	// defer server.Stop()

	server.Use(middleware.NewCorsMiddleware().Handle)
	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	svcGroup := service.NewServiceGroup()
	defer svcGroup.Stop()
	svcGroup.Add(server)
	if c.Mode == service.DevMode || c.Mode == service.TestMode {
		svcGroup.Add(pprofServer{})
	}

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	// server.Start()
	svcGroup.Start()
}

type pprofServer struct{}

func (pprofServer) Start() {
	addr := "0.0.0.0:6669"
	fmt.Printf("Start pprof server, listen addr %s\n", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		logx.Error(err)
	}
}

func (pprofServer) Stop() {
	fmt.Printf("Stop pprof server\n")
}
