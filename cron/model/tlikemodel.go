package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TLikeModel = (*customTLikeModel)(nil)

type (
	// TLikeModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTLikeModel.
	TLikeModel interface {
		tLikeModel
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
