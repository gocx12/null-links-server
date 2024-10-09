package webset

import (
	"context"
	"fmt"
	"strings"

	"null-links/http_service/internal/infrastructure/model"
	"null-links/http_service/internal/svc"
	"null-links/http_service/internal/types"
	"null-links/rpc_service/webset/pb/webset"

	"null-links/internal"

	"github.com/demdxx/gocast"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type PublishActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublishActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishActionLogic {
	return &PublishActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublishActionLogic) PublishAction(req *types.PublishActionReq) (resp *types.PublishActionResp, err error) {
	resp = &types.PublishActionResp{}

	var publishActionRpcReq webset.PublishActionReq
	if req.ActionType == 1 || req.ActionType == 2 {
		// 发布 或 修改
		weblinkListRpcReq := make([]*webset.WebLink, 0, len(req.WebLinkList))
		for _, weblink := range req.WebLinkList {
			weblinkListRpcReq = append(weblinkListRpcReq, &webset.WebLink{
				Url:      weblink.Url,
				Describe: weblink.Describe,
			})
		}
		publishActionRpcReq = webset.PublishActionReq{
			ActionType: req.ActionType,
			UserId:     req.AuthorId,
			Webset: &webset.Webset{
				Title:       req.Title,
				Describe:    req.Describe,
				CoverUrl:    req.CoverUrl,
				WebLinkList: weblinkListRpcReq,
			},
		}
	} else {
		resp.StatusCode = internal.StatusParamErr
		resp.StatusMsg = "未知操作类型"
		return
	}
	publishActionRpcResp, err := l.publishAction(&publishActionRpcReq)

	if err != nil || publishActionRpcResp.StatusCode != internal.StatusSuccess {
		if err != nil {
			logx.Error("call WebsetRpc failed, err: ", err)
			err = nil
		} else if publishActionRpcResp.StatusCode != internal.StatusSuccess {
			logx.Error("call WebsetRpc failed, err: ", publishActionRpcResp.StatusMsg)
		}
		resp.StatusCode = internal.StatusRpcErr
		if req.ActionType == 1 {
			resp.StatusMsg = "发布失败"
		} else if req.ActionType == 2 {
			resp.StatusMsg = "修改失败"
		} else if req.ActionType == 3 {
			resp.StatusMsg = "删除失败"
		} else {
			resp.StatusMsg = "操作失败"
		}
		return
	}

	resp.StatusCode = internal.StatusSuccess
	if req.ActionType == 1 {
		resp.StatusMsg = "发布成功"
	} else if req.ActionType == 2 {
		resp.StatusMsg = "修改成功"
	} else if req.ActionType == 3 {
		resp.StatusMsg = "删除成功"
	} else {
		resp.StatusMsg = "操作成功"
	}

	return
}

func (l *PublishActionLogic) publishAction(in *webset.PublishActionReq) (*webset.PublishActionResp, error) {
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
