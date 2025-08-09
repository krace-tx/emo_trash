package svc

import (
	"github.com/krace-tx/emo_trash/app/rpc/sso/internal/config"
	"github.com/krace-tx/emo_trash/pkg/db/rdb"
)

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	if err := rdb.InitMySQL(c.MysqlConf); err != nil {
		panic(err)
	}

	return &ServiceContext{
		Config: c,
	}
}
