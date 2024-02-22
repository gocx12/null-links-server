package logic

import (
	"fmt"
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
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Printf("publishActionResp: %v\n", publishActionResp)

}
