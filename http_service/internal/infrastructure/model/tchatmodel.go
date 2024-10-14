package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TChatModel = (*customTChatModel)(nil)

type (
	// TChatModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTChatModel.
	TChatModel interface {
		tChatModel
		withSession(session sqlx.Session) TChatModel
	}

	customTChatModel struct {
		*defaultTChatModel
	}
)

// NewTChatModel returns a model for the database table.
func NewTChatModel(conn sqlx.SqlConn) TChatModel {
	return &customTChatModel{
		defaultTChatModel: newTChatModel(conn),
	}
}

func (m *customTChatModel) withSession(session sqlx.Session) TChatModel {
	return NewTChatModel(sqlx.NewSqlConnFromSession(session))
}
