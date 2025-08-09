package config

import (
	"github.com/krace-tx/emo_trash/pkg/db/rdb"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	ServiceConf service.ServiceConf `json:"service"`
	ZrpcConf    zrpc.RpcServerConf  `json:"zrpc"`
	MysqlConf   rdb.DBConfig        `json:"mysql"`
}
