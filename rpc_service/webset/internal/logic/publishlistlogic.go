package logic

import (
	"context"

	"null-links/rpc_service/user/pb/user"
	"null-links/rpc_service/webset/internal/svc"
	"null-links/rpc_service/webset/pb/webset"

	"github.com/zeromicro/go-zero/core/logx"
	"null-links/internal"
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
		logx.Error("get publish list from db error: ", err, " ,userId: ", in.UserId)
		return &webset.PublishListResp{
			StatusCode: internal.StatusRpcErr,
			StatusMsg:  "get publish list failed from db error",
		}, nil
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
			StatusCode: internal.StatusRpcErr,
			StatusMsg:  "failed",
		}, nil
	}
	mapUserIdUserInfo := make(map[int64]*user.UserInfo)
	for _, item := range userInfoListResp.UserList {
		mapUserIdUserInfo[item.Id] = item
	}

	WebsetListRpcResp := make([]*webset.WebsetShort, 0, len(publishListDb))
	// TODO(chancy): 增加在线状态
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
		StatusCode: internal.StatusSuccess,
		StatusMsg:  "success",
		WebsetList: WebsetListRpcResp,
	}, nil
}
