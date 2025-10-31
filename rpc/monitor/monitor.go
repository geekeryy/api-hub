package main

import (
	"flag"
	"fmt"

	"github.com/geekeryy/api-hub/rpc/monitor/internal/config"
	monitorserviceServer "github.com/geekeryy/api-hub/rpc/monitor/internal/server/monitorservice"
	"github.com/geekeryy/api-hub/rpc/monitor/internal/svc"
	"github.com/geekeryy/api-hub/rpc/monitor/monitor"

	coreerror "github.com/geekeryy/api-hub/core/error"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/monitor.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		monitor.RegisterMonitorServiceServer(grpcServer, monitorserviceServer.NewMonitorServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	s.AddUnaryInterceptors(coreerror.ErrorInterceptor)
	defer s.Stop()

	logx.DisableStat()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
