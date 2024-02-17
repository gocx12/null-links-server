package user

import (
	"context"

	"null-links/http_service/internal/svc"
	"null-links/http_service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetValidationCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetValidationCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetValidationCodeLogic {
	return &GetValidationCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetValidationCodeLogic) GetValidationCode(req *types.GetValidationCodeReq) (resp *types.GetValidationCodeResp, err error) {
	logx.Debug("email: ", req.Email)
	resp = &types.GetValidationCodeResp{
		StatusCode: 1,
		StatusMsg:  "success",
	}
	err = nil
	return
}
