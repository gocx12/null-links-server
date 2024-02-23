package logic

import (
	"context"
	"sync"

	"null-links/rpc_service/user/pb/user"
	"null-links/rpc_service/webset/internal/svc"
	"null-links/rpc_service/webset/pb/webset"

	"github.com/demdxx/gocast"
	"github.com/zeromicro/go-zero/core/logx"
	"null-links/internal"
)

type FeedLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFeedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FeedLogic {
	return &FeedLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FeedLogic) Feed(in *webset.FeedReq) (*webset.FeedResp, error) {
	// TODO(chancyGao): 推荐系统
	// status: 0 未定义，1 发布，2 待审核，3 审核不通过，4 定时发布，5 删除
	websetsDb, err := l.svcCtx.WebsetModel.FindRecent(l.ctx, in.Page, in.PageSize)
	if err != nil {
		logx.Error("get webset list from db error: ", err)
		return &webset.FeedResp{
			StatusCode: internal.StatusRpcErr,
			StatusMsg:  "get webset list from db error",
		}, nil
	}

	websetIdList := make([]int64, 0, len(websetsDb))
	authorIdList := make([]int64, 0, len(websetsDb))
	for _, webset := range websetsDb {
		websetIdList = append(websetIdList, webset.Id)
		authorIdList = append(authorIdList, webset.AuthorId)
	}

	var wg sync.WaitGroup
	wg.Add(3) // 点赞，收藏，作者信息

	websetLikeMap := make(map[int64]bool)
	websetFavoriteMap := make(map[int64]bool)
	websetAuthorMap := make(map[int64]*webset.UserInfoShort)

	go func() {
		// 获取点赞信息
		if in.UserId == -1 {
			// 没有id信息则全部未点赞
			for _, item := range websetsDb {
				websetLikeMap[item.Id] = false
			}
		} else {
			likeInfosDb, err := l.svcCtx.LikeModel.GetLikeWebsetUserInfos(l.ctx, websetIdList, in.UserId)
			if err != nil {
				logx.Error("get like info failed, user, err: ", err)
			}
			for index, item := range websetsDb {
				isFound := false
				for _, likeInfo := range likeInfosDb {

					likeCnt, err := l.svcCtx.RedisClient.HGet(l.ctx, RdsKeyWebsetLikedCnt, gocast.ToString(item.Id)).Result()
					if err != nil {
						logx.Error("get like count from redis failed, err: ", err)
					}
					// mysql + redis = 总点赞数
					websetsDb[index].LikeCnt = item.LikeCnt + gocast.ToInt64(likeCnt)

					if item.Id == likeInfo.WebsetId {
						websetLikeMap[item.Id] = (likeInfo.Status == 1)
						isFound = true
						break
					}
				}
				if !isFound {
					websetLikeMap[item.Id] = false
				}
			}
		}
		wg.Done()
	}()

	go func() {
		// 获取收藏信息
		if in.UserId == -1 {
			// 没有id信息则全部未点赞
			for _, item := range websetsDb {
				websetFavoriteMap[item.Id] = false
			}
		} else {
			favoriteInfosDb, err := l.svcCtx.FavoriteModel.GetFavoriteWebsetUserInfos(l.ctx, websetIdList, in.UserId)
			if err != nil {
				logx.Error("get like info failed, user, err: ", err)
			}
			for _, item := range websetsDb {
				isFound := false
				for _, favoriteInfo := range favoriteInfosDb {
					if item.Id == favoriteInfo.WebsetId {
						websetFavoriteMap[item.Id] = (favoriteInfo.Status == 1)
						isFound = true
						break
					}
				}
				if !isFound {
					websetFavoriteMap[item.Id] = false
				}
			}
		}
		wg.Done()
	}()

	go func() {
		// 获取作者信息
		userInfoListRpcResp, err := l.svcCtx.UserRpc.UserInfoList(l.ctx, &user.UserInfoListReq{
			UserIdList: authorIdList,
		})
		if err != nil {
			logx.Error("get user info failed, err: ", err)
		}
		if userInfoListRpcResp.StatusCode == 1 {
			for _, item := range websetsDb {
				isFound := false
				for _, userInfo := range userInfoListRpcResp.UserList {
					if item.AuthorId == userInfo.Id {
						websetAuthorMap[item.AuthorId] = &webset.UserInfoShort{
							Id:            userInfo.Id,
							Name:          userInfo.Name,
							AvatarUrl:     userInfo.AvatarUrl,
							FollowCount:   userInfo.FollowCount,
							FollowerCount: userInfo.FollowerCount,
							IsFollow:      false, // TODO(chancyGao):从relation系统获取
						}
						isFound = true
						break
					}
				}
				if !isFound {
					websetAuthorMap[item.AuthorId] = nil
				}
			}
		}
		wg.Done()
	}()

	wg.Wait()

	// 打包返回结果
	websetListResp := make([]*webset.WebsetShort, 0, len(websetsDb))
	for _, item := range websetsDb {
		websetListResp = append(websetListResp, &webset.WebsetShort{
			Id:            item.Id,
			Title:         item.Title,
			AuthorInfo:    websetAuthorMap[item.AuthorId],
			CoverUrl:      item.CoverUrl,
			LikeCount:     item.LikeCnt,
			IsLike:        websetLikeMap[item.Id],
			FavoriteCount: item.FavoriteCnt,
			IsFavorite:    websetFavoriteMap[item.Id],
		})
	}

	return &webset.FeedResp{
		StatusCode: internal.StatusSuccess,
		StatusMsg:  "success",
		WebsetList: websetListResp,
	}, nil
}
