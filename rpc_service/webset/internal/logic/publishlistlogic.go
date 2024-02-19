package logic

import (
	"context"

	"null-links/rpc_service/user/pb/user"
	"null-links/rpc_service/webset/internal/svc"
	"null-links/rpc_service/webset/pb/webset"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublishListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishListLogic {
	return &PublishListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PublishListLogic) PublishList(in *webset.PublishListReq) (*webset.PublishListResp, error) {
	publishListDb, err := l.svcCtx.WebsetModel.FindPublishList(l.ctx, in.UserId, in.Page, in.PageSize)
	if err != nil {
		logx.Error("find publish list failed, err: ", err, " userId: ", in.UserId)
		return &webset.PublishListResp{
			StatusCode: 0,
			StatusMsg:  "failed",
		}, err
	}

	UserIdList := make([]int64, 0, len(publishListDb))
	for _, item := range publishListDb {
		UserIdList = append(UserIdList, item.AuthorId)
	}

	// 获取用户信息
	userInfoListResp, err := l.svcCtx.UserRpc.UserInfoList(l.ctx, &user.UserInfoListReq{
		UserIdList: UserIdList,
	})
	if err != nil {
		logx.Error("get user info failed, err: ", err)
		return &webset.PublishListResp{
			StatusCode: 0,
			StatusMsg:  "failed",
		}, err
	}
	mapUserIdUserInfo := make(map[int64]*user.UserInfo)
	for _, item := range userInfoListResp.UserList {
		mapUserIdUserInfo[item.Id] = item
	}

	WebsetListRpcResp := make([]*webset.WebsetShort, 0, len(publishListDb))
	for _, item := range publishListDb {
		WebsetListRpcResp = append(WebsetListRpcResp, &webset.WebsetShort{
			Id:            item.Id,
			Title:         item.Title,
			CoverUrl:      item.CoverUrl,
			ViewCount:     item.ViewCnt,
			LikeCount:     item.LikeCnt,
			IsLike:        false,
			FavoriteCount: item.FavoriteCnt,
			IsFavorite:    false,
			AuthorInfo: &webset.UserInfoShort{
				Id:            item.AuthorId,
				Name:          mapUserIdUserInfo[item.AuthorId].Name,
				AvatarUrl:     mapUserIdUserInfo[item.AuthorId].AvatarUrl,
				FollowCount:   mapUserIdUserInfo[item.AuthorId].FollowCount,
				FollowerCount: mapUserIdUserInfo[item.AuthorId].FollowerCount,
				IsFollow:      false, // TODO(chancyGao):从relation系统获取
			},
		})
	}
	return &webset.PublishListResp{
		StatusCode: 1,
		StatusMsg:  "success",
		WebsetList: WebsetListRpcResp,
	}, nil
}
