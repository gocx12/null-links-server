package svc

import (
	"null-links/chat_service/internal/config"
	"null-links/chat_service/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config    config.Config
	ChatModel model.TChatModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		ChatModel: model.NewTChatModel(sqlx.NewMysql(c.DataSource)),
	}
}
