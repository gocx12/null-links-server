package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TBusinessModel = (*customTBusinessModel)(nil)

type (
	// TBusinessModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTBusinessModel.
	TBusinessModel interface {
		tBusinessModel
		withSession(session sqlx.Session) TBusinessModel
	}

	customTBusinessModel struct {
		*defaultTBusinessModel
	}
)

// NewTBusinessModel returns a model for the database table.
func NewTBusinessModel(conn sqlx.SqlConn) TBusinessModel {
	return &customTBusinessModel{
		defaultTBusinessModel: newTBusinessModel(conn),
	}
}

func (m *customTBusinessModel) withSession(session sqlx.Session) TBusinessModel {
	return NewTBusinessModel(sqlx.NewSqlConnFromSession(session))
}
