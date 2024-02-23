package logic

import (
	"fmt"
	"testing"

	"null-links/rpc_service/webset/pb/webset"
)

func TestFeed(t *testing.T) {
	l := NewFeedLogic(ctx, svcCtx)

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

	for i := range tests {
		t.Run(tests[i].testName, func(t *testing.T) {
			resp, err := l.Feed(&tests[i].args)
			if (err != nil) != tests[i].wantErr {
				t.Errorf("Register() error: %v, wantErr %v", err, tests[i].wantErr)
				return
			}
			fmt.Printf("resp: %v\n", resp)

		})
	}
}
