package webset

import (
	"context"
	"null-links/http_service/internal/svc"
	"null-links/http_service/internal/types"
	"null-links/internal"
	"null-links/rpc_service/webset/pb/webset"

	"github.com/zeromicro/go-zero/core/logx"
)

type WebsetInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWebsetInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WebsetInfoLogic {
	return &WebsetInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WebsetInfoLogic) WebsetInfo(req *types.WebsetInfoReq) (resp *types.WebsetInfoResp, err error) {
	websetInfoRpcReq, err := l.svcCtx.WebsetRpc.WebsetInfo(l.ctx, &webset.WebsetInfoReq{
		UserId:   req.UserID,
		WebsetId: req.WebsetID,
	})

	if err != nil {
		websetInfoRpcReq.StatusCode = internal.StatusRpcErr
	}

	websetInfoRpcReq.StatusCode = internal.StatusSuccess
	websetInfoRpcReq.StatusMsg = "获取网页单成功"
	return
}
