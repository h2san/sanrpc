package protocol

import (
	"context"
	"io"
	"net/http"
)

//Message 完整的包体
type Message interface{}

// protocol
type MsgProtocol interface {

}


//RpcMsgProtocol rpc protocol inteface
type RpcMsgProtocol interface {
	DecodeMessage(r io.Reader) (Message, error)
	HandleMessage(ctx context.Context, req Message) (resp Message, err error)
	EncodeMessage(res Message) ([]byte, error)
}

//HttpMsgProtocol http protocol interface
type HttpMsgProtocol interface {
	ServeHTTP(w http.ResponseWriter,req *http.Request)
}

type RegisterServicer interface {
	RegisterService(rcvr interface{}) error
}