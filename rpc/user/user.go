package main

import (
	"flag"
	"fmt"

	"github.com/geekeryy/api-hub/rpc/user/internal/config"
	adminserviceServer "github.com/geekeryy/api-hub/rpc/user/internal/server/adminservice"
	memberserviceServer "github.com/geekeryy/api-hub/rpc/user/internal/server/memberservice"
	"github.com/geekeryy/api-hub/rpc/user/internal/svc"
	"github.com/geekeryy/api-hub/rpc/user/user"

	coreerror "github.com/geekeryy/api-hub/core/error"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/user.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		user.RegisterAdminServiceServer(grpcServer, adminserviceServer.NewAdminServiceServer(ctx))
		user.RegisterMemberServiceServer(grpcServer, memberserviceServer.NewMemberServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	s.AddUnaryInterceptors(coreerror.ErrorInterceptor)
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
