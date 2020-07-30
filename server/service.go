package server

import (
	"context"
	"errors"
	"fmt"
	log "github.com/hillguo/sanlog"
	"github.com/hillguo/sanrpc/errs"
	"github.com/hillguo/sanrpc/filter"
	"github.com/hillguo/sanrpc/metrics"
	"github.com/hillguo/sanrpc/msg"
	"github.com/hillguo/sanrpc/naming/registry"
)

// Service 服务端调用结构
type Service interface {
	// 注册路由信息
	Register(serviceDesc interface{}, serviceImpl interface{}) error
	// 启动服务
	Serve() error
	// 关闭服务
	Close(chan struct{}) error
}

// FilterFunc 内部解析reqbody 并返回拦截器 传给server stub
type FilterFunc func(reqbody interface{}) (filter.Chain, error)

// Method 服务rpc方法信息
type Method struct {
	Name string
	Func func(svr interface{}, ctx context.Context, f FilterFunc) (rspbody interface{}, err error)
}

// ServiceDesc 服务描述service定义
type ServiceDesc struct {
	ServiceName string
	HandlerType interface{}
	Methods     []Method
}

// Handler trpc默认的handler
type Handler func(ctx context.Context, f FilterFunc) (rspbody interface{}, err error)

// service Service实现
type service struct {
	ctx    context.Context    // service关闭
	cancel context.CancelFunc // service关闭

	opts *Options // service选项

	handlers map[string]Handler // rpcname => handler
}

// Serve 启动服务
func (s *service) Serve() (err error) {

	// 确保正常监听之后才能启动服务注册
	if err = s.opts.Transport.ListenAndServe(s.ctx, s.opts.ServeOptions...); err != nil {
		log.Errorf("service:%s ListenAndServe fail:%v", s.opts.ServiceName, err)
		return err
	}

	if s.opts.Registry != nil {
		err = s.opts.Registry.Register(s.opts.ServiceName, registry.WithAddress(s.opts.Address))
		if err != nil {
			// 有注册失败，关闭service，并返回给上层错误
			log.Errorf("service:%s register fail:%v", s.opts.ServiceName, err)
			return err
		}
	}

	log.Infof("%s service:%s launch success, address:%s, serving ...", s.opts.protocol, s.opts.ServiceName, s.opts.Address)

	metrics.Counter("ServiceStart").Incr()
	<-s.ctx.Done()
	return nil
}

// Handle server transport收到请求包后调用此函数
func (s *service) Handle(ctx context.Context, reqbuf []byte) (rspbuf []byte, err error) {

	// 无法回包，只能丢弃
	if s.opts.Codec == nil {
		log.Error(ctx, "server codec empty")
		metrics.Counter("ServerCodecEmpty").Incr()
		return nil, errors.New("server codec empty")
	}

	var rspbodybuf []byte

	msg := msg.Message(ctx)

	rspbody, err := s.handle(ctx, msg, reqbuf)
	if err != nil {
		// 不回包
		if err == errs.ErrServerNoResponse {
			return nil, err
		}
		// 处理失败 给客户端返回错误码, 忽略rspbody
		metrics.Counter("ServiceHandleFail").Incr()
		return s.encode(ctx, msg, rspbodybuf, err)
	}

	// 业务处理成功 才开始打包body
	rspbodybuf, err = codec.Marshal(msg.SerializationType(), rspbody)
	if err != nil {
		metrics.Counter("ServiceCodecMarshalFail").Incr()
		err = errs.NewFrameError(errs.RetServerEncodeFail, "service codec Marshal: "+err.Error())
		// 处理失败 给客户端返回错误码
		return s.encode(ctx, msg, rspbodybuf, err)
	}

	// 处理成功 才开始压缩body
	rspbodybuf, err = codec.Compress(msg.CompressType(), rspbodybuf)
	if err != nil {
		metrics.Counter("ServiceCodecCompressFail").Incr()
		err = errs.NewFrameError(errs.RetServerEncodeFail, "service codec Compress: "+err.Error())
		// 处理失败 给客户端返回错误码
		return s.encode(ctx, msg, rspbodybuf, err)
	}

	return s.encode(ctx, msg, rspbodybuf, nil)
}

func (s *service) handle(ctx context.Context, msg msg.Msg, reqbuf []byte) (rspbody interface{}, err error) {

	msg.WithNamespace(s.opts.Namespace)           // server 的命名空间
	msg.WithEnvName(s.opts.EnvName)               // server 的环境
	msg.WithSetName(s.opts.SetName)               // server 的set
	msg.WithCalleeServiceName(s.opts.ServiceName) // 以server角度看，caller是上游，callee是自身

	reqbodybuf, err := s.opts.Codec.Decode(msg, reqbuf)
	if err != nil {
		metrics.Counter("ServiceCodecDecodeFail").Incr()
		return nil, errs.NewFrameError(errs.RetServerDecodeFail, "service codec Decode: "+err.Error())
	}

	// 再赋值一遍，防止decode更改了
	msg.WithNamespace(s.opts.Namespace)           // server 的命名空间
	msg.WithEnvName(s.opts.EnvName)               // server 的环境
	msg.WithSetName(s.opts.SetName)               // server 的set
	msg.WithCalleeServiceName(s.opts.ServiceName) // 以server角度看，caller是上游，callee是自身

	handler, ok := s.handlers[msg.ServerRPCName()]
	if !ok {
		handler, ok = s.handlers["*"] // 支持通配符全匹配转发处理
		if !ok {
			metrics.Counter("ServiceHandleRpcNameInvalid").Incr()
			return nil, errs.NewFrameError(errs.RetServerNoFunc,
				fmt.Sprintf("service handle: rpc name %s invalid", msg.ServerRPCName()))
		}
	}

	timeout := s.opts.Timeout
	if msg.RequestTimeout() > 0 && !s.opts.DisableRequestTimeout {
		if msg.RequestTimeout() < timeout || timeout == 0 {
			timeout = msg.RequestTimeout()
		}
	}
	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	return handler(ctx, s.filterFunc(msg, reqbodybuf))
}