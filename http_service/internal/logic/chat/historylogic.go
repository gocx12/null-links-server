package chat

import (
	"context"

	"null-links/http_service/internal/infrastructure/model"
	"null-links/http_service/internal/svc"
	"null-links/http_service/internal/types"
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
	var chatDb []*model.TChat
	switch req.Type {
	case 1:
		// 获取初始消息
		chatDb, err = l.svcCtx.ChatModel.FindChatList(l.ctx, req.WebsetID, req.Page, req.PageSize)
	case 2:
		// 查询历史消息
		chatDb, err = l.svcCtx.ChatModel.FindChatListChatId(l.ctx, req.WebsetID, req.LastChatId, req.Page, req.PageSize)
	default:
		logx.Error("invalid history type")
	}

	if err != nil {
		logx.Error("get chat history from mysql failed, error:", err)
		resp.StatusCode = internal.StatusGatewayErr
		resp.StatusMsg = "获取聊天记录失败"
	}
	// 上划加载历史消息
	resp.StatusCode = internal.StatusSuccess
	resp.StatusMsg = "成功"

	// 查询用户名称
	// userid 去重
	userIdNameMap := make(map[int64]string)
	for _, chat := range chatDb {
		userIdNameMap[chat.UserId] = ""
	}
	userIdList := make([]int64, 0, len(userIdNameMap))
	for k := range userIdNameMap {
		userIdList = append(userIdList, k)
	}
	UserInfoListDb, err := l.svcCtx.UserModel.FindMulti(l.ctx, userIdList)
	if err != nil {
		logx.Error("get user info from db error: ", err)
	}
	for _, userInfo := range UserInfoListDb {
		userIdNameMap[userInfo.Id] = userInfo.Username
	}

	chatList := make([]types.Chat, 0, len(chatDb))
	for _, chat := range chatDb {
		chatList = append(chatList, types.Chat{
			UserID:    chat.UserId,
			ChatID:    chat.ChatId,
			UserName:  userIdNameMap[chat.UserId],
			Content:   chat.Content,
			CreatedAt: chat.CreatedAt.Format("2006-01-02 15:04"),
		})
	}

	resp.ChatList = chatList
	return
}
