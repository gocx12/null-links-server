package common

import (
	"context"

	"null-links/http_service/internal/svc"
	"null-links/http_service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdviceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdviceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdviceLogic {
	return &AdviceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdviceLogic) Advice(req *types.AdviceReq) (resp *types.AdviceResp, err error) {
	// todo: add your logic here and delete this line

	return
}
