package user

import (
	"context"

	"null-links/http_service/internal/svc"
	"null-links/http_service/internal/types"
	"null-links/rpc_service/user/pb/user"

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
	respRpc, err := l.svcCtx.UserRpc.CheckUsername(l.ctx, &user.CheckUsernameReq{
		Username: req.Username,
	})
	if err != nil {
		logx.Error("call UserRpc failed, err: ", err)
		resp = &types.CheckUsernameResp{
			StatusCode: internal.StatusRpcErr,
			StatusMsg:  "检查用户错误",
			Result:     respRpc.Result,
		}
		err = nil
		return
	}

	resp = &types.CheckUsernameResp{
		StatusCode: internal.StatusSuccess,
		StatusMsg:  "检查用户名成功",
		Result:     respRpc.Result,
	}
	return
}
