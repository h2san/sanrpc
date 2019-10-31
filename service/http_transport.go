package service

import (
	log "github.com/hillguo/sanlog"
	"github.com/hillguo/sanrpc/protocol"
	"net/http"
)

type httpTransport struct {
	s *service
}

func NewHTTPTransport(s *service) ServerTransport {
	return &httpTransport{s}
}

func (t *httpTransport) ListenAndServer() error{

	log.Infof("http listening network:%s ,address:%s", t.s.opts.NetWork, t.s.opts.Address)
	err := http.ListenAndServe(t.s.opts.Address, t)
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

	if app,ok := t.s.opts.MsgProtocol.(protocol.HttpMsgProtocol); ok {
		app.ServeHTTP(w,req)
	}else {
		log.Error("sanrpc: msg protocol is not http protocol")
	}

}

