package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TWebsetModel = (*customTWebsetModel)(nil)

type (
	// TWebsetModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTWebsetModel.
	TWebsetModel interface {
		tWebsetModel
		FindRecent(ct1x context.Context, page, pageSize int32) (websets []*TWebset, err error)
		FindPublishList(ctx context.Context, authorId int64, page, pageSize int32) (websets []*TWebset, err error)
		InsertTrans(ctx context.Context, data *TWebset, session sqlx.Session) (sql.Result, error)
		UpdateLikeCntTrans(ctx context.Context, incr int32, websetId int64, session sqlx.Session) (sql.Result, error)
		UpdateStatus(ctx context.Context, status int64, id int64) error
		UpdateWebsetInfo(ctx context.Context, data *TWebset) error
		GetConn() sqlx.SqlConn
	}

	customTWebsetModel struct {
		*defaultTWebsetModel
	}
)

// NewTWebsetModel returns a model for the database table.
func NewTWebsetModel(conn sqlx.SqlConn) TWebsetModel {
	return &customTWebsetModel{
		defaultTWebsetModel: newTWebsetModel(conn),
	}
}

func (m *customTWebsetModel) FindRecent(ctx context.Context, page, pageSize int32) (websets []*TWebset, err error) {
	offset := (page - 1) * pageSize

	query := fmt.Sprintf("select %s from %s where status = 1 order by created_at desc limit ? offset ?", tWebsetRows, m.table)
	err = m.conn.QueryRowsCtx(ctx, &websets, query, pageSize, offset)
	return
}

func (m *customTWebsetModel) FindPublishList(ctx context.Context, authorId int64, page, pageSize int32) (websets []*TWebset, err error) {
	query := fmt.Sprintf("select %s from %s where author_id=? order by created_at desc limit ?, ?", tWebsetRows, m.table)
	err = m.conn.QueryRowsCtx(ctx, &websets, query, authorId, page, pageSize)
	return
}

func (m *customTWebsetModel) GetConn() sqlx.SqlConn {
	return m.conn
}

func (m *defaultTWebsetModel) InsertTrans(ctx context.Context, data *TWebset, session sqlx.Session) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, tWebsetRowsExpectAutoSet)
	ret, err := session.ExecCtx(ctx, query, data.Title, data.AuthorId, data.Describe, data.CoverUrl, data.Category, data.ViewCnt, data.LikeCnt, data.FavoriteCnt, data.Status, data.UpdatedAt)
	return ret, err
}

func (m *defaultTWebsetModel) UpdateLikeCntTrans(ctx context.Context, incr int32, websetId int64, session sqlx.Session) (sql.Result, error) {
	query := fmt.Sprintf("update %s set like_cnt=like_cnt+? where id=?", m.table)
	ret, err := session.ExecCtx(ctx, query, incr, websetId)
	return ret, err
}

// 更新webset信息，点赞数量等字段不更新
func (m *defaultTWebsetModel) UpdateWebsetInfo(ctx context.Context, data *TWebset) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, strings.Join([]string{"title", "author_id", "describe", "cover_url", "category", "status", "updated_at"}, "=?, ")+"=?")
	_, err := m.conn.ExecCtx(ctx, query, data.Title, data.AuthorId, data.Describe, data.CoverUrl, data.Category, data.Status, data.UpdatedAt, data.Id)
	return err
}

func (m *defaultTWebsetModel) UpdateStatus(ctx context.Context, status int64, id int64) error {
	query := fmt.Sprintf("update %s set `status` =  ? where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, status, id)
	return err
}
