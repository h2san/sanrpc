package transport

// Transport 传输层接口
type ServerTransport interface {
	ListenAndServer(opt ...TransportOption) error
	Close() error
}

var DefaultTransport = NewTCPTransport()
var DefaultHttpTransport = NewHTTPTransport()