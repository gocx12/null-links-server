package webset

import (
	"context"
	"sync"

	"null-links/http_service/internal/common"
	"null-links/http_service/internal/svc"
	"null-links/http_service/internal/types"
	"null-links/internal"
	"null-links/rpc_service/webset/pb/webset"

	// "null-links/rpc_service/user/pb/user"
	// "null-links/rpc_service/webset/pb/webset"

	"github.com/demdxx/gocast"
	"github.com/zeromicro/go-zero/core/logx"
)

type FeedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFeedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FeedLogic {
	return &FeedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FeedLogic) Feed(req *types.FeedReq) (resp *types.FeedResp, err error) {
	if req.PageSize > 50 {
		req.PageSize = 50
	}

	resp = &types.FeedResp{}

	// respRpc, err := l.svcCtx.WebsetRpc.Feed(l.ctx, &webset.FeedReq{
	// 	UserId:   req.UserId,
	// 	Page:     req.Page,
	// 	PageSize: req.PageSize,
	// })
	// if err != nil {
	// 	resp.StatusCode = internal.StatusRpcErr
	// 	resp.StatusMsg = "获取网页单流失败"
	// 	err = nil
	// 	return
	// }

	// TODO(chancyGao): 推荐系统
	// status: 0 未定义，1 发布，2 待审核，3 审核不通过，4 定时发布，5 删除
	websetListDB, err := l.svcCtx.WebsetModel.FindRecent(l.ctx, req.Page, req.PageSize)
	if err != nil {
		logx.Error("get webset list from db error: ", err)
		resp.StatusCode = internal.StatusGatewayErr
		resp.StatusMsg = "get webset list from db error"
		return
	}
	if len(websetListDB) == 0 {
		resp.StatusCode = internal.StatusSuccess
		resp.StatusMsg = "success"
		resp.WebsetList = make([]types.WebsetShort, 0)
		return
	}

	websetIdList := make([]int64, 0, len(websetListDB))
	authorIdList := make([]int64, 0, len(websetListDB))
	for _, webset := range websetListDB {
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
		if req.UserId == -1 {
			// 没有id信息则全部未点赞
			for _, item := range websetListDB {
				websetLikeMap[item.Id] = false
			}
		} else {
			likeInfosDb, err := l.svcCtx.LikeModel.GetLikeWebsetUserInfos(l.ctx, websetIdList, req.UserId)
			if err != nil {
				logx.Error("get like info failed, user, err: ", err)
			}
			for index, item := range websetListDB {
				isFound := false
				for _, likeInfo := range likeInfosDb {

					likeCnt, err := l.svcCtx.RedisClient.HGet(l.ctx, common.RdsKeyWebsetLikedCnt, gocast.ToString(item.Id)).Result()
					if err != nil {
						logx.Error("get like count from redis failed, err: ", err)
					}
					// mysql + redis = 总点赞数
					websetListDB[index].LikeCnt = item.LikeCnt + gocast.ToInt64(likeCnt)

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
		if req.UserId == -1 {
			// 没有id信息则全部未点赞
			for _, item := range websetListDB {
				websetFavoriteMap[item.Id] = false
			}
		} else {
			favoriteInfosDb, err := l.svcCtx.FavoriteModel.GetFavoriteWebsetUserInfos(l.ctx, websetIdList, req.UserId)
			if err != nil {
				logx.Error("get like info failed, user, err: ", err)
			}
			for _, item := range websetListDB {
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
		userInfoListDB, err := l.svcCtx.UserModel.FindMulti(l.ctx, authorIdList)
		if err != nil {
			logx.Error("get user info list failed, err: ", err)
		}
		for _, item := range websetListDB {
			isFound := false
			for _, userInfo := range userInfoListDB {
				if item.AuthorId == userInfo.Id {
					websetAuthorMap[item.AuthorId] = &webset.UserInfoShort{
						Id:            userInfo.Id,
						Name:          userInfo.Username,
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
		wg.Done()
	}()

	wg.Wait()

	// // 打包返回结果
	// websetListResp := make([]*webset.WebsetShort, 0, len(websetListDB))
	// for _, item := range websetListDB {
	// 	websetListResp = append(websetListResp, &webset.WebsetShort{
	// 		Id:            item.Id,
	// 		Title:         item.Title,
	// 		AuthorInfo:    websetAuthorMap[item.AuthorId],
	// 		CoverUrl:      item.CoverUrl,
	// 		LikeCount:     item.LikeCnt,
	// 		IsLike:        websetLikeMap[item.Id],
	// 		FavoriteCount: item.FavoriteCnt,
	// 		IsFavorite:    websetFavoriteMap[item.Id],
	// 	})
	// }

	resp.StatusCode = internal.StatusSuccess
	resp.StatusMsg = "获取网页单流成功"
	resp.WebsetList = make([]types.WebsetShort, 0, len(websetListDB))
	for _, webset := range websetListDB {
		authorInfo := websetAuthorMap[webset.AuthorId]
		resp.WebsetList = append(resp.WebsetList, types.WebsetShort{
			ID:            webset.Id,
			Title:         webset.Title,
			CoverUrl:      webset.CoverUrl,
			IsLike:        websetLikeMap[webset.Id],
			LikeCount:     webset.LikeCnt,
			ViewCount:     webset.ViewCnt,
			FavoriteCount: webset.FavoriteCnt,
			AuthorInfo: types.UserShort{
				Id:        authorInfo.Id,
				Name:      authorInfo.Name,
				AvatarUrl: authorInfo.AvatarUrl,
			},
		})
	}

	return
}
