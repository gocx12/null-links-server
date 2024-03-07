package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TWeblinkModel = (*customTWeblinkModel)(nil)

type (
	// TWeblinkModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTWeblinkModel.
	TWeblinkModel interface {
		tWeblinkModel
		FindByWebsetId(ctx context.Context, websetId int64) (weblinks []*TWeblink, err error)
		BulkInsert(ctx context.Context, data []TWeblink) (sql.Result, error)
		BulkInsertTrans(ctx context.Context, data []TWeblink, session sqlx.Session) (sql.Result, error)
	}

	customTWeblinkModel struct {
		*defaultTWeblinkModel
	}
)

// NewTWeblinkModel returns a model for the database table.
func NewTWeblinkModel(conn sqlx.SqlConn) TWeblinkModel {
	return &customTWeblinkModel{
		defaultTWeblinkModel: newTWeblinkModel(conn),
	}
}

func (m *customTWeblinkModel) FindByWebsetId(ctx context.Context, websetId int64) (weblinks []*TWeblink, err error) {
	query := fmt.Sprintf("select %s from %s where `webset_id` = ?", tWeblinkRows, m.table)
	err = m.conn.QueryRowsCtx(ctx, &weblinks, query, websetId)
	return
}

func (m *customTWeblinkModel) getBulkInsertQuery(data []TWeblink) (string, []interface{}) {
	placeHolderStrs := make([]string, 0, len(data))
	valueArgs := make([]interface{}, 0, len(data)*7)
	for _, d := range data {
		placeHolderStrs = append(placeHolderStrs, "(?, ?, ?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs, d.LinkId)
		valueArgs = append(valueArgs, d.WebsetId)
		valueArgs = append(valueArgs, d.AuthorId)
		valueArgs = append(valueArgs, d.Describe)
		valueArgs = append(valueArgs, d.Url)
		valueArgs = append(valueArgs, d.CoverUrl)
		valueArgs = append(valueArgs, d.Status)
	}
	query := fmt.Sprintf("insert into %s (`link_id`, `webset_id`, `author_id`, `describe`, `url`, `cover_url`, `status`) values %s", m.table, strings.Join(placeHolderStrs, ","))
	return query, valueArgs
}

func (m *customTWeblinkModel) BulkInsert(ctx context.Context, data []TWeblink) (sql.Result, error) {
	query, valueArgs := m.getBulkInsertQuery(data)
	return m.conn.ExecCtx(ctx, query, valueArgs...)
}

func (m *customTWeblinkModel) BulkInsertTrans(ctx context.Context, data []TWeblink, session sqlx.Session) (sql.Result, error) {
	query, valueArgs := m.getBulkInsertQuery(data)
	return session.ExecCtx(ctx, query, valueArgs...)
}
