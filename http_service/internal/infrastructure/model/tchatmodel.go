package model

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TChatModel = (*customTChatModel)(nil)

type (
	// TChatModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTChatModel.
	TChatModel interface {
		tChatModel
		withSession(session sqlx.Session) TChatModel
		FindChatList(ctx context.Context, websetId int64, page, pageSize int32) ([]*TChat, error)
		FindChatListChatId(ctx context.Context, websetId, lastChatId int64, page, pageSize int32) ([]*TChat, error)
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

func (c *customTChatModel) FindChatList(ctx context.Context, websetId int64, page, pageSize int32) ([]*TChat, error) {
	query := "select * from t_chat where webset_id = ? order by created_at desc limit ?, ?"
	var resp []*TChat
	err := c.conn.QueryRowsCtx(ctx, &resp, query, websetId, (page-1)*pageSize, pageSize)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (c *customTChatModel) FindChatListChatId(ctx context.Context, websetId, lastChatId int64, page, pageSize int32) ([]*TChat, error) {
	query := "select * from t_chat where webset_id = ? and chat_id < ? order by created_at desc limit ?, ?"
	var resp []*TChat
	err := c.conn.QueryRowsCtx(ctx, &resp, query, websetId, lastChatId, (page-1)*pageSize, pageSize)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}
