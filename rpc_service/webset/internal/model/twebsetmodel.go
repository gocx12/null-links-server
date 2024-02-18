package model

import (
	"context"
	"fmt"

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
	query := fmt.Sprintf("select %s from %s where status=1 order by update_at desc limit ?, ?", tWebsetRows, m.table)
	err = m.conn.QueryRowsCtx(ctx, &websets, query, page, pageSize)
	return
}

func (m *customTWebsetModel) FindPublishList(ctx context.Context, authorId int64, page, pageSize int32) (websets []*TWebset, err error) {
	query := fmt.Sprintf("select %s from %s where author_id=? order by update_at desc limit ?, ?", tWebsetRows, m.table)
	err = m.conn.QueryRowsCtx(ctx, &websets, query, authorId, page, pageSize)
	return
}
