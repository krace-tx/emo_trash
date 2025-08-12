package svc

import (
	"github.com/krace-tx/emo_trash/app/rpc/sso/internal/config"
	"github.com/krace-tx/emo_trash/pkg/db/no_sql" // 新增：导入no_sql包
	"github.com/krace-tx/emo_trash/pkg/db/rdb"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

type ServiceContext struct {
	Config config.Config
	Redis  *redis.Redis
	Mongo  *mongo.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	if err := rdb.InitMySQL(c.MysqlConf); err != nil {
		panic(err)
	}

	mongo, err := no_sql.InitMongo(c.MongoConf)
	if err != nil {
		panic(err)
	}

	redis := redis.MustNewRedis(c.ZrpcConf.Redis.RedisConf)

	return &ServiceContext{
		Config: c,
		Redis:  redis,
		Mongo:  mongo,
	}
}
