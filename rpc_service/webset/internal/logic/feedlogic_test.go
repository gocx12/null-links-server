package logic

import (
	"fmt"
	"testing"

	"null-links/rpc_service/webset/pb/webset"
)

func TestFeed(t *testing.T) {
	l := NewFeedLogic(ctx, svcCtx)

	type args struct {
		username string
		password string
	}
	tests := []struct {
		testName string
		args     webset.FeedReq
		want     bool
		wantErr  bool
	}{
		{
			testName: "not exsiting name",
			args: webset.FeedReq{
				UserId:   1,
				Page:     1,
				PageSize: 10,
			},
			want:    true,
			wantErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			resp, err := l.Feed(&test.args)
			if (err != nil) != test.wantErr {
				t.Errorf("Register() error: %v, wantErr %v", err, test.wantErr)
				return
			}
			fmt.Printf("resp: %v\n", resp)

		})
	}
}
