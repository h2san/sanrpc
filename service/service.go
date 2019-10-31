package service

import (
	log "github.com/hillguo/sanlog"
	"github.com/hillguo/sanrpc/config"
	"github.com/hillguo/sanrpc/errs"
	"github.com/hillguo/sanrpc/protocol"
	"github.com/hillguo/sanrpc/protocol/httpx"
	"github.com/hillguo/sanrpc/protocol/sanrpc"
)

type Service interface {
	Name()string
	Serve() error
	Register(serviceDesc interface{}) error
}

type service struct {
	opts *Options
	ServeTransport ServerTransport
}

func NewServicesWithConfig(server_config *config.ServerConfig) []Service {
	if server_config == nil {
		panic("server_config nil")
		return nil
	}
	log.Info(server_config)
	ss := make([]Service,0,len(server_config.Services))
	log.Info(len(ss))
	for _, svr :=range server_config.Services {
		log.Info(svr)
		s := &service{
			opts: &Options{
				Name:           svr.Name,
				InMsgChanSize:  svr.InMsgChanSize,
				OutMsgChanSize: svr.OutMsgChanSize,
				ReadTimeout:    svr.ReadTimeout,
				WriteTimeout:   svr.OutMsgChanSize,
				Address:        svr.Address,
				NetWork:        svr.NetWork,
				TLSCertFile:    svr.TLSCertFile,
				TLSKeyFile:     svr.TLSKeyFile,
			},
		}
		if svr.Protocol == "http" {
			s.ServeTransport = NewHTTPTransport(s)
			s.opts.MsgProtocol = &httpx.HTTProtocol{}
		} else {
			s.ServeTransport = NewTCPTransport(s)
			s.opts.MsgProtocol = &sanrpc.SanRPCProtocol{}
		}
		ss= append(ss,s)
	}
	log.Infof("%v",len(ss))
	return  ss
}

func New(opts ...Option) Service {
	var s *service
	s = &service{
		opts: &Options{
			Name:           "",
			InMsgChanSize:  1024,
			OutMsgChanSize: 1024,
			ReadTimeout:    0,
			WriteTimeout:   0,
			Address:        "",
			NetWork:        "tcp",
			TLSCertFile:    "",
			TLSKeyFile:     "",
			MsgProtocol:    &sanrpc.SanRPCProtocol{},
		},
	}
	for _, o := range opts {
		o(s.opts)
	}
	if s.opts.NetWork == "tcp" {
		s.ServeTransport =  NewTCPTransport(s)
	}else if s.opts.NetWork == "http" {
		s.ServeTransport =  NewHTTPTransport(s)
	}
	return s
}

func (s *service) Name() string{
	return s.opts.Name
}
func (s *service) Serve() error {
	err := s.ServeTransport.ListenAndServer()
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
