package config

import (
	authx "github.com/krace-tx/emo_trash/pkg/auth"
	"github.com/krace-tx/emo_trash/pkg/db/no_sql"
	"github.com/krace-tx/emo_trash/pkg/db/rdb"
	"github.com/krace-tx/emo_trash/pkg/snowflake"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	ServiceConf   service.ServiceConf     `json:"Service"`
	ZrpcConf      zrpc.RpcServerConf      `json:"Zrpc"`
	MysqlConf     rdb.DBConf              `json:"Mysql"`
	JWT           authx.JWTConf           `json:"JWT"`
	MongoConf     no_sql.MongoConf        `json:"Mongo"`
	SnowflakeConf snowflake.SnowflakeConf `json:"Snowflake"`
}
