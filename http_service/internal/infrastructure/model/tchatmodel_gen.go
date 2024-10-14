// Code generated by goctl. DO NOT EDIT.
// versions:
//  goctl version: 1.7.2

package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	tChatFieldNames          = builder.RawFieldNames(&TChat{})
	tChatRows                = strings.Join(tChatFieldNames, ",")
	tChatRowsExpectAutoSet   = strings.Join(stringx.Remove(tChatFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	tChatRowsWithPlaceHolder = strings.Join(stringx.Remove(tChatFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	tChatModel interface {
		Insert(ctx context.Context, data *TChat) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*TChat, error)
		Update(ctx context.Context, data *TChat) error
		Delete(ctx context.Context, id int64) error
	}

	defaultTChatModel struct {
		conn  sqlx.SqlConn
		table string
	}

	TChat struct {
		Id        int64     `db:"id"`         // 主键id
		ChatId    int64     `db:"chat_id"`    // 聊天消息id
		UserId    int64     `db:"user_id"`    // 用户id
		WebsetId  int64     `db:"webset_id"`  // 网页单id
		TopicId   int64     `db:"topic_id"`   // 话题id
		Content   string    `db:"content"`    // 消息内容
		Type      string    `db:"type"`       // 消息类型
		Status    int64     `db:"status"`     // 在库状态
		CreatedAt time.Time `db:"created_at"` // 创建时间
		UpdatedAt time.Time `db:"updated_at"` // 更新时间
	}
)

func newTChatModel(conn sqlx.SqlConn) *defaultTChatModel {
	return &defaultTChatModel{
		conn:  conn,
		table: "`t_chat`",
	}
}

func (m *defaultTChatModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultTChatModel) FindOne(ctx context.Context, id int64) (*TChat, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", tChatRows, m.table)
	var resp TChat
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultTChatModel) Insert(ctx context.Context, data *TChat) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?)", m.table, tChatRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.ChatId, data.UserId, data.WebsetId, data.TopicId, data.Content, data.Type, data.Status)
	return ret, err
}

func (m *defaultTChatModel) Update(ctx context.Context, data *TChat) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, tChatRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.ChatId, data.UserId, data.WebsetId, data.TopicId, data.Content, data.Type, data.Status, data.Id)
	return err
}

func (m *defaultTChatModel) tableName() string {
	return m.table
}
