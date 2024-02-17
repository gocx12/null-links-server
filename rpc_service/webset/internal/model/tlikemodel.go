package model

import (
	"context"
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

func (c *customTLikeModel) GetLikeWebsetUserInfos(ctx context.Context, websetIds []int64, userId int64) ([]*TLike, error) {
	query := fmt.Sprintf("select %s from %s where `webset_id` in (?) and `user_id`=?", tLikeRows, c.table)
	var resp []*TLike
	err := c.conn.QueryRowsCtx(ctx, &resp, query, websetIds, userId)
	switch err {
	case nil:
		return resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
