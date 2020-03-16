package main

import (
	log "github.com/hillguo/sanlog"
	"github.com/hillguo/sanrpc"
	"github.com/hillguo/sanrpc/example"
	"github.com/hillguo/sanrpc/protocol/sanhttp"
	"github.com/hillguo/sanrpc/protocol/sanhttp/ctx"
	"github.com/hillguo/sanrpc/service"
)

func CAdd(ctx *ctx.Context, a *example.Req , b *example.Resq) error  {
	log.Info("CAdd")
	return nil
}

func main()  {
	svr := sanrpc.NewServer()
	app := sanhttp.Default()
	app.Any("/hello", sanhttp.HF(CAdd))
	ss := service.New(service.WithServiceName("http"),service.WithAddress("127.0.0.1:8081"),service.WithNetWork("http"),service.WithMsgProtocol(app))

	svr.AddService("http", ss)
	svr.Serve()
}
