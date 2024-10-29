package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TAdviceModel = (*customTAdviceModel)(nil)

type (
	// TAdviceModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTAdviceModel.
	TAdviceModel interface {
		tAdviceModel
		withSession(session sqlx.Session) TAdviceModel
	}

	customTAdviceModel struct {
		*defaultTAdviceModel
	}
)

// NewTAdviceModel returns a model for the database table.
func NewTAdviceModel(conn sqlx.SqlConn) TAdviceModel {
	return &customTAdviceModel{
		defaultTAdviceModel: newTAdviceModel(conn),
	}
}

func (m *customTAdviceModel) withSession(session sqlx.Session) TAdviceModel {
	return NewTAdviceModel(sqlx.NewSqlConnFromSession(session))
}
