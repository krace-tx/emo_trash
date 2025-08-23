package svc

import (
	"github.com/krace-tx/emo_trash/app/rpc/sso/internal/config"
	"github.com/krace-tx/emo_trash/pkg/db/no_sql" // 新增：导入no_sql包
	"github.com/krace-tx/emo_trash/pkg/db/rdb"
	"github.com/krace-tx/emo_trash/pkg/snowflake"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

type ServiceContext struct {
	Config    config.Config
	Redis     *redis.Redis
	Mongo     *mongo.Client
	Snowflake *snowflake.Snowflake
}

func NewServiceContext(c config.Config) *ServiceContext {
	rdb.MustInitMySQL(c.MysqlConf)

	return &ServiceContext{
		Config:    c,
		Redis:     redis.MustNewRedis(c.ZrpcConf.Redis.RedisConf),
		Mongo:     no_sql.MustInitMongo(c.MongoConf),
		Snowflake: snowflake.MustNewSnowflake(c.SnowflakeConf),
	}
}
