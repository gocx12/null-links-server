package logic

import (
	"context"

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

	WebsetListRpcResp := make([]*webset.Webset, 0, len(publishListDb))
	for _, item := range publishListDb {
		WebsetListRpcResp = append(WebsetListRpcResp, &webset.Webset{
			Id:       item.Id,
			Title:    item.Title,
			Describe: item.Describe,
			AuthorId: item.AuthorId,
			CoverUrl: item.CoverUrl,
			Category: item.Category,
			ViewCnt:  item.ViewCnt,
			LikeCnt:  item.LikeCnt,
			Status:   item.Status,
		})
	}
	return &webset.PublishListResp{
		StatusCode: 1,
		StatusMsg:  "success",
		WebsetList: WebsetListRpcResp,
	}, nil
}
