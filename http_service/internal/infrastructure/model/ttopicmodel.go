package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TTopicModel = (*customTTopicModel)(nil)

type (
	// TTopicModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTTopicModel.
	TTopicModel interface {
		tTopicModel
		withSession(session sqlx.Session) TTopicModel
	}

	customTTopicModel struct {
		*defaultTTopicModel
	}
)

// NewTTopicModel returns a model for the database table.
func NewTTopicModel(conn sqlx.SqlConn) TTopicModel {
	return &customTTopicModel{
		defaultTTopicModel: newTTopicModel(conn),
	}
}

func (m *customTTopicModel) withSession(session sqlx.Session) TTopicModel {
	return NewTTopicModel(sqlx.NewSqlConnFromSession(session))
}
