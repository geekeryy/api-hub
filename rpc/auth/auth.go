package main

import (
	"flag"
	"fmt"

	"github.com/geekeryy/api-hub/rpc/auth/auth"
	"github.com/geekeryy/api-hub/rpc/auth/internal/config"
	authserviceServer "github.com/geekeryy/api-hub/rpc/auth/internal/server/authservice"
	"github.com/geekeryy/api-hub/rpc/auth/internal/svc"

	coreerror "github.com/geekeryy/api-hub/core/error"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/auth.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		auth.RegisterAuthServiceServer(grpcServer, authserviceServer.NewAuthServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	s.AddUnaryInterceptors(coreerror.ErrorInterceptor)
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
