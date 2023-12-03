package logic

import (
	"context"

	"nulltv/rpc_service/user/internal/model"
	"nulltv/rpc_service/user/internal/svc"
	"nulltv/rpc_service/user/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserInfoLogic) UserInfo(in *user.UserInfoReq) (*user.UserInfoResp, error) {
	resp, err := l.svcCtx.UserModel.FindOne(l.ctx, in.UserId)

	if err != nil && err != model.ErrNotFound {
		return nil, err
	}

	if resp == nil {
		return &user.UserInfoResp{
			StatusCode: 1004,
			StatusMsg:  "用户不存在",
			UserInfo:   nil,
		}, nil
	}

	var respUser user.UserInfo

	respUser.Id = resp.Id
	respUser.Name = resp.Username
	respUser.Signature = resp.Signature
	respUser.AvatarUrl = resp.AvatarUrl
	respUser.IsFollow = resp.IsFollow != 0
	respUser.BackgroundUrl = resp.BackgroundUrl
	respUser.FollowCount = resp.FollowCount
	respUser.FollowerCount = resp.FollowerCount

	return &user.UserInfoResp{
		StatusCode: 200,
		StatusMsg:  "成功",
		UserInfo:   &respUser,
	}, nil
}
