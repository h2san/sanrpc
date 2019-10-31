package main

import (
	"context"
	"github.com/hillguo/sanrpc"
	"github.com/hillguo/sanrpc/example"
	"github.com/hillguo/sanrpc/config"
	"github.com/hillguo/sanrpc/service"
	log "github.com/hillguo/sanlog"
)

type Test struct {

}

func (t *Test) Add(ctx context.Context, a *example.Req , b *example.Resq) error  {
	log.Info("11")

	return nil
}

func main(){
	svr := sanrpc.NewServer()
	c := config.GetConfig()
	log.Info(c)
	ss := service.NewServicesWithConfig(&c.Server)
	log.Info(ss)
	for _,s := range ss {
		if s.Name() == "test" {
			s.Register(&Test{})
		}
		svr.AddService(s.Name(), s)
	}
	svr.Serve()
}
