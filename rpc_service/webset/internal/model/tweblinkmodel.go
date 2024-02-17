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
		FindMulti(ctx context.Context, websetIds []int64) (weblinks []*TWeblink, err error)
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

func (m *customTWeblinkModel) FindMulti(ctx context.Context, websetIds []int64) (weblinks []*TWeblink, err error) {
	query := fmt.Sprintf("select %s from %s where `webset_id` = ?", tWeblinkRows, m.table)
	err = m.conn.QueryRowsCtx(ctx, &weblinks, query, websetIds)
	return

}
