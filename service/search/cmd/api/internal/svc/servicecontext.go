package svc

import (
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/zrpc"
	"go-zero-easy/service/search/cmd/api/internal/config"
	"go-zero-easy/service/search/model"
	"go-zero-easy/service/user/cmd/rpc/userclient"
)

type ServiceContext struct {
	Config    config.Config
	BookModel model.BookmsgModel
	LendModel model.LendModel
	UserRpc   userclient.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:    c,
		BookModel: model.NewBookmsgModel(conn, c.CacheRedis),
		LendModel: model.NewLendModel(conn, c.CacheRedis),
		UserRpc:   userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
	}
}
