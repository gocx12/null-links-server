package logic

import (
	"context"
	"sync"

	"null-links/rpc_service/user/pb/user"
	"null-links/rpc_service/webset/internal/svc"
	"null-links/rpc_service/webset/pb/webset"

	"null-links/internal"

	"github.com/demdxx/gocast"
	"github.com/zeromicro/go-zero/core/logx"
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
	if len(websetsDb) == 0 {
		return &webset.FeedResp{
			StatusCode: internal.StatusSuccess,
			StatusMsg:  "success",
			WebsetList: []*webset.WebsetShort{},
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
		for _, item := range websetsDb {
			websetLikeMap[item.Id] = false
		}
		// 获取点赞信息
		if in.UserId != -1 {
			websetIdListFiltered := make([]int64, 0, len(websetIdList))
			for _, websetId := range websetIdList {
				// 检查该webset对应的布隆过滤器是否存在
				res, err := l.svcCtx.RedisClient.Exists(l.ctx, "BF_LIKE_"+gocast.ToString(websetId)).Result()
				if err != nil {
					logx.Error("check whether BF_LIKE_ exists from redis error: ", err)
					websetIdListFiltered = append(websetIdListFiltered, websetId)
				} else if res == 1 {
					// 布隆过滤器存在
					res, err := l.svcCtx.RedisClient.BFExists(l.ctx, "BF_LIKE_"+gocast.ToString(websetId), in.UserId).Result()
					if err != nil {
						logx.Error("get like info from redis bloom filter error: ", err)
						websetIdListFiltered = append(websetIdListFiltered, websetId)
					} else if !res {
						// 用户未点赞
						logx.Debug("user not liked")
						websetLikeMap[websetId] = false
					}
				}
			}
			// 从数据库获取未被布隆过滤器过滤掉的webset的点赞信息
			likeInfosDb, err := l.svcCtx.LikeModel.GetLikeWebsetUserInfos(l.ctx, websetIdListFiltered, in.UserId)
			if err != nil {
				logx.Error("get like info failed, user, err: ", err)
			} else {
				for _, likeInfo := range likeInfosDb {
					websetLikeMap[likeInfo.WebsetId] = (likeInfo.Status == 1)
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
		if userInfoListRpcResp.StatusCode == internal.StatusSuccess {
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
					websetAuthorMap[item.AuthorId] = &webset.UserInfoShort{}
				}
			}
		} else {
			logx.Error("get user info from rpc failed, err: ", userInfoListRpcResp.StatusMsg)
			for _, item := range websetsDb {
				websetAuthorMap[item.AuthorId] = &webset.UserInfoShort{}
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
