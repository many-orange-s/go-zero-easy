package svc

import (
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"go-zero-easy/service/user/cmd/api/internal/config"
	"go-zero-easy/service/user/model"
)

type ServiceContext struct {
	Config    config.Config
	UserModel model.UsermsgModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:    c,
		UserModel: model.NewUsermsgModel(conn, c.CacheRedis),
	}
}
