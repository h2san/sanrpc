package transport

import "github.com/hillguo/sanrpc/protocol"

type TransportOptions struct {
	InMsgChanSize uint32
	OutMsgChanSize uint32
	ReadTimeout uint32
	WriteTimeout uint32

	Address string
	NetWork string
	TLSCertFile string
	TLSKeyFile string

	MsgProtocol protocol.MsgProtocol

}

type TransportOption func(*TransportOptions)

func WithInMsgChanSize(size uint32) TransportOption{
	return func(options *TransportOptions){
		options.InMsgChanSize = size
	}
}

func WithOutMsgChanSize(size uint32) TransportOption{
	return func(options *TransportOptions){
		options.OutMsgChanSize = size
	}
}

func WithWriteTimeout(size uint32) TransportOption{
	return func(options *TransportOptions){
		options.WriteTimeout = size
	}
}

func WithReadTimeout(size uint32) TransportOption{
	return func(options *TransportOptions){
		options.ReadTimeout = size
	}
}

func WithAddress(address string) TransportOption{
	return func(options *TransportOptions){
		options.Address = address
	}
}

func WithNetWork(network string) TransportOption{
	return func(options *TransportOptions){
		options.NetWork = network
	}
}

func WithTLS(certFile, keyFile string) TransportOption{
	return func(options *TransportOptions){
		options.TLSCertFile = certFile
		options.TLSKeyFile = keyFile
	}
}

func WithMsgProtocol(p protocol.MsgProtocol)TransportOption{
	return func(options *TransportOptions){
		options.MsgProtocol = p
	}
}

