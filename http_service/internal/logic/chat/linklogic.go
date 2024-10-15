package chat

import (
	"context"

	"null-links/http_service/internal/svc"
	"null-links/http_service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LinkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLinkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LinkLogic {
	return &LinkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LinkLogic) Link(req *types.ChatLinkReq) (resp *types.ChatLinkResp, err error) {
	// todo: add your logic here and delete this line

	return
}
