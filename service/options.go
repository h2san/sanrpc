package service

import (
	"github.com/hillguo/sanrpc/protocol"
	"github.com/hillguo/sanrpc/transport"
)

type Options struct {
	Name string
	ServeTransport transport.ServerTransport
	ServeTransportOptions []transport.TransportOption
	MsgProtocol protocol.MsgProtocol
}

type Option func(options *Options)

func WithServiceName(s string) Option{
	return func(o *Options){
		o.Name = s
	}
}

func WithAddress(s string) Option{
	return func(o *Options) {
		o.ServeTransportOptions = append(o.ServeTransportOptions, transport.WithAddress(s))
	}
}

func WithNetWork(network string) Option{
	return func(o *Options) {
		o.ServeTransportOptions = append(o.ServeTransportOptions, transport.WithNetWork(network))
	}
}

func WithReadTimeout(timeout uint32) Option{
	return func(o *Options) {
		o.ServeTransportOptions = append(o.ServeTransportOptions, transport.WithReadTimeout(timeout))
	}
}

func WithWriteTimeout(timeout uint32) Option{
	return func(o *Options) {
		o.ServeTransportOptions = append(o.ServeTransportOptions, transport.WithWriteTimeout(timeout))
	}
}

func WithInMsgChanSize(size uint32) Option{
	return func(o *Options){
		o.ServeTransportOptions = append(o.ServeTransportOptions, transport.WithInMsgChanSize(size))
	}
}

func WithOutMsgChanSize(size uint32) Option{
	return func(o *Options){
		o.ServeTransportOptions = append(o.ServeTransportOptions, transport.WithOutMsgChanSize(size))
	}
}

// WithTransport 替换底层server通信层
func WithTransport(t transport.ServerTransport) Option {
	return func(o *Options) {
		o.ServeTransport = t
	}
}

func WithMsgProtocol(p protocol.MsgProtocol) Option {
	return func(o *Options){
		o.MsgProtocol = p
	}
}