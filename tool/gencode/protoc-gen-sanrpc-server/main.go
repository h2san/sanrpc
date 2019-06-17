package main

import (
	"fmt"

	"github.com/hillguo/sanrpc/tool/gencode"
)

var formatServerPart1 = `package %s

import (
	"context"
	"github.com/hillguo/sanrpc/server"
)

type %sImpl struct {

}
`
var formatServerFunc = `
func (impl *%sImpl)(ctx *context.Context, req *%s, resp *%s) error {
	panic("unimplemented")
} 
`
var formatServerPart2 = `
func main() {
	
}
`

func genServer(protoInfo *gencode.ProtoFileInfo) (string, string) {
	data := fmt.Sprintf(formatServerPart1, protoInfo.PackageName, protoInfo.ModuleName)
	for _, methodInfo := range protoInfo.Methods {
		data += fmt.Sprintf(formatServerFunc, methodInfo.MethodName, methodInfo.InputType, methodInfo.OutputType)
	}
	data += formatServerPart2
	return protoInfo.ModuleName + "server.go", data
}

func main() {
	gencode.Main(genServer)
}
