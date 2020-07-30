package client

import (
	"time"
)

// Options 客户端调用参数
type Options struct {
	ServiceName string // 调用服务名
	MethodName  string // 调用方法名
	ConnectTimeout time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration

	//
	Network string
	Discovery string
	Selector string
	//
	Protocol    string
}

// Option 调用参数工具函数
type Option func(*Options)