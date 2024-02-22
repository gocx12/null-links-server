package user

import (
	"context"

	"null-links/http_service/internal/svc"
	"null-links/http_service/internal/types"
	"null-links/internal"
	"null-links/rpc_service/user/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo(req *types.UserInfoReq) (resp *types.UserInfoResp, err error) {
	respRpc, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &user.UserInfoReq{
		UserId: req.UserID,
	})
	if err != nil {
		logx.Error("call UserRpc failed, err: ", err)
		resp = &types.UserInfoResp{
			StatusCode: internal.StatusRpcErr,
			StatusMsg:  "获取信息失败",
			User:       types.User{},
		}
		err = nil
		return
	} else if respRpc.StatusCode != internal.StatusSuccess {
		logx.Error("call UserRpc failed, err msg: ", respRpc.StatusMsg)
		resp = &types.UserInfoResp{
			StatusCode: internal.StatusRpcErr,
			StatusMsg:  "获取信息失败",
			User:       types.User{},
		}
		err = nil
		return
	}

	resp = &types.UserInfoResp{
		StatusCode: internal.StatusSuccess,
		StatusMsg:  "获取信息成功",
		User: types.User{
			Id:            respRpc.UserInfo.Id,
			Name:          respRpc.UserInfo.Name,
			Email:         respRpc.UserInfo.Email,
			AvatarUrl:     respRpc.UserInfo.AvatarUrl,
			BackgroundUrl: respRpc.UserInfo.BackgroundUrl,
			FollowCount:   respRpc.UserInfo.FollowCount,
			FollowerCount: respRpc.UserInfo.FollowerCount,
			IsFollow:      respRpc.UserInfo.IsFollow,
			Signature:     respRpc.UserInfo.Signature,
			WorkCount:     respRpc.UserInfo.WorkCount,
		},
	}

	return
}
