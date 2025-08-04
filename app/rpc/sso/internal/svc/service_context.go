package svc

import (
	"github.com/krace-tx/emo_trash/app/rpc/sso/internal/config"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type ServiceContext struct {
	Config        config.Config
	Rds           *redis.Redis
	KafkaProducer *kq.Pusher
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:        c,
		Rds:           redis.MustNewRedis(c.Redis.RedisConf),
		KafkaProducer: kq.NewPusher(c.KqPusherConf.Brokers, c.KqPusherConf.Topic),
	}
}
