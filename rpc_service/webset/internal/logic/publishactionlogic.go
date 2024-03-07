package logic

import (
	"context"
	"fmt"
	"strings"

	"null-links/rpc_service/webset/internal/model"
	"null-links/rpc_service/webset/internal/svc"
	"null-links/rpc_service/webset/pb/webset"

	"null-links/internal"

	"github.com/demdxx/gocast"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
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
		// 发布
		websetDb := model.TWebset{
			Title:       in.Webset.Title,
			Describe:    in.Webset.Describe,
			AuthorId:    in.UserId,
			CoverUrl:    in.Webset.CoverUrl,
			Category:    0, // 分区，暂时不用
			ViewCnt:     0,
			LikeCnt:     0,
			FavoriteCnt: 0,
			Status:      2, //status==2 待审核
		}
		// 插入weblinks
		webLinkListDb := make([]model.TWeblink, 0, len(in.Webset.WebLinkList))
		for i, webLink := range in.Webset.WebLinkList {
			if !strings.HasPrefix(webLink.Url, "http") || !strings.HasPrefix(webLink.Url, "https") {
				webLink.Url = "https://" + webLink.Url
			}
			webLinkListDb = append(webLinkListDb, model.TWeblink{
				LinkId:   int64(i),
				WebsetId: in.WebsetId,
				AuthorId: in.UserId,
				Describe: webLink.Describe,
				Url:      webLink.Url,
				CoverUrl: webLink.CoverUrl,
				Status:   2, // status==2 待审核
			})
		}

		// weblink 与 webset 在同一个数据库中，因此可以使用本地事务
		err := l.svcCtx.WebsetModel.GetConn().TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
			// 插入webset
			r, err := l.svcCtx.WebsetModel.InsertTrans(l.ctx, &websetDb, session)
			if err != nil {
				return err
			}
			if rowsAffected, err := r.RowsAffected(); err != nil {
				return err
			} else if rowsAffected == 0 {
				return fmt.Errorf("insert webset failed, rows affected: %d", rowsAffected)
			}

			lastInsertId, err := r.LastInsertId()
			if err != nil {
				return err
			}
			for i := range webLinkListDb {
				webLinkListDb[i].WebsetId = lastInsertId
			}

			// kafka pusher
			data := gocast.ToString(lastInsertId)
			if err := l.svcCtx.WlCoverKqConsumser.Push(data); err != nil {
				logx.Error("WlCoverKqConsumser Push error:", err)
			}

			// 批量插入weblink
			r, err = l.svcCtx.WeblinkModel.BulkInsertTrans(l.ctx, webLinkListDb, session)
			if err != nil {
				return err
			}

			if rowsAffected, err := r.RowsAffected(); err != nil {
				return err
			} else if rowsAffected == 0 {
				return fmt.Errorf("insert weblinks failed, rows affected: %d", rowsAffected)
			}

			return nil
		})

		if err != nil {
			logx.Error("insert weblinks failed, err: ", err)
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
		// 更新
		l.svcCtx.WebsetModel.Update(l.ctx, &model.TWebset{
			Id:       in.Webset.Id,
			Title:    in.Webset.Title,
			Describe: in.Webset.Describe,
			AuthorId: in.Webset.AuthorInfo.Id,
			CoverUrl: in.Webset.CoverUrl,
			Category: 0,
			Status:   2, // status==2 待审核
		})
	}

	// 未知操作类型
	logx.Error("unknown publish action type, action type: ", in.ActionType)
	return &webset.PublishActionResp{
		StatusCode: internal.StatusRpcErr,
		StatusMsg:  "fail, unknown action type",
	}, nil
}
