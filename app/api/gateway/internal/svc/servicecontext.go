package svc

import (
	"github.com/krace-tx/emo_trash/app/api/gateway/internal/config"
	"github.com/krace-tx/emo_trash/app/api/gateway/internal/middleware"
	"github.com/krace-tx/emo_trash/app/rpc/post/client/post"
	"github.com/krace-tx/emo_trash/app/rpc/sso/client/auth"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config
	Filter rest.Middleware
	Log    rest.Middleware
	Auth   rest.Middleware
	Sso    auth.Auth
	Post   post.Post
}

func NewServiceContext(c config.Config) *ServiceContext {
	authRpc := auth.NewAuth(zrpc.MustNewClient(c.Rpc.Auth))
	return &ServiceContext{
		Config: c,
		Filter: middleware.NewFilterMiddleware().Handle,
		Log:    middleware.NewLogMiddleware().Handle,
		Auth:   middleware.NewAuthMiddleware(authRpc).Handle,
		Sso:    authRpc,
		Post:   post.NewPost(zrpc.MustNewClient(c.Rpc.Post)),
	}
}
