package svc

import (
	"null-links/rpc_service/user/internal/config"
	"null-links/rpc_service/user/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config    config.Config
	UserModel model.TUserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		UserModel: model.NewTUserModel(sqlx.NewMysql(c.DataSource)),
	}
}
