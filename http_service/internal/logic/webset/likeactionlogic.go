package webset

import (
	"context"
	"fmt"

	"null-links/http_service/internal/infrastructure/model"
	"null-links/http_service/internal/svc"
	"null-links/http_service/internal/types"
	"null-links/internal"

	"github.com/demdxx/gocast"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type LikeActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLikeActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikeActionLogic {
	return &LikeActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LikeActionLogic) LikeAction(req *types.LikeActionReq) (resp *types.LikeActionResp, err error) {
	resp = &types.LikeActionResp{}
	logx.Debug("LikeAction|req=", req)

	userId := gocast.ToInt64(l.ctx.Value("userId"))
	switch internal.LikeActionTypeEnum(req.ActionType) {
	case internal.DoLike:
		err = l.doLike(userId, req.WebsetId)
		if err != nil {
			resp.StatusMsg = "like failed"
			resp.StatusCode = internal.StatusErr
		}
	case internal.DoCancelLike:
		err = l.doCancelLike(userId, req.WebsetId)
		if err != nil {
			resp.StatusMsg = "cancel like failed"
			resp.StatusCode = internal.StatusErr
		}
	default:
		resp.StatusCode = internal.StatusRpcErr
		resp.StatusMsg = "unknown like action type"
	}

	if err != nil {
		logx.Error("LikeAction|err=", err)
		return
	}

	resp.StatusCode = internal.StatusSuccess
	resp.StatusMsg = "success"
	return
}

func (l *LikeActionLogic) doLike(userId int64, websetId int64) (err error) {
	// like 与 webset 在同一个数据库中，因此可以使用本地事务
	err = l.svcCtx.WebsetModel.GetConn().TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {

		likeStatus, err := l.svcCtx.LikeModel.FindStatusWebsetIdUserIdTrans(l.ctx, websetId, userId, session)
		if err != nil && err != sqlx.ErrNotFound {
			return err
		} else if err == sqlx.ErrNotFound {
			// 点赞记录不存在，插入点赞记录
			res, err := l.svcCtx.LikeModel.Insert(l.ctx, &model.TLike{
				UserId:   userId,
				WebsetId: websetId,
				Status:   gocast.ToInt64(internal.Like.Code()),
			})
			if err != nil {
				return err
			}
			if rowsAffected, err := res.RowsAffected(); err != nil {
				return err
			} else if rowsAffected == 0 {
				return fmt.Errorf("insert like record failed, rows affected=%d", rowsAffected)
			}
		} else {
			if likeStatus == internal.Like.Code() {
				return fmt.Errorf("like record already exists, user_id=%d, webset_id=%d", userId, websetId)
			}

			// 点赞记录已存在，修改状态
			res, err := l.svcCtx.LikeModel.UpdateStatusTrans(l.ctx, websetId, userId, internal.Like.Code(), session)
			if err != nil {
				return err
			}
			if rowsAffected, err := res.RowsAffected(); err != nil {
				return err
			} else if rowsAffected == 0 {
				return fmt.Errorf("update like record failed, rows affected: %d", rowsAffected)
			}
		}

		// 更新webset点赞数
		r, err := l.svcCtx.WebsetModel.UpdateLikeCntTrans(l.ctx, 1, websetId, session)
		if err != nil {
			return err
		}
		if rowsAffected, err := r.RowsAffected(); err != nil {
			return err
		} else if rowsAffected == 0 {
			return fmt.Errorf("update webset failed, rows affected: %d", rowsAffected)
		}

		return nil
	})

	if err != nil {
		logx.Error("like webset failed, err: ", err)
	}

	return
}

func (l *LikeActionLogic) doCancelLike(userId int64, websetId int64) (err error) {
	// like 与 webset 在同一个数据库中，因此可以使用本地事务
	err = l.svcCtx.WebsetModel.GetConn().TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		likeStatus, err := l.svcCtx.LikeModel.FindStatusWebsetIdUserIdTrans(l.ctx, websetId, userId, session)
		if err != nil && err != sqlx.ErrNotFound {
			return err
		} else if err == sqlx.ErrNotFound {
			// 点赞记录不存在，插入点赞记录
			res, err := l.svcCtx.LikeModel.Insert(l.ctx, &model.TLike{
				UserId:   userId,
				WebsetId: websetId,
				Status:   gocast.ToInt64(internal.UnLike.Code()),
			})
			if err != nil {
				return err
			}
			if rowsAffected, err := res.RowsAffected(); err != nil {
				return err
			} else if rowsAffected == 0 {
				return fmt.Errorf("insert like record failed, rows affected: %d", rowsAffected)
			}
		} else {
			if likeStatus == internal.UnLike.Code() {
				return fmt.Errorf("cancel like record already exists, user_id:%d, webset_id:%d", userId, websetId)
			}

			// 点赞记录已存在，修改状态
			res, err := l.svcCtx.LikeModel.UpdateStatusTrans(l.ctx, websetId, userId, internal.UnLike.Code(), session)
			if err != nil {
				return err
			}
			if rowsAffected, err := res.RowsAffected(); err != nil {
				return err
			} else if rowsAffected == 0 {
				return fmt.Errorf("update like record failed, rows affected: %d", rowsAffected)
			}
		}

		// 更新webset点赞数
		r, err := l.svcCtx.WebsetModel.UpdateLikeCntTrans(l.ctx, -1, websetId, session)
		if err != nil {
			return err
		}
		if rowsAffected, err := r.RowsAffected(); err != nil {
			return err
		} else if rowsAffected == 0 {
			return fmt.Errorf("update webset failed, rows affected: %d", rowsAffected)
		}

		return nil
	})

	if err != nil {
		logx.Error("like webset failed, err: ", err)
	}

	return
}
