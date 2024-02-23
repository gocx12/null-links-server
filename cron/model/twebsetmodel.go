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
		UpdateLikeCnt(ctx context.Context, likeCnt int64, websetId int64) error
		UpdateFavoriteCnt(ctx context.Context, favoriteCnt int64, websetId int64) error
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

func (m *customTWebsetModel) UpdateLikeCnt(ctx context.Context, likeCnt int64, websetId int64) error {
	query := fmt.Sprintf("update %s set `like_cnt` = ? where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, likeCnt, websetId)
	return err
}

func (m *customTWebsetModel) UpdateFavoriteCnt(ctx context.Context, favoriteCnt int64, websetId int64) error {
	query := fmt.Sprintf("update %s set `favorite_cnt` = ? where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, favoriteCnt, websetId)
	return err
}
