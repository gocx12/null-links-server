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
	tUserFieldNames          = builder.RawFieldNames(&TUser{})
	tUserRows                = strings.Join(tUserFieldNames, ",")
	tUserRowsExpectAutoSet   = strings.Join(stringx.Remove(tUserFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	tUserRowsWithPlaceHolder = strings.Join(stringx.Remove(tUserFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	tUserModel interface {
		Insert(ctx context.Context, data *TUser) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*TUser, error)
		FindOneByEmail(ctx context.Context, email string) (*TUser, error)
		FindOneByUsername(ctx context.Context, username string) (*TUser, error)
		Update(ctx context.Context, data *TUser) error
		Delete(ctx context.Context, id int64) error
	}

	defaultTUserModel struct {
		conn  sqlx.SqlConn
		table string
	}

	TUser struct {
		Id            int64     `db:"id"`             // 主键id
		Username      string    `db:"username"`       // 用户名
		Email         string    `db:"email"`          // 邮箱地址
		Password      string    `db:"password"`       // 密码
		AvatarUrl     string    `db:"avatar_url"`     // 头像地址
		BackgroundUrl string    `db:"background_url"` // 背景地址
		Signature     string    `db:"signature"`      // 个性签名
		FollowCount   int64     `db:"follow_count"`   // 关注数
		FollowerCount int64     `db:"follower_count"` // 粉丝数
		Status        int64     `db:"status"`         // 在库状态
		CreatedAt     time.Time `db:"created_at"`     // 创建时间
		UpdatedAt     time.Time `db:"updated_at"`     // 更新时间
	}
)

func newTUserModel(conn sqlx.SqlConn) *defaultTUserModel {
	return &defaultTUserModel{
		conn:  conn,
		table: "`t_user`",
	}
}

func (m *defaultTUserModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultTUserModel) FindOne(ctx context.Context, id int64) (*TUser, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", tUserRows, m.table)
	var resp TUser
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

func (m *defaultTUserModel) FindOneByEmail(ctx context.Context, email string) (*TUser, error) {
	var resp TUser
	query := fmt.Sprintf("select %s from %s where `email` = ? limit 1", tUserRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, email)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultTUserModel) FindOneByUsername(ctx context.Context, username string) (*TUser, error) {
	var resp TUser
	query := fmt.Sprintf("select %s from %s where `username` = ? limit 1", tUserRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, username)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultTUserModel) Insert(ctx context.Context, data *TUser) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, tUserRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.Username, data.Email, data.Password, data.AvatarUrl, data.BackgroundUrl, data.Signature, data.FollowCount, data.FollowerCount, data.Status)
	return ret, err
}

func (m *defaultTUserModel) Update(ctx context.Context, newData *TUser) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, tUserRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, newData.Username, newData.Email, newData.Password, newData.AvatarUrl, newData.BackgroundUrl, newData.Signature, newData.FollowCount, newData.FollowerCount, newData.Status, newData.Id)
	return err
}

func (m *defaultTUserModel) tableName() string {
	return m.table
}
