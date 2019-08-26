package transport

import (
	log "github.com/hillguo/sanlog"
	"github.com/hillguo/sanrpc/protocol/sanhttp"
	"net/http"
)

type httpTransport struct {
	opts *TransportOptions

}

func NewHTTPTransport(opt ...TransportOption) ServerTransport{
	opts := &TransportOptions{

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
		t.opts.MsgProtocol = &sanhttp.SanHTTPProtocol{}
	}
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
	if req.URL.Path != "/SanHTTP" {
		http.NotFound(w,req)
		return
	}
	msg := sanhttp.Msg{
		Req:req,
		Resp:w,
	}
	_, err := t.opts.MsgProtocol.HandleMessage(req.Context(),msg)
	if err != nil {
		log.Error("msg protocol handle message fail", err)
	}
}

