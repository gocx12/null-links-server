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
	}

	resp.ChatList = make([]types.Chat, 0)
	chatListDb, err := l.svcCtx.ChatModel.FindByTopicId(l.ctx, req.TopicId)

	resp.TopicTitle = topicDb.Title
	resp.StatusCode = internal.StatusSuccess
	resp.StatusMsg = "success"
	return
}
