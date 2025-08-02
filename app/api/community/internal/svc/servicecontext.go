package svc

import (
	"github.com/krace-tx/emo_trash/app/api/community/internal/config"
	"github.com/krace-tx/emo_trash/app/api/community/internal/middleware"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config config.Config
	Filter rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Filter: middleware.NewFilterMiddleware().Handle,
	}
}
