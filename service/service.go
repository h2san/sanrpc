package service

import (
	log "github.com/hillguo/sanlog"
	"github.com/hillguo/sanrpc/errs"
	"github.com/hillguo/sanrpc/protocol"
	"github.com/hillguo/sanrpc/protocol/sanrpc"
	"github.com/hillguo/sanrpc/transport"
)

type Service interface {
	Serve() error
	Register(serviceDesc interface{}) error
}

type service struct {
	opts *Options
}

func New(opts ...Option) Service {
	s := &service{
		opts: &Options{
			ServeTransport: transport.DefaultTransport,
			MsgProtocol:&sanrpc.SanRPCProtocol{},
		},
	}
	for _, o := range opts {
		o(s.opts)
	}
	return s
}

func (s *service) Serve() error {
	s.opts.ServeTransportOptions = append(s.opts.ServeTransportOptions, transport.WithMsgProtocol(s.opts.MsgProtocol))
	err := s.opts.ServeTransport.ListenAndServer(s.opts.ServeTransportOptions...)
	if err != nil {
		log.Error(err)
	}
	return err
}

func (s *service) Register(serviceDesc interface{}) error{
	p := s.opts.MsgProtocol
	if p == nil {
		return errs.ErrServerNoMsgProtocol
	}
	if rs, ok := p.(protocol.RegisterServicer); ok {
		return  rs.RegisterService(serviceDesc)
	}
	return nil
}
