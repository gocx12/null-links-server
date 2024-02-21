package logic

import (
	"context"

	"null-links/rpc_service/content_security/internal/svc"
	"null-links/rpc_service/content_security/pb/content_security"

	"github.com/zeromicro/go-zero/core/logx"
)

type WebsetCheckQueryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewWebsetCheckQueryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WebsetCheckQueryLogic {
	return &WebsetCheckQueryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *WebsetCheckQueryLogic) WebsetCheckQuery(in *content_security.WebsetCheckQueryReq) (*content_security.WebsetCheckQueryResp, error) {
	// todo: add your logic here and delete this line

	return &content_security.WebsetCheckQueryResp{}, nil
}
