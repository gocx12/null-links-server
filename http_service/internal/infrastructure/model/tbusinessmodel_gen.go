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
	tBusinessFieldNames          = builder.RawFieldNames(&TBusiness{})
	tBusinessRows                = strings.Join(tBusinessFieldNames, ",")
	tBusinessRowsExpectAutoSet   = strings.Join(stringx.Remove(tBusinessFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	tBusinessRowsWithPlaceHolder = strings.Join(stringx.Remove(tBusinessFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	tBusinessModel interface {
		Insert(ctx context.Context, data *TBusiness) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*TBusiness, error)
		Update(ctx context.Context, data *TBusiness) error
		Delete(ctx context.Context, id int64) error
	}

	defaultTBusinessModel struct {
		conn  sqlx.SqlConn
		table string
	}

	TBusiness struct {
		Id        int64     `db:"id"`         // 主键id
		Business  string    `db:"business"`   // 业务名
		CreatedAt time.Time `db:"created_at"` // 创建时间
		UpdatedAt time.Time `db:"updated_at"` // 更新时间
	}
)

func newTBusinessModel(conn sqlx.SqlConn) *defaultTBusinessModel {
	return &defaultTBusinessModel{
		conn:  conn,
		table: "`t_business`",
	}
}

func (m *defaultTBusinessModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultTBusinessModel) FindOne(ctx context.Context, id int64) (*TBusiness, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", tBusinessRows, m.table)
	var resp TBusiness
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

func (m *defaultTBusinessModel) Insert(ctx context.Context, data *TBusiness) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?)", m.table, tBusinessRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.Business)
	return ret, err
}

func (m *defaultTBusinessModel) Update(ctx context.Context, data *TBusiness) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, tBusinessRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.Business, data.Id)
	return err
}

func (m *defaultTBusinessModel) tableName() string {
	return m.table
}
