package gencode

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	pluginGo "github.com/golang/protobuf/protoc-gen-go/plugin"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type MethodInfo struct {
	MethodName string
	InputType  string
	OutputType string
}

type ProtoFileInfo struct {
	FullProtoName string // ./mmhelloworld/mmhelloworld.proto
	GoPackageName string // git.code.oa.com/wxg-contrib/svrkit-go/mmhelloworld
	PackageName   string // mmhelloworld
	ProtoDir      string // ./mmhelloworld
	ServiceName   string // MMHelloWorld
	ModuleName    string // mmhelloworld
	Methods       []*MethodInfo
}

func NewProtoFileInfo(fd *descriptor.FileDescriptorProto) *ProtoFileInfo {
	info := &ProtoFileInfo{}
	info.FullProtoName = fd.GetName()
	info.ProtoDir = filepath.Dir(info.FullProtoName)
	info.GoPackageName = fd.GetOptions().GetGoPackage()
	if len(info.GoPackageName) == 0 {
		info.PackageName = fd.GetPackage()
	} else {
		info.PackageName = filepath.Base(info.GoPackageName)
	}
	if len(fd.Service) != 1 {
		log.Fatalf("invalid service num%d in proto:%s", len(fd.Service), info.FullProtoName)
	}
	sd := fd.Service[0]
	info.ServiceName = sd.GetName()
	info.ModuleName = strings.ToLower(info.ServiceName)
	info.Methods = make([]*MethodInfo, len(sd.Method))
	for i, md := range sd.Method {
		methodInfo := &MethodInfo{}
		methodInfo.MethodName = md.GetName()
		// TODO 需要适配input或output import自不同package的proto
		methodInfo.InputType = filepath.Ext(md.GetInputType())[1:]
		methodInfo.OutputType = filepath.Ext(md.GetOutputType())[1:]
		info.Methods[i] = methodInfo
	}
	return info
}

type GenOenProtoTypeFunc func(*ProtoFileInfo) (string, string)

func Main(genOenProtoTypeFunc ...GenOenProtoTypeFunc) {
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("fail ReadAll input,err:%v", err)
	}
	req := &pluginGo.CodeGeneratorRequest{}
	if err := proto.Unmarshal(data, req); err != nil {
		log.Fatalf("fail Unmarshal input proto,err:%v", err)
	}
	if len(req.FileToGenerate) == 0 {
		log.Fatal("no files to generate")
	}
	needGenMap := make(map[string]byte)
	for _, f := range req.FileToGenerate {
		needGenMap[f] = 0
	}
	resp := &pluginGo.CodeGeneratorResponse{}
	for _, fd := range req.ProtoFile {
		if _, ok := needGenMap[fd.GetName()]; !ok {
			continue
		}
		protoInfo := NewProtoFileInfo(fd)
		for _, gen := range genOenProtoTypeFunc {
			fileName, fileData := gen(protoInfo)
			if protoInfo.ProtoDir != "." {
				fileName = protoInfo.ProtoDir + "/" + fileName
			}
			if _, err := os.Stat(fileName); err == nil {
				log.Printf("ignore existed file:%s", fileName)
			} else if os.IsNotExist(err) {
				log.Printf("add new file:%s", fileName)
				newFile := &pluginGo.CodeGeneratorResponse_File{}
				newFile.Name = &fileName
				newFile.Content = &fileData
				resp.File = append(resp.File, newFile)
			} else {
				log.Fatalf("fail Stat,err:%v", err)
			}
		}
	}
	if data, err = proto.Marshal(resp); err != nil {
		log.Fatalf("fail Marshal output proto,err:%v", err)
	}
	if _, err = os.Stdout.Write(data); err != nil {
		log.Fatalf("fail Write output proto,err:%v", err)
	}
}
