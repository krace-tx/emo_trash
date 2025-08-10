package svc

import (
	"github.com/krace-tx/emo_trash/app/rpc/sso/internal/config"
	"github.com/krace-tx/emo_trash/pkg/db/rdb"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type ServiceContext struct {
	Config config.Config
	Redis  *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	if err := rdb.InitMySQL(c.MysqlConf); err != nil {
		panic(err)
	}

	return &ServiceContext{
		Config: c,
		Redis:  redis.MustNewRedis(c.ZrpcConf.Redis.RedisConf),
	}
}
