package service

import (
	"github.com/hillguo/sanrpc/protocol"
)

type Options struct {
	Name string
	InMsgChanSize uint32
	OutMsgChanSize uint32
	ReadTimeout uint32
	WriteTimeout uint32

	Address        string
	NetWork        string

	MsgProtocol    protocol.MsgProtocol
}

type Option func(options *Options)

func WithServiceName(s string) Option{
	return func(o *Options){
		o.Name = s
	}
}

func WithAddress(s string) Option{
	return func(o *Options) {
		o.Address = s
	}
}

func WithNetWork(network string) Option{
	return func(o *Options) {
		o.NetWork = network
	}
}

func WithReadTimeout(timeout uint32) Option{
	return func(o *Options) {
		o.ReadTimeout = timeout
	}
}

func WithWriteTimeout(timeout uint32) Option{
	return func(o *Options) {
		o.WriteTimeout = timeout
	}
}

func WithInMsgChanSize(size uint32) Option{
	return func(o *Options){
		o.InMsgChanSize = size
	}
}

func WithOutMsgChanSize(size uint32) Option{
	return func(o *Options){
		o.OutMsgChanSize = size
	}
}

func WithMsgProtocol(p protocol.MsgProtocol) Option {
	return func(o *Options){
		o.MsgProtocol = p
	}
}