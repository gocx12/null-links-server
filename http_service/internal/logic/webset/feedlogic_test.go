package webset

import (
	"null-links/http_service/internal/types"
	"testing"
)

func TestFeed(t *testing.T) {
	testCases := []struct {
		testName string
		args     types.FeedReq
		want     bool
		wantErr  bool
	}{
		{
			testName: "normal",
			args: types.FeedReq{
				Page:     1,
				PageSize: 20,
			},
			want:    true,
			wantErr: false,
		},
	}

	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {
			_, err := NewFeedLogic(ctx, svcCtx).Feed(&test.args)
			if (err != nil) != test.wantErr {
				t.Errorf("Feed() error: %v, wantErr %v", err, test.wantErr)
				return
			}
		})
	}
}
