package main

import (
	"flag"
	"fmt"

	{{.importPackages}}

	coreHandler "{{.projectPkg}}/core/handler"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var configFile = flag.String("f", "etc/{{.serviceName}}.yaml", "the config file")

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
