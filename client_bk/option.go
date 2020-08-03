package client_bk

import (
	"crypto/tls"
	"time"
)

//Option ...
type Option struct {
	Retries        int
	TLSConfig      *tls.Config
	ConnectTimeout time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration

	Heartbeat         bool
	HeartbeatInterval time.Duration
}

// DefaultOption is a common option configuration for client_bk.
var DefaultOption = Option{
	Retries:        3,
	ConnectTimeout: 10 * time.Second,
}
