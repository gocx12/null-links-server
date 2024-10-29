package pay

import (
	"context"

	"null-links/http_service/internal/svc"
	"null-links/http_service/internal/types"

	"github.com/demdxx/gocast"
	"github.com/zeromicro/go-zero/core/logx"
)

type InfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InfoLogic {
	return &InfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InfoLogic) Info(req *types.PayInfoReq) (resp *types.PayInfoResp, err error) {
	userId := gocast.ToInt64(l.ctx.Value("userId"))
	l.svcCtx.BalanceModel.FindOne(l.ctx, userId)

	resp = &types.PayInfoResp{}

	return
}
