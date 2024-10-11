package user

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

func TestRegister(t *testing.T) {
	type args struct {
		username string
		password string
	}

	testCases := []struct {
		testName string
		args     args
		want     bool
		wantErr  bool
	}{
		{
			testName: "not exsiting name",
			args: args{
				username: "Alex",
				password: "123456",
			},
			want:    true,
			wantErr: false,
		},
		{
			testName: "existing name",
			args: args{
				username: "gao",
				password: "123456",
			},
			want:    false,
			wantErr: false,
		},
	}

	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {
			resp, err := NewRegisterLogic(ctx, svcCtx).Register(&types.RegisterReq{
				Username: test.args.username,
				Password: test.args.password,
			})
			if (err != nil) != test.wantErr {
				t.Errorf("Register() error: %v, wantErr %v", err, test.wantErr)
				return
			}
			if (resp.UserId != -1) != test.want {
				t.Errorf("Register() error: want %v", test.want)
				return
			}
		})
	}
}
