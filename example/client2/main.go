package main

import (
	"context"
	log "github.com/hillguo/sanlog"
	"github.com/hillguo/sanrpc/client"
	"github.com/hillguo/sanrpc/example"
)

func main()  {

	opts := []client.Option{
		client.WithAddress("127.0.0.1:8001"),
		client.WithServiceName("test"),
		client.WithMethodName("add"),
	}
	c := client.NewClient(opts...)

	a :=&example.Req{
		A:100,
	}
	b :=&example.Resq{}

	err := c.Invoke(context.Background(),a, b)
	if err != nil {
		log.Debug(err)
	}
	log.Debug(*a,*b)
}