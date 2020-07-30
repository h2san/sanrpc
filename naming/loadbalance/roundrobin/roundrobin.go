package roundrobin

import (
	"sync"
	"time"

	"git.code.oa.com/trpc-go/trpc-go/naming/loadbalance"
	"git.code.oa.com/trpc-go/trpc-go/naming/registry"
)

var defaultUpdateRate time.Duration = time.Second * 10

func init() {
	loadbalance.Register("round_robin", NewRoundRobin(defaultUpdateRate))
}

// NewRoundRobin 新建实例
func NewRoundRobin(interval time.Duration) *RoundRobin {
	if interval == 0 {
		interval = defaultUpdateRate
	}
	return &RoundRobin{
		pickers:  new(sync.Map),
		interval: interval,
	}
}

// RoundRobin round robin 对象
type RoundRobin struct {
	pickers  *sync.Map
	interval time.Duration
}

// Select 实现 loadbalance 接口
func (rr *RoundRobin) Select(serviceName string, list []*registry.Node,
	opt ...loadbalance.Option) (*registry.Node, error) {

	opts := &loadbalance.Options{}
	for _, o := range opt {
		o(opts)
	}
	p, ok := rr.pickers.Load(serviceName)
	if ok {
		return p.(*rrPicker).Pick(list, opts)
	}

	newPicker := &rrPicker{
		interval: rr.interval,
	}
	v, ok := rr.pickers.LoadOrStore(serviceName, newPicker)
	if !ok {
		return newPicker.Pick(list, opts)
	}
	return v.(*rrPicker).Pick(list, opts)
}

type rrPicker struct {
	list     []*registry.Node
	updated  time.Time
	mu       sync.Mutex
	next     int
	interval time.Duration
}

// Pick 选择一个地址
func (p *rrPicker) Pick(list []*registry.Node, opts *loadbalance.Options) (*registry.Node, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.updateState(list)
	if len(p.list) == 0 {
		return nil, loadbalance.ErrNoServerAvailable
	}
	node := p.list[p.next]
	p.next = (p.next + 1) % len(p.list)

	return node, nil
}

func (p *rrPicker) updateState(list []*registry.Node) {
	if len(p.list) == 0 ||
		len(p.list) != len(list) ||
		time.Since(p.updated) > p.interval {
		p.list = list
		p.updated = time.Now()
		p.next = 0
		return
	}
}
