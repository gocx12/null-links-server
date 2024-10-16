package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TLikeModel = (*customTLikeModel)(nil)

type (
	// TLikeModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTLikeModel.
	TLikeModel interface {
		tLikeModel
		GetLikeWebsetUserInfos(ctx context.Context, websetIds []int64, userId int64) ([]*TLike, error)
		GetLikeWebsetUserInfo(ctx context.Context, websetId int64, userId int64) (*TLike, error)
		UpdateStatusTrans(ctx context.Context, websetId, userId int64, action int32, session sqlx.Session) (sql.Result, error)
		FindStatusWebsetIdUserIdTrans(ctx context.Context, websetId, userId int64, session sqlx.Session) (int32, error)
	}

	customTLikeModel struct {
		*defaultTLikeModel
	}
)

// NewTLikeModel returns a model for the database table.
func NewTLikeModel(conn sqlx.SqlConn) TLikeModel {
	return &customTLikeModel{
		defaultTLikeModel: newTLikeModel(conn),
	}
}

func (c *customTLikeModel) GetLikeWebsetUserInfo(ctx context.Context, websetId int64, userId int64) (*TLike, error) {
	query := fmt.Sprintf("select %s from %s where `webset_id`=? and `user_id`=?", tLikeRows, c.table)
	var resp TLike
	err := c.conn.QueryRowCtx(ctx, &resp, query, websetId, userId)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (c *customTLikeModel) GetLikeWebsetUserInfos(ctx context.Context, websetIds []int64, userId int64) ([]*TLike, error) {
	query := fmt.Sprintf("select %s from %s where `webset_id` in (?) and `user_id`=?", tLikeRows, c.table)
	args := make([]interface{}, 0, len(websetIds)+1)
	for _, v := range websetIds {
		args = append(args, v)
	}
	args = append(args, userId)

	var resp []*TLike
	err := c.conn.QueryRowsCtx(ctx, &resp, query, args...)
	switch err {
	case nil:
		return resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (c *customTLikeModel) UpdateStatusTrans(ctx context.Context, websetId, userId int64, action int32, session sqlx.Session) (sql.Result, error) {
	query := fmt.Sprintf("update %s set `status` = ? where `user_id` = ? and `webset_id` = ?", c.table)
	res, err := session.ExecCtx(ctx, query, action, userId, websetId)
	switch err {
	case nil:
		return res, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (c *customTLikeModel) InsertTrans(ctx context.Context, data *TLike, session sqlx.Session) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (`webset_id`, `user_id`, `status`) values (?, ?, ?)", c.table)
	ret, err := session.ExecCtx(ctx, query, data.WebsetId, data.UserId, data.Status)
	return ret, err
}

func (c *customTLikeModel) FindStatusWebsetIdUserIdTrans(ctx context.Context, websetId, userId int64, session sqlx.Session) (int32, error) {
	query := fmt.Sprintf("select `status` from %s where `webset_id` = ? and `user_id` = ?", c.table)
	var status int32
	err := session.QueryRowCtx(ctx, &status, query, websetId, userId)
	switch err {
	case nil:
		return status, nil
	case sqlx.ErrNotFound:
		return 0, ErrNotFound
	default:
		return 0, err
	}
}
