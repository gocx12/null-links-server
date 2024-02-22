package user

import (
	"context"

	"null-links/http_service/internal/svc"
	"null-links/http_service/internal/types"
	"null-links/internal"
	"null-links/rpc_service/user/pb/user"

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

	respRpc, err := l.svcCtx.UserRpc.GetValidtaionCode(l.ctx, &user.GetValidtaionCodeReq{
		Email: req.Email,
	})
	if err != nil {
		logx.Error("call UserRpc failed, err: ", err)
		resp = &types.GetValidationCodeResp{
			StatusCode: internal.StatusRpcErr,
			StatusMsg:  "获取验证码失败",
		}
		err = nil
		return
	} else if respRpc.StatusCode != internal.StatusSuccess {
		logx.Error("call UserRpc failed, err: ", resp.StatusMsg)
		resp = &types.GetValidationCodeResp{
			StatusCode: internal.StatusRpcErr,
			StatusMsg:  "获取验证码失败",
		}
		return
	}

	resp = &types.GetValidationCodeResp{
		StatusCode: internal.StatusSuccess,
		StatusMsg:  "success",
	}
	return
}
