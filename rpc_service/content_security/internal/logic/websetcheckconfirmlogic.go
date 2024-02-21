package logic

import (
	"context"

	"null-links/rpc_service/content_security/internal/svc"
	"null-links/rpc_service/content_security/pb/content_security"

	"github.com/zeromicro/go-zero/core/logx"
)

type WebsetCheckConfirmLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewWebsetCheckConfirmLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WebsetCheckConfirmLogic {
	return &WebsetCheckConfirmLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *WebsetCheckConfirmLogic) WebsetCheckConfirm(in *content_security.WebsetCheckConfirmReq) (*content_security.WebsetCheckConfirmResp, error) {
	// todo: add your logic here and delete this line

	return &content_security.WebsetCheckConfirmResp{}, nil
}
