package main

import (
	"context"
	log "github.com/hillguo/sanlog"
	"github.com/hillguo/sanrpc/service"
)

type Test struct {

}

func (t *Test) Add(ctx context.Context, a int , b *int) error {
	log.Debug("11")
	return nil
}

func main(){
	svr := service.New(service.WithNetWork("tcp"),service.WithAddress("127.0.0.1:9999"))
	svr.Register(&Test{})
	svr.Serve()
}
