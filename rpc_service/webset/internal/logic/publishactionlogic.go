package logic

import (
	"context"

	"null-links/rpc_service/webset/internal/model"
	"null-links/rpc_service/webset/internal/svc"
	"null-links/rpc_service/webset/pb/webset"

	"github.com/zeromicro/go-zero/core/logx"
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
	// todo: add your logic here and delete this line
	if in.ActionType == 1 {
		// 发布
		websetDb := model.TWebset{
			Title:       in.Webset.Title,
			Describe:    in.Webset.Describe,
			AuthorId:    in.Webset.AuthorInfo.Id,
			CoverUrl:    in.Webset.CoverUrl,
			Category:    0,
			ViewCnt:     0,
			LikeCnt:     0,
			FavoriteCnt: 0,
			Status:      0,
		}
		insRes, err := l.svcCtx.WebsetModel.Insert(l.ctx, &websetDb)
		if err != nil {
			logx.Error("insert webset failed, err: ", err)
			return nil, err
		}
		rowsAffected, err := insRes.RowsAffected()
		if err != nil {
			logx.Error("insert webset failed, err: ", err)
			return nil, err
		}
		if rowsAffected == 0 {
			logx.Error("insert webset failed, rows affected: ", rowsAffected)
			return nil, err
		}
	} else if in.ActionType == 2 {
		// 更新

	} else if in.ActionType == 3 {
		// 删除

	} else {
		// 未知操作类型
		logx.Error("unknown publish action type, action type: ", in.ActionType)
	}

	return &webset.PublishActionResp{}, nil
}
