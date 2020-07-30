package discovery

import (
	"github.com/hillguo/sanrpc/client/node"
	"sync"
)

type Discovery interface {
	List(serviceName string) ([]*node.Node,error)
}

var (
	discoveries = make(map[string]Discovery)
	lock        = sync.RWMutex{}
)

var DefaultDiscovery = &IPDiscovery{}


// Register 注册discovery
func Register(name string, s Discovery) {
	lock.Lock()
	discoveries[name] = s
	lock.Unlock()
}

// Get 获取discovery
func Get(name string) Discovery {
	lock.RLock()
	d := discoveries[name]
	lock.RUnlock()
	return d
}