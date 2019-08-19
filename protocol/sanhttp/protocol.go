package sanhttp

import (
	"context"
	"github.com/hillguo/sanrpc/protocol"
	"io"
)

type SanHTTPProtocol struct {
	protocol.BaseService
}

func(p *SanHTTPProtocol) DecodeMessage(r io.Reader) (protocol.Message, error){
	return nil, nil
}
func(p *SanHTTPProtocol) HandleMessage(ctx context.Context, req protocol.Message) (resp protocol.Message, err error){

	return nil, nil
}
func(p *SanHTTPProtocol) EncodeMessage(res protocol.Message) ([]byte, error){
	return nil, nil
}