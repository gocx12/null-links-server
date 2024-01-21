package logic

import (
	"context"

	"nulltv/rpc_service/video/internal/svc"
	"nulltv/rpc_service/video/pb/video"

	"github.com/zeromicro/go-zero/core/logx"
)

type ParseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewParseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ParseLogic {
	return &ParseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ParseLogic) Parse(in *video.ParseReq) (*video.ParseResp, error) {
	// todo: add your logic here and delete this line

	return &video.ParseResp{}, nil
}
