package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TFavoriteModel = (*customTFavoriteModel)(nil)

type (
	// TFavoriteModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTFavoriteModel.
	TFavoriteModel interface {
		tFavoriteModel
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
