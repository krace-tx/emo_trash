package config

import (
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	ServiceConf service.ServiceConf `json:"service"`
	ZrpcConf    zrpc.RpcServerConf  `json:"zrpc"`
}
