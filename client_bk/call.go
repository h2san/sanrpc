package client1

import (
	log "github.com/hillguo/sanlog"
	"github.com/hillguo/sanrpc/codec"
)

type Call struct {
	ServicePath   string
	ServiceMethod string
	Metadata      map[string]string
	ResMetadata   map[string]string
	Args          interface{}
	Reply         interface{}
	Error         error
	CallDoneChan          chan CallDoneSignal
	Raw           bool

	// 每个调用都可以设置编码与压缩类型
	SerializeType codec.SerializeType
	CompressType  codec.CompressType
}

type CallDoneSignal struct {

}
var CallDone CallDoneSignal

func (call *Call) done() {
	select {
	case call.CallDoneChan <- CallDone:
	default:
		log.Info("rpc: discarding Call reply due to insufficient Done chan capacity")

	}
}
