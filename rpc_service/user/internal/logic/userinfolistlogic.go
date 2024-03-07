package logic

import (
	"context"

	"null-links/internal"
	"null-links/rpc_service/user/internal/svc"
	"null-links/rpc_service/user/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserInfoListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoListLogic {
	return &UserInfoListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserInfoListLogic) UserInfoList(in *user.UserInfoListReq) (*user.UserInfoListResp, error) {
	resp := &user.UserInfoListResp{}

	userListDb, err := l.svcCtx.UserModel.FindMulti(l.ctx, in.UserIdList)
	if err != nil {
		logx.Error("get user list from mysql error: ", err)

		resp.StatusCode = internal.StatusRpcErr
		resp.StatusMsg = "fail to get user's info"
		return resp, nil
	}

	userList := make([]*user.UserInfo, len(userListDb))
	for _, userDb := range userListDb {
		userInfo := &user.UserInfo{
			Id: userDb.Id,
			// IsFollow:      userDb.IsFollow != 0,
			Name:          userDb.Username,
			Signature:     userDb.Signature,
			AvatarUrl:     userDb.AvatarUrl,
			BackgroundUrl: userDb.BackgroundUrl,
			FollowCount:   userDb.FollowCount,
			FollowerCount: userDb.FollowerCount,
		}
		userList = append(userList, userInfo)
	}

	resp.StatusCode = internal.StatusSuccess
	resp.StatusMsg = "success"
	resp.UserList = userList
	return resp, nil
}

func (l *UserInfoListLogic) getUserInfo(userId int64) (*user.UserInfo, error) {
	resp, err := l.svcCtx.UserModel.FindOne(l.ctx, userId)

	if err != nil {
		return nil, err
	}

	userInfo := &user.UserInfo{
		Id: resp.Id,
		// IsFollow:      resp.IsFollow != 0,
		Name:          resp.Username,
		Signature:     resp.Signature,
		AvatarUrl:     resp.AvatarUrl,
		BackgroundUrl: resp.BackgroundUrl,
		FollowCount:   resp.FollowCount,
		FollowerCount: resp.FollowerCount,
	}

	return userInfo, nil
}
