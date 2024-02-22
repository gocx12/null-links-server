package webset

import (
	"context"

	"null-links/http_service/internal/svc"
	"null-links/http_service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublishListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishListLogic {
	return &PublishListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublishListLogic) PublishList(req *types.PublishListReq) (resp *types.PublishListResp, err error) {
	l.svcCtx.WebsetRpc.PublishList(l.ctx, &types.PublishListReq{
		WebsetId: req.WebsetId,
	})

	return
}
