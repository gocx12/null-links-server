package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TPayHistoryModel = (*customTPayHistoryModel)(nil)

type (
	// TPayHistoryModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTPayHistoryModel.
	TPayHistoryModel interface {
		tPayHistoryModel
		withSession(session sqlx.Session) TPayHistoryModel
	}

	customTPayHistoryModel struct {
		*defaultTPayHistoryModel
	}
)

// NewTPayHistoryModel returns a model for the database table.
func NewTPayHistoryModel(conn sqlx.SqlConn) TPayHistoryModel {
	return &customTPayHistoryModel{
		defaultTPayHistoryModel: newTPayHistoryModel(conn),
	}
}

func (m *customTPayHistoryModel) withSession(session sqlx.Session) TPayHistoryModel {
	return NewTPayHistoryModel(sqlx.NewSqlConnFromSession(session))
}
