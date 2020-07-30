package consistenthash

import (
	"errors"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
	"time"

	"git.code.oa.com/trpc-go/trpc-go/naming/loadbalance"
	"git.code.oa.com/trpc-go/trpc-go/naming/registry"
)

// 默认虚拟节点系数
var defaultReplicas int = 100

type Hash func(data []byte) uint32

// 默认使用CRC32算法
var defaultHashFunc Hash = crc32.ChecksumIEEE

func init() {
	loadbalance.Register("consistent_hash", NewConsistentHash())
}

// NewConsistentHash 新建实例
func NewConsistentHash() *ConsistentHash {
	return &ConsistentHash{
		pickers: new(sync.Map),
	}
}

// ConsistentHash consistent hash 对象
type ConsistentHash struct {
	pickers  *sync.Map
	interval time.Duration
}

// Select 实现 loadbalance 接口
func (ch *ConsistentHash) Select(serviceName string, list []*registry.Node,
	opt ...loadbalance.Option) (*registry.Node, error) {

	opts := &loadbalance.Options{}
	for _, o := range opt {
		o(opts)
	}
	p, ok := ch.pickers.Load(serviceName)
	if ok {
		return p.(*chPicker).Pick(list, opts)
	}

	newPicker := &chPicker{
		interval: ch.interval,
	}
	v, ok := ch.pickers.LoadOrStore(serviceName, newPicker)
	if !ok {
		return newPicker.Pick(list, opts)
	}
	return v.(*chPicker).Pick(list, opts)
}

type chPicker struct {
	list     []*registry.Node
	keys     UInt32Slice               // 已排序的节点的哈希的slice, 其长度为节点数*replicas
	hashMap  map[uint32]*registry.Node // 保存hash-node映射关系的map
	updated  time.Time
	mu       sync.Mutex
	interval time.Duration
}

// Pick 选择一个地址
func (p *chPicker) Pick(list []*registry.Node, opts *loadbalance.Options) (*registry.Node, error) {
	if len(list) == 0 {
		return nil, loadbalance.ErrNoServerAvailable
	}
	//必须要传入opts的key，不然直接返回错误
	if opts.Key == "" {
		return nil, errors.New("missing key")
	}
	tmpKeys, tmpMap, err := p.updateState(list, opts.Replicas)
	if err != nil {
		return nil, err
	}
	hash := defaultHashFunc([]byte(opts.Key))
	//二分查找获取最优节点，第一个节点hash值大于对象hash值的就是最优节点
	idx := sort.Search(len(tmpKeys), func(i int) bool { return tmpKeys[i] >= hash })
	if idx == len(tmpKeys) {
		idx = 0
	}
	node, ok := tmpMap[tmpKeys[idx]]
	if !ok {
		return nil, loadbalance.ErrNoServerAvailable
	}
	return node, nil
}

//每隔一段时间，若节点变化，才进行重新计算列表
func (p *chPicker) updateState(list []*registry.Node, replicas int) (UInt32Slice, map[uint32]*registry.Node, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	//检查上次更新的node list和本次的node是否相同，若相同，那就不必更新哈希环
	if isNodeSliceEqualBCE(p.list, list) {
		return p.keys, p.hashMap, nil
	}
	actualReplicas := replicas
	if actualReplicas <= 0 {
		actualReplicas = defaultReplicas
	}
	//更新节点
	//记录本次的node list
	p.list = list
	p.hashMap = make(map[uint32]*registry.Node)
	p.keys = make(UInt32Slice, len(list)*actualReplicas)
	for i, node := range list {
		if node == nil {
			//不允许有空的node
			return nil, nil, errors.New("list contains nil node")
		}
		for j := 0; j < actualReplicas; j++ {
			hash := defaultHashFunc([]byte(strconv.Itoa(j) + node.Address))
			p.keys[i*(actualReplicas)+j] = hash
			p.hashMap[hash] = node
		}
	}
	sort.Sort(p.keys)
	return p.keys, p.hashMap, nil
}

type UInt32Slice []uint32

func (s UInt32Slice) Len() int {
	return len(s)
}

func (s UInt32Slice) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s UInt32Slice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

//golang支持BCE，判断两个node list是否相等
func isNodeSliceEqualBCE(a, b []*registry.Node) bool {
	if len(a) != len(b) {
		return false
	}
	if (a == nil) != (b == nil) {
		return false
	}
	b = b[:len(a)]
	for i, v := range a {
		if (v == nil) != (b[i] == nil) {
			return false
		}
		if v.Address != b[i].Address {
			return false
		}
	}
	return true
}
