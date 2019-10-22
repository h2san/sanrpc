package main

import (
	"context"
	log "github.com/hillguo/sanlog"
	"github.com/hillguo/sanrpc"
	"github.com/hillguo/sanrpc/example"
	"github.com/hillguo/sanrpc/service"
	"github.com/hillguo/sanrpc/transport"
)

type Test struct {

}

func (t *Test) Add(ctx context.Context, a *example.Req , b *example.Resq) error  {
	log.Debug("11")

	return nil
}

func main(){

	svr1 := service.New(service.WithNetWork("tcp"),service.WithAddress("127.0.0.1:9999"), service.WithTransport(transport.NewHTTPTransport()))
	svr1.Register(&Test{})

	svr2 := service.New(service.WithNetWork("tcp"),service.WithAddress("127.0.0.1:9998"), service.WithTransport(transport.NewTCPTransport()))
	svr2.Register(&Test{})

	svr := sanrpc.NewServer()
	svr.AddService("rpc", svr1)
	svr.AddService("http", svr2)
	svr.Serve()
}
