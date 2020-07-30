package server

import (
	"github.com/hillguo/sanrpc/codec"
	"github.com/hillguo/sanrpc/filter"
	"github.com/hillguo/sanrpc/naming/registry"
	"time"
)

// Options 服务端调用参数
type Options struct {
	Namespace   string // 当前服务命名空间 正式环境 Production 测试环境 Development
	EnvName     string // 当前环境
	SetName     string // set分组
	ServiceName string // 当前服务的 service name

	Address               string        // 监听地址 ip:port
	Timeout               time.Duration // 请求最长处理时间
	DisableRequestTimeout bool          // 禁用继承上游的超时时间

	protocol   string // 业务协议 trpc http ...
	handlerSet bool   // 用户是否自己定义handler

	// ServeOptions []transport.ListenServeOption
	// Transport    transport.ServerTransport

	Registry registry.Registry
	Codec    codec.Codec

	Filters filter.Chain // 链式拦截器
}
