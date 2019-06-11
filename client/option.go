package client

import (
	"crypto/tls"
	"time"

	"github.com/h2san/sanrpc/share"
)

type Option struct {
	Retries        int
	TLSConfig      *tls.Config
	RPCPath        string
	ConnectTimeout time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration

	Heartbeat         bool
	HeartbeatInterval time.Duration
}

// DefaultOption is a common option configuration for client.
var DefaultOption = Option{
	Retries:        3,
	RPCPath:        share.DefaultRPCPath,
	ConnectTimeout: 10 * time.Second,
}
