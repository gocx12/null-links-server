package chat

import (
	"context"
	"flag"
	"fmt"
	"null-links/chat_service/internal/config"
	"null-links/chat_service/internal/svc"
	"testing"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "../../../etc/service.yaml", "the config file")
var svcCtx *svc.ServiceContext

func TestMain(m *testing.M) {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	_ = rest.MustNewServer(c.RestConf)
	svcCtx = svc.NewServiceContext(c)
	m.Run()
}

func TestGenChatId(t *testing.T) {
	client := &Client{
		WebsetId: 1, UserId: 1,
		UserName: "admin",
		Ctx:      context.Background(), SvcCtx: svcCtx,
	}
	_, err := client.genChatMsgId(1, 1)
	if err != nil {
		fmt.Print("err:", err)
	}
	// fmt.Print("chat id: ", chatId)
}
