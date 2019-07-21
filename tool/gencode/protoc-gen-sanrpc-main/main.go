package main

import (
	"fmt"

	"github.com/hillguo/sanrpc/tool/gencode"
)

var main_tpl = `package main

import (
	"github.com/hillguo/sanrpc/server"
)

func main() {
	s := server.NewRpcServer()
	s.RegisterService(&%s{})
	// s.Serve("http", "0.0.0.0:8080")
	s.Serve("tcp", "0.0.0.0:8080")
}
`

func genMain(protoInfo *gencode.ProtoFileInfo) (string, string) {
	data := fmt.Sprintf(main_tpl, protoInfo.ServiceName)
	return "main.go", data
}

func main() {
	gencode.Main(genMain)
}
