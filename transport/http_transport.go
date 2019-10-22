package transport

import (
	log "github.com/hillguo/sanlog"
	"github.com/hillguo/sanrpc/protocol"
	"github.com/hillguo/sanrpc/protocol/httpx"
	"net/http"
)

type httpTransport struct {
	opts *TransportOptions

}

func NewHTTPTransport(opt ...TransportOption) ServerTransport{
	opts := &TransportOptions{
		MsgProtocol:    &httpx.HTTProtocol{},
	}
	for _, o := range opt {
		o(opts)
	}
	return &httpTransport{opts:opts}
}

func (t *httpTransport) ListenAndServer(opt ...TransportOption) error{
	for _, o := range opt {
		o(t.opts)
	}

	if t.opts.MsgProtocol == nil {
		t.opts.MsgProtocol = httpx.DefaultHTTProtocol
	}


	log.Infof("http listening network:%s ,address:%s", t.opts.NetWork, t.opts.Address)
	err := http.ListenAndServe(t.opts.Address, t)
	if err != nil {
		log.Errorf("ListenAndServe fail:", err)
		return err
	}
	return  nil
}

func (t *httpTransport) Close() error{
	return nil
}

func (t *httpTransport) ServeHTTP (w http.ResponseWriter, req *http.Request) {

	if app,ok := t.opts.MsgProtocol.(protocol.HttpMsgProtocol); ok {
		app.ServeHTTP(w,req)
	}else {
		log.Error("sanrpc: msg protocol is not http protocol")
	}

}

