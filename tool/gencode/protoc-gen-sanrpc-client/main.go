package main

import (
	"fmt"

	"github.com/hillguo/sanrpc/tool/gencode"
)

var formatPart1 = `package client1

import (
	"context"

	"github.com/hillguo/sanrpc/client1/selector"
	"github.com/hillguo/sanrpc/client1/servicediscovery"
	 xclient "github.com/hillguo/sanrpc/client1"
	pb "%s/pb"
)

type %sClient struct {
	xclient.XClient
}

func New%sClient() *%sClient {
	c := NoteSvrClient{}
	c.XClient= xclient.NewXClient(
		"%s",
		xclient.Failtry,
		selector.RoundRobin,
		servicediscovery.NewP2PDiscovery("127.0.0.1:8080",""),
		xclient.DefaultOption)
	return &c
}
`

var formatFunc = `
func (c *%sClient) %s(ctx context.Context, req *pb.%s, resp *pb.%s) error {
	return c.Call(ctx,  "%s", req, resp)
}
`

func genClient(protoInfo *gencode.ProtoFileInfo) (string, string) {
	data := fmt.Sprintf(formatPart1, protoInfo.ModuleName, protoInfo.ServiceName, protoInfo.ServiceName,
		protoInfo.ServiceName, protoInfo.ServiceName)

	for _, methodInfo := range protoInfo.Methods {
		data += fmt.Sprintf(formatFunc, protoInfo.ServiceName, methodInfo.MethodName,
			methodInfo.InputType, methodInfo.OutputType, methodInfo.MethodName)
	}
	return protoInfo.ModuleName + "client1.go", data
}

func main() {
	gencode.Main(genClient)
}
