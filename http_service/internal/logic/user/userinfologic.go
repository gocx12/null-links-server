package user

import (
	"context"

	"null-links/cron/model"
	"null-links/http_service/internal/svc"
	"null-links/http_service/internal/types"
	"null-links/internal"

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
	if req.UserID < 0 {
		logx.Error()
		return &types.UserInfoResp{
			StatusCode: internal.StatusParamErr,
			StatusMsg:  "user id is invalid",
		}, nil
	}

	userInfoDb, err := l.svcCtx.UserModel.FindOne(l.ctx, req.UserID)
	if err != nil && err != model.ErrNotFound {
		logx.Error("get user info from db error. err=", err)
		return nil, err
	} else if err == model.ErrNotFound {
		return &types.UserInfoResp{
			StatusCode: internal.StatusGatewayErr,
			StatusMsg:  "the user does not exist",
		}, nil
	}

	// TODO(chancyGao): waiting for relation module
	// respUser.IsFollow = resp.IsFollow != 0

	resp = &types.UserInfoResp{
		StatusCode: internal.StatusSuccess,
		StatusMsg:  "success",
		User: types.User{
			Id:            userInfoDb.Id,
			Name:          userInfoDb.Username,
			Email:         userInfoDb.Email,
			AvatarUrl:     userInfoDb.AvatarUrl,
			BackgroundUrl: userInfoDb.BackgroundUrl,
			FollowCount:   userInfoDb.FollowCount,
			FollowerCount: userInfoDb.FollowerCount,
			IsFollow:      false,
			Signature:     userInfoDb.Signature,
			WorkCount:     0,
		},
	}
	return
}
