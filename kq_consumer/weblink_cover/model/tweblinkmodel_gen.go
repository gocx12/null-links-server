// Code generated by goctl. DO NOT EDIT.

package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	tWeblinkFieldNames          = builder.RawFieldNames(&TWeblink{})
	tWeblinkRows                = strings.Join(tWeblinkFieldNames, ",")
	tWeblinkRowsExpectAutoSet   = strings.Join(stringx.Remove(tWeblinkFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	tWeblinkRowsWithPlaceHolder = strings.Join(stringx.Remove(tWeblinkFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	tWeblinkModel interface {
		Insert(ctx context.Context, data *TWeblink) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*TWeblink, error)
		Update(ctx context.Context, data *TWeblink) error
		Delete(ctx context.Context, id int64) error
	}

	defaultTWeblinkModel struct {
		conn  sqlx.SqlConn
		table string
	}

	TWeblink struct {
		Id        int64        `db:"id"`         // 主键id
		LinkId    int64        `db:"link_id"`    // 网页id
		WebsetId  int64        `db:"webset_id"`  // 网页单id
		AuthorId  int64        `db:"author_id"`  // 添加者id
		Describe  string       `db:"describe"`   // 描述
		Url       string       `db:"url"`        // 网址
		CoverUrl  string       `db:"cover_url"`  // 封面地址
		ClickCnt  int64        `db:"click_cnt"`  // 点击数
		Status    int64        `db:"status"`     // 在库状态
		CreatedAt time.Time    `db:"created_at"` // 创建时间
		UpdatedAt time.Time    `db:"updated_at"` // 更新时间
		DeletedAt sql.NullTime `db:"deleted_at"` // 删除时间
	}
)

func newTWeblinkModel(conn sqlx.SqlConn) *defaultTWeblinkModel {
	return &defaultTWeblinkModel{
		conn:  conn,
		table: "`t_weblink`",
	}
}

func (m *defaultTWeblinkModel) withSession(session sqlx.Session) *defaultTWeblinkModel {
	return &defaultTWeblinkModel{
		conn:  sqlx.NewSqlConnFromSession(session),
		table: "`t_weblink`",
	}
}

func (m *defaultTWeblinkModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultTWeblinkModel) FindOne(ctx context.Context, id int64) (*TWeblink, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", tWeblinkRows, m.table)
	var resp TWeblink
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultTWeblinkModel) Insert(ctx context.Context, data *TWeblink) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, tWeblinkRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.LinkId, data.WebsetId, data.AuthorId, data.Describe, data.Url, data.CoverUrl, data.ClickCnt, data.Status, data.DeletedAt)
	return ret, err
}

func (m *defaultTWeblinkModel) Update(ctx context.Context, data *TWeblink) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, tWeblinkRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.LinkId, data.WebsetId, data.AuthorId, data.Describe, data.Url, data.CoverUrl, data.ClickCnt, data.Status, data.DeletedAt, data.Id)
	return err
}

func (m *defaultTWeblinkModel) tableName() string {
	return m.table
}