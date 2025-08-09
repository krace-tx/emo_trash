package main

import (
	"flag"
	"fmt"

	"github.com/krace-tx/emo_trash/app/rpc/sso/internal/config"
	authServer "github.com/krace-tx/emo_trash/app/rpc/sso/internal/server/auth"
	"github.com/krace-tx/emo_trash/app/rpc/sso/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/sso/pb"
	"github.com/krace-tx/emo_trash/pkg/interceptor"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/sso.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.ZrpcConf, func(grpcServer *grpc.Server) {
		pb.RegisterAuthServer(grpcServer, authServer.NewAuthServer(ctx))

		if c.ZrpcConf.Mode == service.DevMode || c.ZrpcConf.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})

	// 添加一元拦截器
	s.AddUnaryInterceptors(interceptor.UnaryRecoverInterceptor)
	interceptor.SignalInterceptor()

	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ZrpcConf.ListenOn)
	s.Start()
}
