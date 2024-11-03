package chat

import (
	"context"

	"null-links/http_service/internal/svc"
	"null-links/http_service/internal/types"
	"null-links/internal"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTopicLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTopicLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTopicLogic {
	return &GetTopicLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTopicLogic) GetTopic(req *types.ChatGetTopicReq) (resp *types.ChatGetTopicResp, err error) {
	resp = &types.ChatGetTopicResp{}
	topicDb, err := l.svcCtx.TopicModel.FindOne(l.ctx, req.TopicId)
	if err != nil {
		logx.Error("getTopic error | err=", err)
		resp.StatusCode = internal.StatusErr
		resp.StatusMsg = "get topic error"
		return
	}

	resp.ChatList = make([]types.Chat, 0)
	chatListDb, err := l.svcCtx.ChatModel.FindChatListByTopicId(l.ctx, req.TopicId, 0, 10)
	if err != nil {
		logx.Error("getChatListByTopicId error | err=", err)
		resp.StatusCode = internal.StatusErr
		resp.StatusMsg = "get topic error"
		return
	}

	// user去重
	userIdSet := make(map[int64]bool)
	for _, chat := range chatListDb {
		userIdSet[chat.UserId] = true
	}
	userIdList := make([]int64, 0)
	for userId := range userIdSet {
		userIdList = append(userIdList, userId)
	}

	userInfoDb, err := l.svcCtx.UserModel.FindMulti(l.ctx, userIdList)
	if err != nil {
		logx.Error("findMulti error | err=", err)
		resp.StatusCode = internal.StatusErr
		resp.StatusMsg = "get topic error"
		return
	}

	resp.TopicTitle = topicDb.Title
	resp.ChatList = make([]types.Chat, 0)
	for _, chat := range chatListDb {
		resp.ChatList = append(resp.ChatList, types.Chat{
			ChatID:   chat.Id,
			Content:  chat.Content,
			UserID:   chat.UserId,
			UserName: userInfoDb[chat.UserId].Username,
		})
	}
	resp.StatusCode = internal.StatusSuccess
	resp.StatusMsg = "success"
	return
}
