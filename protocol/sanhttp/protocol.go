package sanhttp

import (
	"context"
	"errors"
	"github.com/hillguo/sanrpc/protocol"
	"io"
	"net/http"
)

type SanHTTPProtocol struct {
	protocol.BaseService
}

func(p *SanHTTPProtocol) DecodeMessage(r io.Reader) (protocol.Message, error){
	return nil, nil
}
func(p *SanHTTPProtocol) HandleMessage(ctx context.Context, msg protocol.Message) ( protocol.Message,  error){

	m,ok := msg.(Msg);
	if !ok {
		return nil, errors.New("msg is not http message")
	}

	r := m.Req
	w := m.Resp

	req := &ReqMsg{}
	resp := &RespMsg{}


	err := req.Decode(r.Body)
	if err != nil {
		resp.Code = 10000
		resp.Msg = "xxxxx"
		err = resp.Write(w)
	}

	if err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return nil, nil
	}







	return nil, nil
}
func(p *SanHTTPProtocol) EncodeMessage(res protocol.Message) ([]byte, error){
	return nil, nil
}