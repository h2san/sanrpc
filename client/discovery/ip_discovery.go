package discovery

import (
	"errors"
	"github.com/hillguo/sanrpc/client/node"
	"strings"
)

func init() {
	Register("ip", &IPDiscovery{})
}

// IPDiscovery  ip列表服务发现
type IPDiscovery struct{}

// List 返回原始ipport
func (*IPDiscovery) List(serviceName string) ([]*node.Node, error) {
	ips := strings.Split(serviceName,",")
	if len(ips) ==0 {
		return nil, errors.New("no ip discovery")
	}
	nodes := make([]*node.Node,0,len(ips))
	for _,ip := range ips {
		nodes = append(nodes, &node.Node{Address: ip})
	}
	return nodes, nil
}
