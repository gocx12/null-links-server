package webset

import (
	"context"
	"sync"

	"null-links/http_service/internal/common"
	"null-links/http_service/internal/infrastructure/model"
	"null-links/http_service/internal/svc"
	"null-links/http_service/internal/types"
	"null-links/internal"

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

	var wg sync.WaitGroup
	wg.Add(3) // 点赞，收藏，作者信息

	websetHasLikeMap := make(map[int64]bool)
	websetFavoriteMap := make(map[int64]bool)
	websetAuthorMap := make(map[int64]types.UserShort)

	go func() {
		l.getLikeInfo(req.UserId, websetListDB, &websetHasLikeMap)
		wg.Done()
	}()

	go func() {
		l.getFavoriteInfo(req.UserId, websetListDB, &websetFavoriteMap)
		wg.Done()
	}()

	go func() {
		l.getAuthorInfo(websetListDB, &websetAuthorMap)
		wg.Done()
	}()

	wg.Wait()

	resp.StatusCode = internal.StatusSuccess
	resp.StatusMsg = "获取网页单流成功"
	resp.WebsetList = make([]types.WebsetShort, 0, len(websetListDB))
	for _, webset := range websetListDB {
		resp.WebsetList = append(resp.WebsetList, types.WebsetShort{
			Id:            webset.Id,
			Title:         webset.Title,
			CoverUrl:      webset.CoverUrl,
			IsLike:        websetHasLikeMap[webset.Id],
			LikeCount:     webset.LikeCnt,
			ViewCount:     webset.ViewCnt,
			FavoriteCount: webset.FavoriteCnt,
			AuthorInfo:    websetAuthorMap[webset.AuthorId],
		})
	}

	return
}

func (l *FeedLogic) getLikeInfo(userId int64, websetList []*model.TWebset, websetHasLikeMap *map[int64]bool) {
	// 获取点赞信息
	if userId == -1 {
		// 没有id信息则全部未点赞
		for _, webset := range websetList {
			(*websetHasLikeMap)[webset.Id] = false
		}
	} else {
		websetIdList := make([]int64, 0, len(websetList))
		for _, webset := range websetList {
			websetIdList = append(websetIdList, webset.Id)
		}
		likeInfosDB, err := l.svcCtx.LikeModel.GetLikeWebsetUserInfos(l.ctx, websetIdList, userId)
		if err != nil {
			logx.Error("get like info failed, user, err: ", err)
		}
		for index, webset := range websetList {
			hasLike := false
			for _, likeInfo := range likeInfosDB {

				// 1. 综合DB和缓存，统计该webset的总点赞数量
				likeCnt, err := l.svcCtx.RedisClient.HGet(l.ctx, common.RdsKeyWebsetLikedCnt, gocast.ToString(webset.Id)).Result()
				if err != nil {
					logx.Error("get like count from redis failed, err: ", err)
				}
				// mysql + redis = 总点赞数
				websetList[index].LikeCnt = webset.LikeCnt + gocast.ToInt64(likeCnt)

				// 2. 当前用户是否有点赞该webset
				if webset.Id == likeInfo.WebsetId {
					hasLike = (likeInfo.Status == 1)
					break
				}
			}
			(*websetHasLikeMap)[webset.Id] = hasLike
		}
	}

}

func (l *FeedLogic) getFavoriteInfo(userId int64, websetList []*model.TWebset, websetFavoriteMap *map[int64]bool) {
	// 获取收藏信息
	if userId == -1 {
		// 没有id信息则全部未点赞
		for _, item := range websetList {
			(*websetFavoriteMap)[item.Id] = false
		}
	} else {
		websetIdList := make([]int64, 0, len(websetList))
		for _, webset := range websetList {
			websetIdList = append(websetIdList, webset.Id)
		}
		favoriteInfosDb, err := l.svcCtx.FavoriteModel.GetFavoriteWebsetUserInfos(l.ctx, websetIdList, userId)
		if err != nil {
			logx.Error("get like info failed, user, err: ", err)
		}
		for _, webset := range websetList {
			isFound := false
			for _, favoriteInfo := range favoriteInfosDb {
				if webset.Id == favoriteInfo.WebsetId {
					(*websetFavoriteMap)[webset.Id] = (favoriteInfo.Status == 1)
					isFound = true
					break
				}
			}
			if !isFound {
				(*websetFavoriteMap)[webset.Id] = false
			}
		}
	}
}

func (l *FeedLogic) getAuthorInfo(websetList []*model.TWebset, websetAuthorMap *map[int64]types.UserShort) {
	authorIdList := make([]int64, 0, len(websetList))
	for _, webset := range websetList {
		authorIdList = append(authorIdList, webset.AuthorId)
	}

	// 获取作者信息
	userInfoListDB, err := l.svcCtx.UserModel.FindMulti(l.ctx, authorIdList)
	if err != nil {
		logx.Error("get user info list failed, err: ", err)
	}
	for _, item := range websetList {
		isFound := false
		for _, userInfo := range userInfoListDB {
			if item.AuthorId == userInfo.Id {
				(*websetAuthorMap)[item.AuthorId] = types.UserShort{
					Id:        userInfo.Id,
					Name:      userInfo.Username,
					AvatarUrl: userInfo.AvatarUrl,
				}
				isFound = true
				break
			}
		}
		if !isFound {
			(*websetAuthorMap)[item.AuthorId] = types.UserShort{}
		}
	}
}
