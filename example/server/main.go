package main

import (
	"context"
	log "github.com/hillguo/sanlog"
	"github.com/hillguo/sanrpc"
	"github.com/hillguo/sanrpc/example"
	"github.com/hillguo/sanrpc/service"
)

type Test struct {

}


func (t *Test) Del(ctx context.Context, a *example.Req , b *example.Resq) error  {
	log.Info("11")

	return nil
}

type Test2 struct {

}

func (t *Test) Add(ctx context.Context, a *example.Req , b *example.Resq) error  {
	log.Info("11")
	b.B = a.A + 1
	return nil
}

type ATest struct {
	Test
	Test2
}

func (t *ATest) CAdd(ctx context.Context, a *example.Req , b *example.Resq) error  {
	log.Info("CAdd")
	b.B = a.A + 1
	return nil
}

func main(){

	svr := sanrpc.NewServer()

	ss := service.New(service.WithServiceName("rpc"), service.WithAddress("127.0.0.1:8000"))
	svr.AddService("rpc",ss)
	svr.Register(&Test{})
	svr.Serve()
}
