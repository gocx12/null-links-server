package user

import (
	"context"

	"null-links/cron/model"
	"null-links/http_service/internal/svc"
	"null-links/http_service/internal/types"

	"null-links/internal"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckUsernameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckUsernameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckUsernameLogic {
	return &CheckUsernameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckUsernameLogic) CheckUsername(req *types.CheckUsernameReq) (resp *types.CheckUsernameResp, err error) {
	// respRpc, err := l.svcCtx.UserRpc.CheckUsername(l.ctx, &user.CheckUsernameReq{
	// 	Username: req.Username,
	// })
	// if err != nil {
	// 	logx.Error("call UserRpc failed, err: ", err)
	// 	resp = &types.CheckUsernameResp{
	// 		StatusCode: internal.StatusRpcErr,
	// 		StatusMsg:  "检查用户错误",
	// 		Result:     0,
	// 	}
	// 	err = nil
	// 	return
	// }

	// resp = &types.CheckUsernameResp{
	// 	StatusCode: internal.StatusSuccess,
	// 	StatusMsg:  "检查用户名成功",
	// 	Result:     respRpc.Result,
	// }
	// return

	resp = &types.CheckUsernameResp{}
	// result 0 用户名不存在，1 用户名已存在
	_, err = l.svcCtx.UserModel.FindOneByName(l.ctx, req.Username)

	switch err {
	case nil:
		resp.StatusCode = internal.StatusSuccess
		resp.StatusMsg = "success"
		resp.Result = 0
		return resp, nil
	case model.ErrNotFound:
		logx.Debug("username: ", req.Username)
		resp.StatusCode = internal.StatusSuccess
		resp.StatusMsg = "success"
		resp.Result = 1
		return resp, nil
	default:
		logx.Error("get user info from db error: ", err)
		resp.StatusCode = internal.StatusRpcErr
		resp.StatusMsg = "get user info from db error"
		resp.Result = -1
		return resp, nil
	}
}
