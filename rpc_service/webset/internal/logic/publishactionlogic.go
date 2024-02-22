package logic

import (
	"context"

	"null-links/rpc_service/webset/internal/model"
	"null-links/rpc_service/webset/internal/svc"
	"null-links/rpc_service/webset/pb/webset"

	"github.com/zeromicro/go-zero/core/logx"
	"null-links/internal"
)

type PublishActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublishActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishActionLogic {
	return &PublishActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PublishActionLogic) PublishAction(in *webset.PublishActionReq) (*webset.PublishActionResp, error) {
	if in.ActionType == 1 {
		// 发布, status==2 待审核
		websetDb := model.TWebset{
			Title:       in.Webset.Title,
			Describe:    in.Webset.Describe,
			AuthorId:    in.Webset.AuthorInfo.Id,
			CoverUrl:    in.Webset.CoverUrl,
			Category:    0,
			ViewCnt:     0,
			LikeCnt:     0,
			FavoriteCnt: 0,
			Status:      2,
		}
		insRes, err := l.svcCtx.WebsetModel.Insert(l.ctx, &websetDb)
		if err != nil {
			logx.Error("insert webset failed, err: ", err)
			return &webset.PublishActionResp{
				StatusCode: internal.StatusRpcErr,
				StatusMsg:  "fail",
			}, nil
		}
		rowsAffected, err := insRes.RowsAffected()
		if err != nil {
			logx.Error("insert webset failed, err: ", err)
			return &webset.PublishActionResp{
				StatusCode: internal.StatusRpcErr,
				StatusMsg:  "fail",
			}, nil
		}
		if rowsAffected == 0 {
			logx.Error("insert webset failed, rows affected: ", rowsAffected)
			return &webset.PublishActionResp{
				StatusCode: internal.StatusRpcErr,
				StatusMsg:  "fail",
			}, nil
		}

		return &webset.PublishActionResp{
			StatusCode: internal.StatusSuccess,
			StatusMsg:  "success",
		}, nil
	} else if in.ActionType == 2 {
		// 更新, status==2 待审核
		l.svcCtx.WebsetModel.Update(l.ctx, &model.TWebset{
			Id:       in.Webset.Id,
			Title:    in.Webset.Title,
			Describe: in.Webset.Describe,
			AuthorId: in.Webset.AuthorInfo.Id,
			CoverUrl: in.Webset.CoverUrl,
			Category: 0,
			Status:   2,
		})
	} else if in.ActionType == 3 {
		// 删除
		err := l.svcCtx.WebsetModel.Update(l.ctx, &model.TWebset{
			Id:     in.Webset.Id,
			Status: 5,
		})
		if err != nil {
			logx.Error("delete webset failed, err: ", err)
			return &webset.PublishActionResp{
				StatusCode: internal.StatusRpcErr,
				StatusMsg:  "fail",
			}, err
		}
		return &webset.PublishActionResp{
			StatusCode: 1,
			StatusMsg:  "success",
		}, nil
	}

	// 未知操作类型
	logx.Error("unknown publish action type, action type: ", in.ActionType)
	return &webset.PublishActionResp{
		StatusCode: internal.StatusRpcErr,
		StatusMsg:  "fail, unknown action type",
	}, nil
}
