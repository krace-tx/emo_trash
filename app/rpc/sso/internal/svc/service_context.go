package svc

import (
	"github.com/krace-tx/emo_trash/app/rpc/sso/internal/config"
	"github.com/krace-tx/emo_trash/pkg/datastore/mongo"
	"github.com/krace-tx/emo_trash/pkg/datastore/redis"
	"github.com/krace-tx/emo_trash/pkg/datastore/sqlstore"
	"github.com/krace-tx/emo_trash/pkg/snowflake"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type ServiceContext struct {
	Config    config.Config
	Redis     *redis.Redis
	Mongo     *mongodriver.Database
	Snowflake *snowflake.Snowflake
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlstore.MustInitMySQL(c.MysqlConf)

	return &ServiceContext{
		Config:    c,
		Redis:     redis.NewRedisClient(c.RedisConf),
		Mongo:     mongo.MustInitMongo(c.MongoConf),
		Snowflake: snowflake.MustNewSnowflake(c.SnowflakeConf),
	}
}
