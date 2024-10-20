package model

import (
	"context"
	"fmt"
	"strings"

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
	questionMarks := strings.Repeat("?,", len(websetIds))
	query := fmt.Sprintf("select %s from %s where `webset_id` in (%s) and `user_id`=?", tLikeRows, c.table, questionMarks[:len(questionMarks)-1])
	args := make([]interface{}, 0, len(websetIds)+1)
	for _, v := range websetIds {
		args = append(args, v)
	}
	args = append(args, userId)

	var resp []*TFavorite
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
