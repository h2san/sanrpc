package selector

import (
	"github.com/hillguo/sanrpc/client/node"
)

// Selector 路由组件接口
type Selector interface {
	Select(list []*node.Node) (node *node.Node, err error)
}

var (
	selectors = make(map[string]Selector)
)

var DefaultSelector = &RandomSelector{}

// Register 注册selector，如l5 dns cmlb tseer
func Register(name string, s Selector) {
	selectors[name] = s
}

// Get 获取selector
func GetSelector(name string) Selector {
	s := selectors[name]
	return s
}

