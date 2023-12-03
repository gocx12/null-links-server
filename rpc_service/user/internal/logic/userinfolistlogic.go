package logic

import (
	"context"

	"nulltv/rpc_service/user/internal/svc"
	"nulltv/rpc_service/user/pb/user"

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
	userListDb, err := l.svcCtx.UserModel.FindMulti(l.ctx, in.UserIdList)
	if err != nil {
		logx.Error("get user list from mysql error: ", err)
		return &user.UserInfoListResp{
			StatusCode: 1004,
			StatusMsg:  "获取用户信息失败",
			UserList:   nil,
		}, err
	}

	userList := make([]*user.UserInfo, len(userListDb))
	for i, userDb := range userListDb {
		userInfo := &user.UserInfo{
			Id:            userDb.Id,
			IsFollow:      userDb.IsFollow != 0,
			Name:          userDb.Username,
			Signature:     userDb.Signature,
			AvatarUrl:     userDb.AvatarUrl,
			BackgroundUrl: userDb.BackgroundUrl,
			FollowCount:   userDb.FollowCount,
			FollowerCount: userDb.FollowerCount,
		}
		userList[i] = userInfo
	}
	return &user.UserInfoListResp{
		StatusCode: 0,
		StatusMsg:  "成功",
		UserList:   userList,
	}, nil

	// DEPRECATED
	// userList := make([]*user.UserInfo, len(in.UserIdList))

	// // When the length of user id list is small, it is not necessary to use goroutine.
	// // The cost of create goroutines is higher than the cost of get user info synchronously
	// if len(in.UserIdList) < 16 {
	// 	for i, userId := range in.UserIdList {
	// 		userInfo, err := l.getUserInfo(userId)
	// 		if err != nil {
	// 			userList[i] = &user.UserInfo{Id: userId}
	// 			logc.Info(l.ctx, "UserInfoList() try to get user %d info error: %v", userId, err)
	// 		}
	// 		userList[i] = userInfo
	// 	}

	// 	return &user.UserInfoListResp{
	// 		StatusCode: 0,
	// 		StatusMsg:  "成功",
	// 		UserList:   userList,
	// 	}, nil
	// }

	// var wg sync.WaitGroup
	// wg.Add(len(in.UserIdList))

	// for i, userId := range in.UserIdList {
	// 	go func(i int, userId int64) {
	// 		userInfo, err := l.getUserInfo(userId)
	// 		if err != nil {
	// 			userList[i] = &user.UserInfo{Id: userId}
	// 			logc.Info(l.ctx, "UserInfoList() try to get user %d info error: %v", userId, err)
	// 		}
	// 		userList[i] = userInfo
	// 		wg.Done()
	// 	}(i, userId)
	// }
	// wg.Wait()

	// return &user.UserInfoListResp{
	// 	StatusCode: 0,
	// 	StatusMsg:  "成功",
	// 	UserList:   userList,
	// }, nil
}

func (l *UserInfoListLogic) getUserInfo(userId int64) (*user.UserInfo, error) {
	resp, err := l.svcCtx.UserModel.FindOne(l.ctx, userId)

	if err != nil {
		return nil, err
	}

	userInfo := &user.UserInfo{
		Id:            resp.Id,
		IsFollow:      resp.IsFollow != 0,
		Name:          resp.Username,
		Signature:     resp.Signature,
		AvatarUrl:     resp.AvatarUrl,
		BackgroundUrl: resp.BackgroundUrl,
		FollowCount:   resp.FollowCount,
		FollowerCount: resp.FollowerCount,
	}

	return userInfo, nil
}
