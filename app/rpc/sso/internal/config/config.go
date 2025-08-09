package config

import (
	authx "github.com/krace-tx/emo_trash/pkg/auth"
	"github.com/krace-tx/emo_trash/pkg/db/rdb"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	ServiceConf service.ServiceConf `json:"Service"`
	ZrpcConf    zrpc.RpcServerConf  `json:"Zrpc"`
	MysqlConf   rdb.DBConfig        `json:"Mysql"`
	JWT         authx.JWTConfig     `json:"JWT"`
}
