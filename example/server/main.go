package main

import (
	"context"
	"github.com/hillguo/sanrpc"
	"github.com/hillguo/sanrpc/example"
	"github.com/hillguo/sanrpc/config"
	"github.com/hillguo/sanrpc/service"
	log "github.com/hillguo/sanlog"
	"reflect"
)

type Test struct {

}


func (t *Test) Del(ctx context.Context, a *example.Req , b *example.Resq) error  {
	log.Info("11")

	return nil
}

type Test2 struct {

}

func (t *Test2) Add(ctx context.Context, a *example.Req , b *example.Resq) error  {
	log.Info("11")

	return nil
}

type ATest struct {
	Test
	Test2
}

func (t *ATest) CAdd(ctx context.Context, a *example.Req , b *example.Resq) error  {
	log.Info("CAdd")
	return nil
}

func main(){
	t := reflect.TypeOf(&ATest{})
	log.Info(t.NumMethod())

	svr := sanrpc.NewServer()
	c := config.GetConfig()
	log.Info(c)
	ss := service.NewServicesWithConfig(&c.Server)
	log.Info(ss)
	for _,s := range ss {
		if s.Name() == "test" {
			err := s.Register(&ATest{})
			log.Info(err)
		}
		svr.AddService(s.Name(), s)
	}
	svr.Serve()
}
