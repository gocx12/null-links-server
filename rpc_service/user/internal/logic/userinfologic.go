package logic

import (
	"context"

	"null-links/internal"
	"null-links/rpc_service/user/internal/model"
	"null-links/rpc_service/user/internal/svc"
	"null-links/rpc_service/user/pb/user"

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
	resp := &user.UserInfoResp{}

	userInfoDb, err := l.svcCtx.UserModel.FindOne(l.ctx, in.UserId)
	if err != nil && err != model.ErrNotFound {
		return nil, err
	} else if err == model.ErrNotFound {

		return &user.UserInfoResp{
			StatusCode: 1004,
			StatusMsg:  "用户不存在",
			UserInfo:   nil,
		}, nil
	}

	respUserInfo := user.UserInfo{
		Id:            userInfoDb.Id,
		Name:          userInfoDb.Username,
		Signature:     userInfoDb.Signature,
		AvatarUrl:     userInfoDb.AvatarUrl,
		BackgroundUrl: userInfoDb.BackgroundUrl,
		FollowCount:   userInfoDb.FollowCount,
		FollowerCount: userInfoDb.FollowerCount,
	}
	// TODO(chancyGao): waiting for relation module
	// respUser.IsFollow = resp.IsFollow != 0

	resp.StatusCode = internal.StatusSuccess
	resp.StatusMsg = "success"
	resp.UserInfo = &respUserInfo

	return resp, nil
}
