package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TFavoriteModel = (*customTFavoriteModel)(nil)

type (
	// TFavoriteModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTFavoriteModel.
	TFavoriteModel interface {
		tFavoriteModel
		GetFavoriteWebsetUserInfos(ctx context.Context, websetIds []int64, userId int64) ([]*TFavorite, error)
	}

	customTFavoriteModel struct {
		*defaultTFavoriteModel
	}
)

// NewTFavoriteModel returns a model for the database table.
func NewTFavoriteModel(conn sqlx.SqlConn) TFavoriteModel {
	return &customTFavoriteModel{
		defaultTFavoriteModel: newTFavoriteModel(conn),
	}
}

func (c *customTFavoriteModel) GetFavoriteWebsetUserInfos(ctx context.Context, websetIds []int64, userId int64) ([]*TFavorite, error) {
	query := fmt.Sprintf("select %s from %s where `webset_id` in (?) and `user_id`=?", tFavoriteRows, c.table)
	var resp []*TFavorite
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
