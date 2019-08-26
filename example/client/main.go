package main

import (
	"context"
	log "github.com/hillguo/sanlog"
	"github.com/hillguo/sanrpc/client"
	"github.com/hillguo/sanrpc/example"
)

func main()  {

	opt := client.Option{

	}
	c := client.NewClient(opt)
	err := c.Connect("tcp","127.0.0.1:9999")
	if err != nil {
		log.Fatal(err)
	}
	a :=&example.Req{}
	b :=&example.Resq{}

	err = c.Call(context.Background(),"test","add",a, b)
	if err != nil {
		log.Debug(err)
	}
	log.Debug(a.String(), *b)
}