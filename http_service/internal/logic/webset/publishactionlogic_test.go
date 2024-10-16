package webset

import (
	"context"
	"flag"
	"null-links/http_service/internal/config"
	"null-links/http_service/internal/svc"
	"null-links/http_service/internal/types"
	"testing"

	"github.com/zeromicro/go-zero/core/conf"
)

var configFile = flag.String("f", "../../../etc/service.yaml", "the config file")
var ctx = context.Background()
var svcCtx *svc.ServiceContext

func TestMain(m *testing.M) {
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c)
	svcCtx = svc.NewServiceContext(c)

	m.Run()
}

func TestPublishAction(t *testing.T) {
	testCases := []struct {
		testName string
		args     types.PublishActionReq
		want     bool
		wantErr  bool
	}{
		{
			testName: "not exsiting name",
			args: types.PublishActionReq{
				ActionType:  0,
				AuthorId:    1,
				Title:       "test",
				Description: "test",
				Category:    0,
				WebLinkList: []types.WebLinkPublish{
					{
						Url:         "https://www.baidu.com",
						Description: "test",
					},
				},
			},
			want:    true,
			wantErr: false,
		},
	}

	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {
			_, err := NewPublishActionLogic(ctx, svcCtx).PublishAction(&test.args)
			if (err != nil) != test.wantErr {
				t.Errorf("PublishAction() error: %v, wantErr %v", err, test.wantErr)
				return
			}
		})
	}
}
