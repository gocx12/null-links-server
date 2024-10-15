package chat

import (
	"context"

	"null-links/http_service/internal/svc"
	"null-links/http_service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllTopicLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAllTopicLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllTopicLogic {
	return &GetAllTopicLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAllTopicLogic) GetAllTopic(req *types.ChatGetAllTopicReq) (resp *types.ChatGetAllTopicResp, err error) {
	// todo: add your logic here and delete this line

	return
}
