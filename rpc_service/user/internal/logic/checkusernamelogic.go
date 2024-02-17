package logic

import (
	"context"

	"null-links/rpc_service/user/internal/svc"
	"null-links/rpc_service/user/pb/user"

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
	userInfo, err := l.svcCtx.UserModel.FindOneByName(l.ctx, in.Username)
	if err != nil {
		return &user.CheckUsernameResp{
			StatusCode: 0,
			StatusMsg:  "fail to check username, err: " + err.Error(),
			Result:     1,
		}, err
	}

	var result int32 = 0
	if userInfo != nil {
		result = 1
	}

	logx.Debug("username: ", in.Username, ", check result: ", result)
	return &user.CheckUsernameResp{
		StatusCode: 1,
		StatusMsg:  "success",
		Result:     result,
	}, nil
}
