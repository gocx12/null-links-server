package logic

import (
	"context"
	"flag"
	"fmt"
	"testing"

	"null-links/rpc_service/webset/internal/config"

	"null-links/rpc_service/webset/pb/webset"

	"github.com/zeromicro/go-zero/core/conf"

	"null-links/rpc_service/webset/internal/svc"
)

var (
	configFile = flag.String("f", "etc/webset.yaml", "the config file")
	svcCtx     *svc.ServiceContext
	ctx        context.Context
)

func TestMain(m *testing.M) {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	svcCtx = svc.NewServiceContext(c)
	ctx = context.Background()

	m.Run()
}

func TestLikeAction(t *testing.T) {
	l := NewLikeActionLogic(ctx, svcCtx)

	case1 := webset.LikeActionReq{
		UserId:     1,
		ActionType: 1,
		WebsetId:   1,
	}
	case2 := webset.LikeActionReq{
		UserId:     1,
		ActionType: 2,
		WebsetId:   1,
	}

	likeActionResp, err := l.LikeAction(&case1)
	if err != nil {
		fmt.Printf("error:", err)
	}
	fmt.Printf("likeActionResp: ", likeActionResp)

	likeActionResp, err = l.LikeAction(&case2)
	if err != nil {
		fmt.Printf("error:", err)
	}
	fmt.Printf("likeActionResp: ", likeActionResp)
}
