package user

import (
	"context"

	"null-links/http_service/internal/svc"
	"null-links/http_service/internal/types"
	"null-links/internal"

	"github.com/zeromicro/go-zero/core/logx"
)

type ModifyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewModifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ModifyLogic {
	return &ModifyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ModifyLogic) Modify(req *types.ModifyReq) (resp *types.ModifyResp, err error) {
	// todo: add your logic here and delete this line

	err = l.svcCtx.UserModel.UpdateAvatarUrl(l.ctx, req.UserId, req.AvatarUrl)
	if err != nil {
		resp = &types.ModifyResp{
			StatusCode: internal.StatusRpcErr,
			StatusMsg:  "修改用户信息失败",
		}
		err = nil
	}
	if resp.StatusCode != internal.StatusSuccess {
		resp = &types.ModifyResp{
			StatusCode: internal.StatusRpcErr,
			StatusMsg:  "修改用户信息失败",
		}
		return
	}

	resp = &types.ModifyResp{
		StatusCode: resp.StatusCode,
		StatusMsg:  "修改用户信息成功",
	}

	return
}
