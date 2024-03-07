package logic

import (
	"context"

	"null-links/rpc_service/user/internal/svc"
	"null-links/rpc_service/user/pb/user"

	"null-links/internal"

	"github.com/zeromicro/go-zero/core/logx"
)

type ModifyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewModifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ModifyLogic {
	return &ModifyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ModifyLogic) Modify(in *user.ModifyReq) (*user.ModifyResp, error) {
	logx.Debug("modify user info request: ", in)
	err := l.svcCtx.UserModel.UpdateAvatarUrl(l.ctx, in.UserId, in.AvatarUrl)
	if err != nil {
		logx.Error("update avatar url error: ", err)
		return &user.ModifyResp{
			StatusCode: internal.StatusRpcErr,
			StatusMsg:  "modify user info failed",
		}, nil
	}
	return &user.ModifyResp{
		StatusCode: internal.StatusSuccess,
		StatusMsg:  "modify user info success",
	}, nil
}
