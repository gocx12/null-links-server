package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TUserModel = (*customTUserModel)(nil)

type (
	// TUserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTUserModel.
	TUserModel interface {
		tUserModel
		FindOneByName(ctx context.Context, username string) (*TUser, error)
		FindMulti(ctx context.Context, userIds []int64) ([]*TUser, error)
		FindPasswordByEmail(ctx context.Context, email string) (*TUser, error)
	}

	customTUserModel struct {
		*defaultTUserModel
	}
)

// NewTUserModel returns a model for the database table.
func NewTUserModel(conn sqlx.SqlConn) TUserModel {
	return &customTUserModel{
		defaultTUserModel: newTUserModel(conn),
	}
}

func (c *customTUserModel) FindOneByName(ctx context.Context, username string) (*TUser, error) {
	query := fmt.Sprintf("select %s from %s where `username` = ? limit 1", tUserRows, c.table)
	var resp TUser
	err := c.conn.QueryRowCtx(ctx, &resp, query, username)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (c *customTUserModel) FindPasswordByEmail(ctx context.Context, email string) (*TUser, error) {
	query := fmt.Sprintf("select %s from %s where `email` = ? limit 1", tUserRows, c.table)
	var resp TUser
	err := c.conn.QueryRowCtx(ctx, &resp, query, email)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (c *customTUserModel) FindMulti(ctx context.Context, userIds []int64) ([]*TUser, error) {
	placeHodlers := make([]string, 0, len(userIds))
	for range userIds {
		placeHodlers = append(placeHodlers, "?")
	}
	query := fmt.Sprintf("select %s from %s where `id` in (%s)", tUserRows, c.table, strings.Join(placeHodlers, ","))
	var resp []*TUser
	err := c.conn.QueryRowsCtx(ctx, &resp, query, userIds)
	switch err {
	case nil:
		return resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
