package svc

import (
	"github.com/krace-tx/emo_trash/app/rpc/sso/internal/config"
	"github.com/krace-tx/emo_trash/pkg/db/no_sql"
	"github.com/krace-tx/emo_trash/pkg/db/rdb"
	"github.com/krace-tx/emo_trash/pkg/snowflake"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type ServiceContext struct {
	Config    config.Config
	Redis     *no_sql.Redis
	Mongo     *mongo.Client
	Snowflake *snowflake.Snowflake
}

func NewServiceContext(c config.Config) *ServiceContext {
	rdb.MustInitMySQL(c.MysqlConf)

	redis.NewClient(&redis.Options{})

	return &ServiceContext{
		Config:    c,
		Redis:     no_sql.NewRedisClient(c.RedisConf),
		Mongo:     no_sql.MustInitMongo(c.MongoConf),
		Snowflake: snowflake.MustNewSnowflake(c.SnowflakeConf),
	}
}
