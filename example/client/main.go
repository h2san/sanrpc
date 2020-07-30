package main

import (
	"context"
	log "github.com/hillguo/sanlog"
	"github.com/hillguo/sanrpc/client1"
	"github.com/hillguo/sanrpc/example"
)

func main()  {

	opt := client1.Option{

	}
	c := client1.NewClient(opt)
	err := c.Connect("tcp","127.0.0.1:8000")
	if err != nil {
		log.Fatal(err)
	}
	a :=&example.Req{
		A:100,
	}
	b :=&example.Resq{}

	err = c.Call(context.Background(),"test","add",a, b)
	if err != nil {
		log.Debug(err)
	}
	log.Debug(*a,*b)
}