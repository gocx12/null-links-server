package chat

import (
	"context"

	"null-links/chat_service/internal/svc"
	"null-links/chat_service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChatHistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatHistoryLogic {
	return &ChatHistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChatHistoryLogic) ChatHistory(req *types.ChatHistoryReq) (resp *types.ChatHistoryResp, err error) {
	// todo: add your logic here and delete this line

	return
}
