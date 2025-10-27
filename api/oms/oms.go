package main

import (
	"flag"
	"fmt"

	"github.com/geekeryy/api-hub/api/oms/internal/config"
	"github.com/geekeryy/api-hub/api/oms/internal/handler"
	"github.com/geekeryy/api-hub/api/oms/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/proc"
	"github.com/zeromicro/go-zero/rest"

	coreHandler "github.com/geekeryy/api-hub/core/handler"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var configFile = flag.String("f", "etc/oms-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	httpx.SetErrorHandlerCtx(coreHandler.ErrorHandler)
	httpx.SetOkHandler(coreHandler.OkHandler)

	proc.AddWrapUpListener(func() {
		ctx.Close()
	})

	logx.DisableStat()

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
