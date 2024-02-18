package logic

import (
	"testing"

	"null-links/rpc_service/webset/pb/webset"

	"null-links/rpc_service/webset/internal/svc"
)

func TestFeed(t *testing.T) {

	l := NewFeedLogic(ctx, svcCtx)

	case1 := webset.FeedReq{
		UserId:   1,
		Page:     1,
		PageSize: 10,
	}

	FeedResp, err := l.Feed(&case1)
}
