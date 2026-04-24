package config

import (
	authx "github.com/krace-tx/emo_trash/pkg/auth"
	"github.com/krace-tx/emo_trash/pkg/datastore/mongo"
	"github.com/krace-tx/emo_trash/pkg/datastore/redis"
	"github.com/krace-tx/emo_trash/pkg/datastore/sqlstore"
	"github.com/krace-tx/emo_trash/pkg/email"
	"github.com/krace-tx/emo_trash/pkg/snowflake"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	ServiceConf   service.ServiceConf     `json:"Service"`
	ZrpcConf      zrpc.RpcServerConf      `json:"Zrpc"`
	MysqlConf     sqlstore.DBConf         `json:"Mysql"`
	JWT           authx.JWTConf           `json:"JWT"`
	Email         email.Pop3              `json:"Email"`
	MongoConf     mongo.MongoConf         `json:"Mongo"`
	RedisConf     redis.RedisConf         `json:"Redis"`
	SnowflakeConf snowflake.SnowflakeConf `json:"Snowflake"`
}
