package main

import (
	"fmt"

	"github.com/hillguo/sanrpc/tool/gencode"
)

var formatServerPart1 = `package main

import (
	"context"
	"fmt"
	
	"%s/pb"
)

type %s struct {

}
`
var formatServerFunc = `
func (impl *%s) %s (ctx context.Context, req *pb.%s, resp *pb.%s) error {
	fmt.Println("unimplemented")
	return nil
} 
`

func genServer(protoInfo *gencode.ProtoFileInfo) (string, string) {
	data := fmt.Sprintf(formatServerPart1, protoInfo.ModuleName, protoInfo.ServiceName)
	for _, methodInfo := range protoInfo.Methods {
		data += fmt.Sprintf(formatServerFunc, protoInfo.ServiceName, methodInfo.MethodName, methodInfo.InputType, methodInfo.OutputType)
	}
	return protoInfo.ModuleName + "serviceimpl.go", data
}

func main() {
	gencode.Main(genServer)
}
