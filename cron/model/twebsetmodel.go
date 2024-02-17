package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TWebsetModel = (*customTWebsetModel)(nil)

type (
	// TWebsetModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTWebsetModel.
	TWebsetModel interface {
		tWebsetModel
	}

	customTWebsetModel struct {
		*defaultTWebsetModel
	}
)

// NewTWebsetModel returns a model for the database table.
func NewTWebsetModel(conn sqlx.SqlConn) TWebsetModel {
	return &customTWebsetModel{
		defaultTWebsetModel: newTWebsetModel(conn),
	}
}
