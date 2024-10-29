package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TBalanceModel = (*customTBalanceModel)(nil)

type (
	// TBalanceModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTBalanceModel.
	TBalanceModel interface {
		tBalanceModel
		withSession(session sqlx.Session) TBalanceModel
	}

	customTBalanceModel struct {
		*defaultTBalanceModel
	}
)

// NewTBalanceModel returns a model for the database table.
func NewTBalanceModel(conn sqlx.SqlConn) TBalanceModel {
	return &customTBalanceModel{
		defaultTBalanceModel: newTBalanceModel(conn),
	}
}

func (m *customTBalanceModel) withSession(session sqlx.Session) TBalanceModel {
	return NewTBalanceModel(sqlx.NewSqlConnFromSession(session))
}
