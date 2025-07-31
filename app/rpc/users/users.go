package main

import (
	"flag"
	"fmt"

	"github.com/krace-tx/emo_trash/app/rpc/users/internal/config"
	socialServer "github.com/krace-tx/emo_trash/app/rpc/users/internal/server/social"
	usersServer "github.com/krace-tx/emo_trash/app/rpc/users/internal/server/users"
	"github.com/krace-tx/emo_trash/app/rpc/users/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/users/users"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/users.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		users.RegisterUsersServer(grpcServer, usersServer.NewUsersServer(ctx))
		users.RegisterSocialServer(grpcServer, socialServer.NewSocialServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
