package logic

import (
	"context"

	"null-links/rpc_service/user/internal/model"
	"null-links/rpc_service/user/internal/svc"
	"null-links/rpc_service/user/pb/user"

	"null-links/internal"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckUsernameLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckUsernameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckUsernameLogic {
	return &CheckUsernameLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CheckUsernameLogic) CheckUsername(in *user.CheckUsernameReq) (*user.CheckUsernameResp, error) {
	resp := &user.CheckUsernameResp{}

	// result 0 用户名不存在，1 用户名已存在
	_, err := l.svcCtx.UserModel.FindOneByName(l.ctx, in.Username)

	switch err {
	case nil:
		resp.StatusCode = internal.StatusSuccess
		resp.StatusMsg = "success"
		resp.Result = 1
		return resp, nil
	case model.ErrNotFound:
		logx.Debug("username: ", in.Username)
		resp.StatusCode = internal.StatusSuccess
		resp.StatusMsg = "success"
		resp.Result = 0
		return resp, nil
	default:
		logx.Error("get user info from db error: ", err)
		resp.StatusCode = internal.StatusRpcErr
		resp.StatusMsg = "get user info from db error"
		resp.Result = -1
		return resp, nil
	}

}
