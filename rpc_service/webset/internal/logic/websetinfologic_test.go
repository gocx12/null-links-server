package logic

import (
	"fmt"
	"testing"

	"null-links/rpc_service/webset/pb/webset"
)

func TestWebsetInfo(t *testing.T) {
	l := NewWebsetInfoLogic(ctx, svcCtx)

	caseList := []*webset.WebsetInfoReq{
		{
			UserId:   1,
			WebsetId: 1,
		},
	}

	for _, in := range caseList {
		websetInfoResp, err := l.WebsetInfo(in)
		if err != nil {
			fmt.Print("error:", err)
		}
		fmt.Printf("websetInfoResp: %v\n", websetInfoResp)
	}

}
