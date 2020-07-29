package service

import (
	"errors"
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

	app,ok := t.s.opts.MsgProtocol.(protocol.HttpMsgProtocol)
	if !ok {
		log.Error("sanrpc: msg protocol is not http protocol")
		return errors.New("sanrpc: msg protocol is not http protocol")
	}

	log.Infof("http listening network:%s ,address:%s", t.s.opts.NetWork, t.s.opts.Address)
	err := http.ListenAndServe(t.s.opts.Address, app)
	if err != nil {
		log.Errorf("ListenAndServe fail:", err)
		return err
	}
	return  nil
}

func (t *httpTransport) Close() error{
	return nil
}

