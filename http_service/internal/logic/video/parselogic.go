package video

import (
	"context"

	"nulltv/http_service/internal/svc"
	"nulltv/http_service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ParseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewParseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ParseLogic {
	return &ParseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ParseLogic) Parse(req *types.ParseReq) (resp *types.ParseResp, err error) {
	// todo: add your logic here and delete this line

	return
}
