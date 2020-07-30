package client

import "context"

type tcpTransport struct {
	Connect
}

type ServerTransport interface {
	ListenAndServer() error
	Close() error
}

// ClientTransport client通讯层接口
type ClientTransport interface {
	Connect() error
	RoundTrip(ctx context.Context, req interface{}, resp interface{})  error
	Close()
}