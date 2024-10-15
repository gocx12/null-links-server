package chat

import (
	"context"

	"null-links/http_service/internal/svc"
	"null-links/http_service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChatWsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatWsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatWsLogic {
	return &ChatWsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChatWsLogic) ChatWs(req *types.ChatWsReq) (resp *types.ChatWsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
