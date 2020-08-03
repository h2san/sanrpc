package client

// Options 客户端调用参数
type Options struct {
	ServiceName string // 调用服务名
	MethodName  string // 调用方法名

	//
	Network string
	Discovery string
	Selector string
	Address string

	ConnectTimeout uint64
	ReadTimeout    uint64
	WriteTimeout   uint64
}

// Option 调用参数工具函数
type Option func(*Options)

func WithServiceName(ServiceName string) Option{
	return func(o *Options){
		o.ServiceName = ServiceName
	}
}

func WithMethodName(MethodName string) Option{
	return func(o *Options){
		o.MethodName = MethodName
	}
}

func WithNetwork(Network string) Option{
	return func(o *Options){
		o.Network = Network
	}
}

func WithDiscovery(Discovery string) Option{
	return func(o *Options){
		o.Discovery = Discovery
	}
}

func WithSelector(Selector string) Option{
	return func(o *Options){
		o.Selector = Selector
	}
}

func WithAddress(Address string) Option{
	return func(o *Options){
		o.Address = Address
	}
}

func WithConnectTimeout(ConnectTimeout uint64) Option{
	return func(o *Options){
		o.ConnectTimeout = ConnectTimeout
	}
}

func WithWriteTimeout(WriteTimeout uint64) Option{
	return func(o *Options){
		o.WriteTimeout = WriteTimeout
	}
}

func WithReadTimeout(ReadTimeout uint64) Option{
	return func(o *Options){
		o.ReadTimeout = ReadTimeout
	}
}