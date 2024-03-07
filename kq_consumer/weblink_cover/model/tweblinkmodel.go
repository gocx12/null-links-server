package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TWeblinkModel = (*customTWeblinkModel)(nil)

type (
	// TWeblinkModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTWeblinkModel.
	TWeblinkModel interface {
		tWeblinkModel
		UpdateCoverUrl(ctx context.Context, websetId, linkId int64, url string)
		FindByWebsetId(ctx context.Context, websetId int64) (weblinks []*TWeblink, err error)
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

func (c *customTWeblinkModel) UpdateCoverUrl(ctx context.Context, websetId, linkId int64, coverUrl string) {
	query := fmt.Sprintf("update %s set `cover_url` = ? where `webset_id` = ? and `link_id` = ?", c.tableName())
	c.conn.ExecCtx(ctx, query, coverUrl, websetId, linkId)
}

func (m *customTWeblinkModel) FindByWebsetId(ctx context.Context, websetId int64) (weblinks []*TWeblink, err error) {
	query := fmt.Sprintf("select %s from %s where `webset_id` = ?", tWeblinkRows, m.table)
	err = m.conn.QueryRowsCtx(ctx, &weblinks, query, websetId)
	return
}
