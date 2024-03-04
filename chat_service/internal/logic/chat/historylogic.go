package chat

import (
	"context"

	"null-links/chat_service/internal/svc"
	"null-links/chat_service/internal/types"
	"null-links/internal"

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
	// get chat history
	resp = &types.ChatHistoryResp{}
	switch req.Type {
	case 1:
		chatDb, err := l.svcCtx.ChatModel.FindChatList(l.ctx, req.WebsetID, req.LastChatId, req.Page, req.PageSize)
		if err != nil {
			logx.Error("get chat history from mysql failed, error:", err)
			resp.StatusCode = internal.StatusGatewayErr
			resp.StatusMsg = "获取聊天记录失败"
		}
		// 上划加载历史消息
		resp.StatusCode = internal.StatusSuccess
		resp.StatusMsg = "成功"

		chatList := make([]types.Chat, 0, len(chatDb))
		for _, chat := range chatDb {
			chatList = append(chatList, types.Chat{
				UserID:    chat.UserId,
				Content:   chat.Content,
				CreatedAt: chat.CreatedAt.Format("2006-01-02 15:04"),
			})
		}

		resp.ChatList = chatList
	case 2:
		// 查询历史消息
	default:
		logx.Error("invalid history type")
	}

	return
}
