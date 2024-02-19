package chat

import (
	"context"

	"null-links/chat_service/internal/svc"
	"null-links/chat_service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HistoryLogic {
	return &HistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HistoryLogic) History(req *types.ChatHistoryReq) (resp *types.ChatHistoryResp, err error) {
	// payload, err := internal.ParseJwtToken(l.svcCtx.Config.Auth.AccessSecret, req.Token)
	// if err != nil {
	// 	logx.Error("parse jwt token failed, err: ", err)
	// 	return
	// }
	// userId, ok := payload["user_id"]
	// if !ok {
	// 	logx.Error("parse jwt token failed, user_id not found")
	// 	return
	// }
	// userIdInt := gocast.ToInt64(userId)

	// get chat history

	return
}
