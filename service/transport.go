package service

// Transport 传输层接口
type ServerTransport interface {
	ListenAndServer() error
	Close() error
}