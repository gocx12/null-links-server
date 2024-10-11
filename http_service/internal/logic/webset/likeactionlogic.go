package webset

import (
	"context"
	"fmt"

	"null-links/http_service/internal/infrastructure/model"
	"null-links/http_service/internal/svc"
	"null-links/http_service/internal/types"
	"null-links/internal"
	"null-links/rpc_service/webset/pb/webset"

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

	LikeActionRpcReq := &webset.LikeActionReq{
		UserId:     req.UserId,
		ActionType: req.ActionType,
		WebsetId:   req.WebsetId,
	}
	likeActionRpcResp, err := l.likeAction(LikeActionRpcReq)

	if likeActionRpcResp.StatusCode != internal.StatusSuccess {
		resp.StatusCode = internal.StatusRpcErr
		if req.ActionType == 1 {
			resp.StatusMsg = "点赞失败"

		} else if req.ActionType == 2 {
			resp.StatusMsg = "取消赞失败"
		}
	}

	resp.StatusCode = internal.StatusSuccess
	if req.ActionType == 1 {
		resp.StatusMsg = "点赞成功"

	} else if req.ActionType == 2 {
		resp.StatusMsg = "取消赞成功"
	}

	return
}

func (l *LikeActionLogic) likeAction(in *webset.LikeActionReq) (*webset.LikeActionResp, error) {
	// hash key: webset_id::user_id value:status
	likeActionResp := webset.LikeActionResp{
		StatusCode: internal.StatusSuccess,
		StatusMsg:  "success",
	}

	if in.ActionType == 1 {
		// type == 1, 点赞

		// // redis事务
		// _, err := l.svcCtx.RedisClient.TxPipelined(l.ctx, func(pipe redis.Pipeliner) error {
		// 	key := gocast.ToString(in.WebsetId) + "::" + gocast.ToString(in.UserId)
		// 	pipe.HSet(l.ctx, RdsKeyUserWebsetLiked, key, 1)
		// 	// 点赞数+1
		// 	pipe.HIncrBy(l.ctx, RdsKeyWebsetLikedCnt, gocast.ToString(in.WebsetId), 1)
		// 	return nil
		// })

		// like 与 webset 在同一个数据库中，因此可以使用本地事务
		err := l.svcCtx.WebsetModel.GetConn().TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {

			likeStatus, err := l.svcCtx.LikeModel.FindStatusWebsetIdUserIdTrans(l.ctx, in.WebsetId, in.UserId, session)
			if err != nil && err != sqlx.ErrNotFound {
				return err
			} else if err == sqlx.ErrNotFound {
				// 点赞记录不存在，插入点赞记录
				res, err := l.svcCtx.LikeModel.Insert(l.ctx, &model.TLike{
					UserId:   in.UserId,
					WebsetId: in.WebsetId,
					Status:   1,
				})
				if err != nil {
					return err
				}
				if rowsAffected, err := res.RowsAffected(); err != nil {
					return err
				} else if rowsAffected == 0 {
					return fmt.Errorf("update like record failed, rows affected: %d", rowsAffected)
				}
			}

			if likeStatus == 1 {
				return fmt.Errorf("like record already exists, user_id:%d, webset_id:%d", in.UserId, in.WebsetId)
			} else {
				// 点赞记录已存在，修改状态
				res, err := l.svcCtx.LikeModel.UpdateStatusTrans(l.ctx, in.WebsetId, in.UserId, 1, session)
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
			r, err := l.svcCtx.WebsetModel.UpdateLikeCntTrans(l.ctx, 1, in.WebsetId, session)
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
			likeActionResp.StatusCode = internal.StatusRpcErr
			likeActionResp.StatusMsg = "like webset failed"
		}

		likeActionResp.StatusCode = internal.StatusSuccess
		likeActionResp.StatusMsg = "success"
	} else if in.ActionType == 2 {
		// type == 2, 取消点赞
		// redis事务
		// _, err := l.svcCtx.RedisClient.TxPipelined(l.ctx, func(pipe redis.Pipeliner) error {
		// 	key := gocast.ToString(in.WebsetId) + "::" + gocast.ToString(in.UserId)
		// 	pipe.HSet(l.ctx, RdsKeyUserWebsetLiked, key, 2)
		// 	// 点赞数-1
		// 	pipe.HIncrBy(l.ctx, RdsKeyWebsetLikedCnt, gocast.ToString(in.WebsetId), -1)
		// 	return nil
		// })

		// like 与 webset 在同一个数据库中，因此可以使用本地事务
		err := l.svcCtx.WebsetModel.GetConn().TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
			likeStatus, err := l.svcCtx.LikeModel.FindStatusWebsetIdUserIdTrans(l.ctx, in.WebsetId, in.UserId, session)
			if err != nil && err != sqlx.ErrNotFound {
				return err
			} else if err == sqlx.ErrNotFound {
				// 点赞记录不存在，插入点赞记录
				res, err := l.svcCtx.LikeModel.Insert(l.ctx, &model.TLike{
					UserId:   in.UserId,
					WebsetId: in.WebsetId,
					Status:   2,
				})
				if err != nil {
					return err
				}
				if rowsAffected, err := res.RowsAffected(); err != nil {
					return err
				} else if rowsAffected == 0 {
					return fmt.Errorf("update like record failed, rows affected: %d", rowsAffected)
				}
			}

			if likeStatus == 2 {
				return fmt.Errorf("cancel like record already exists, user_id:%d, webset_id:%d", in.UserId, in.WebsetId)
			} else {
				// 点赞记录已存在，修改状态
				res, err := l.svcCtx.LikeModel.UpdateStatusTrans(l.ctx, in.WebsetId, in.UserId, 2, session)
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
			r, err := l.svcCtx.WebsetModel.UpdateLikeCntTrans(l.ctx, -1, in.WebsetId, session)
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
			likeActionResp.StatusCode = internal.StatusRpcErr
			likeActionResp.StatusMsg = "like webset failed"
		}

		likeActionResp.StatusCode = internal.StatusSuccess
		likeActionResp.StatusMsg = "success"
	} else {
		// 未知操作类型
		logx.Error("unknown like action type")
		likeActionResp.StatusCode = internal.StatusRpcErr
		likeActionResp.StatusMsg = "unknown like action type"
	}

	return &likeActionResp, nil
}
