package logic

import (
	"testing"

	"null-links/rpc_service/webset/pb/webset"
)

func TestPublishAction(t *testing.T) {
	l := NewPublishActionLogic(ctx, svcCtx)

	case1 := webset.PublishActionReq{
		UserId:     1,
		ActionType: 1,
		WebsetId:   1,
	}

	publishActionResp, err := l.PublishAction(&case1)

}
