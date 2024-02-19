package svc

import (
	"null-links/chat_service/internal/config"
)

type ServiceContext struct {
	Config config.Config
	ChatModel:                model.NewTChatModel(sqlx.NewMysql(c.DataSource)),
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
