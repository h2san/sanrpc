package consistenthash

import (
	"fmt"
	"sync"
	"testing"

	"git.code.oa.com/trpc-go/trpc-go/naming/loadbalance"
	"git.code.oa.com/trpc-go/trpc-go/naming/registry"
	"github.com/stretchr/testify/assert"
)

//测试key是否起作用
//对于同一个key, 在同一个节点列表中，每次寻址，返回的节点应当是同一个
func TestConsistentHashGetOne(t *testing.T) {
	ch := NewConsistentHash()

	//test list 1
	n, err := ch.Select("test", list1, loadbalance.WithKey("123"))
	assert.Nil(t, err)
	expectAddr := n.Address
	n, err = ch.Select("test", list1, loadbalance.WithKey("123"))
	assert.Nil(t, err)
	assert.Equal(t, expectAddr, n.Address)

	n, err = ch.Select("test", list1, loadbalance.WithKey("123456"))
	assert.Nil(t, err)
	expectAddr = n.Address
	n, err = ch.Select("test", list1, loadbalance.WithKey("123456"))
	assert.Nil(t, err)
	assert.Equal(t, expectAddr, n.Address)

	n, err = ch.Select("test", list1, loadbalance.WithKey("12315"))
	assert.Nil(t, err)
	expectAddr = n.Address
	n, err = ch.Select("test", list1, loadbalance.WithKey("12315"))
	assert.Nil(t, err)
	assert.Equal(t, expectAddr, n.Address)

	//test list 4
	n, err = ch.Select("test", list4, loadbalance.WithKey("Pony"))
	assert.Nil(t, err)
	expectAddr = n.Address
	n, err = ch.Select("test", list4, loadbalance.WithKey("Pony"))
	assert.Nil(t, err)
	assert.Equal(t, expectAddr, n.Address)

	n, err = ch.Select("test", list4, loadbalance.WithKey("John"))
	assert.Nil(t, err)
	expectAddr = n.Address
	n, err = ch.Select("test", list4, loadbalance.WithKey("John"))
	assert.Nil(t, err)
	assert.Equal(t, expectAddr, n.Address)

	n, err = ch.Select("test", list4, loadbalance.WithKey("Jack"))
	assert.Nil(t, err)
	expectAddr = n.Address
	n, err = ch.Select("test", list4, loadbalance.WithKey("Jack"))
	assert.Nil(t, err)
	assert.Equal(t, expectAddr, n.Address)
}

//测试空node list
//应当返回预期中的报错
func TestNilList(t *testing.T) {
	ch := NewConsistentHash()
	n, err := ch.Select("test", nil, loadbalance.WithKey("123"))
	assert.Nil(t, n)
	assert.Equal(t, loadbalance.ErrNoServerAvailable, err)
}

//测试空opt
//opt的WithKey是必须的
//应当返回预期中的报错
func TestNilOpts(t *testing.T) {
	ch := NewConsistentHash()

	n, err := ch.Select("test", list1)
	assert.Nil(t, n)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "missing key")

	n, err = ch.Select("test", list1, loadbalance.WithKey("whatever"))
	assert.Nil(t, err)
	assert.NotNil(t, n)
}

//测试长度为1的 node list
//每次都应当返回同一个结果
func TestSingleNode(t *testing.T) {
	ch := NewConsistentHash()
	n, err := ch.Select("test", list2, loadbalance.WithKey("123"))
	assert.Nil(t, err)
	assert.Equal(t, list2[0].Address, n.Address)

	n, err = ch.Select("test", list2, loadbalance.WithKey("456"))
	assert.Nil(t, err)
	assert.Equal(t, list2[0].Address, n.Address)

	n, err = ch.Select("test", list2, loadbalance.WithKey("12306"))
	assert.Nil(t, err)
	assert.Equal(t, list2[0].Address, n.Address)

	n, err = ch.Select("test", list2, loadbalance.WithKey("JackChen"))
	assert.Nil(t, err)
	assert.Equal(t, list2[0].Address, n.Address)
}

//节点长度变化，会马上更新哈希环
func TestInterval(t *testing.T) {
	ch := NewConsistentHash()

	//list长度变化，会马上触发计算新的hash环
	n, err := ch.Select("test", list2, loadbalance.WithKey("123"))
	assert.Nil(t, err)
	assert.Equal(t, list2[0].Address, n.Address)

	n, err = ch.Select("test", list4, loadbalance.WithKey("123"))
	assert.Nil(t, err)
	assert.Equal(t, false, isInList(n.Address, list2))
	assert.Equal(t, true, isInList(n.Address, list4))
}

//测试节点的删除对object映射位置的影响
func TestSubNode(t *testing.T) {
	ch := NewConsistentHash()

	var address1, address2, address3 string
	n, err := ch.Select("test", list1, loadbalance.WithKey("123"))
	assert.Nil(t, err)
	address1 = n.Address

	n, err = ch.Select("test", list1, loadbalance.WithKey("123456"))
	assert.Nil(t, err)
	address2 = n.Address

	n, err = ch.Select("test", list1, loadbalance.WithKey("12315"))
	assert.Nil(t, err)
	address3 = n.Address

	deletedAddress := address1

	//对list1删除deletedAddress
	//只会影响deletedAddress对应的key，其它key不受影响
	listTmp := deleteNode(deletedAddress, list1)

	n, err = ch.Select("test", listTmp, loadbalance.WithKey("123"))
	assert.Nil(t, err)
	if address1 != deletedAddress {
		assert.Equal(t, address1, n.Address)
	} else {
		assert.NotEqual(t, address1, n.Address)
	}

	n, err = ch.Select("test", listTmp, loadbalance.WithKey("123456"))
	assert.Nil(t, err)
	if address2 != deletedAddress {
		assert.Equal(t, address2, n.Address)
	} else {
		assert.NotEqual(t, address2, n.Address)
	}

	n, err = ch.Select("test", listTmp, loadbalance.WithKey("12315"))
	assert.Nil(t, err)
	if address3 != deletedAddress {
		assert.Equal(t, address3, n.Address)
	} else {
		assert.NotEqual(t, address3, n.Address)
	}
}

//测试平衡性
func TestBalance(t *testing.T) {
	ch := NewConsistentHash()
	counter := make(map[string]int)
	for i := 0; i < 200; i++ {
		n, err := ch.Select("test", list1, loadbalance.WithKey(fmt.Sprintf("%d", i)), loadbalance.WithReplicas(100))
		assert.Nil(t, err)
		if _, ok := counter[n.Address]; !ok {
			counter[n.Address] = 0
		} else {
			counter[n.Address]++
		}
	}
	for _, v := range counter {
		//每个node都有节点落入
		assert.NotEqual(t, 0, v)
		fmt.Println(v)
	}
}

func TestIsNodeSliceEqualBCE(t *testing.T) {
	isEqual := isNodeSliceEqualBCE(list1, list2)
	assert.Equal(t, false, isEqual)
	isEqual = isNodeSliceEqualBCE(list1, list1)
	assert.Equal(t, true, isEqual)
	isEqual = isNodeSliceEqualBCE(list1, nil)
	assert.Equal(t, false, isEqual)
}

// 测试并发安全性，每次访问时，list都不一样
// 这是一种极端情况，事实上绝大部分情况下一个service对应的node list是不会频繁变化的，只会在扩缩容时有变化
func TestParallel(t *testing.T) {
	var wg sync.WaitGroup
	ch := NewConsistentHash()
	var lists [][]*registry.Node
	var keys []string
	var results []string

	n, err := ch.Select("test", list1, loadbalance.WithKey("1"))
	assert.Nil(t, err)
	results = append(results, n.Address)
	lists = append(lists, list1)
	keys = append(keys, "1")

	n, err = ch.Select("test", list2, loadbalance.WithKey("2"))
	assert.Nil(t, err)
	results = append(results, n.Address)
	lists = append(lists, list2)
	keys = append(keys, "2")

	n, err = ch.Select("test", list3, loadbalance.WithKey("3"))
	assert.Nil(t, err)
	results = append(results, n.Address)
	lists = append(lists, list3)
	keys = append(keys, "3")

	n, err = ch.Select("test", list4, loadbalance.WithKey("4"))
	assert.Nil(t, err)
	results = append(results, n.Address)
	lists = append(lists, list4)
	keys = append(keys, "4")

	n, err = ch.Select("test", list5, loadbalance.WithKey("5"))
	assert.Nil(t, err)
	results = append(results, n.Address)
	lists = append(lists, list5)
	keys = append(keys, "5")

	//模拟大量协程并发访问

	//这是一种极端情况，事实上绝大部分情况下一个service对应的node list是不会频繁变化的，只会在扩缩容时有变化
	for i := 0; i < 50; i++ {
		go func(i int) {
			wg.Add(1)
			defer wg.Done()
			n, err := ch.Select("test0", lists[i%5], loadbalance.WithKey(keys[i%5]))
			assert.Nil(t, err)
			assert.Equal(t, results[i%5], n.Address)
		}(i)
	}

	for i := 0; i < 50; i++ {
		go func(i int) {
			wg.Add(1)
			defer wg.Done()
			n, err := ch.Select("test1", lists[0], loadbalance.WithKey(keys[0]))
			assert.Nil(t, err)
			assert.Equal(t, results[0], n.Address)
		}(i)
	}
	wg.Wait()
}

// 测试并发访问性能
func BenchmarkParallel(b *testing.B) {
	ch := NewConsistentHash()
	b.SetParallelism(10) //十个并发协程
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = ch.Select("test", list1, loadbalance.WithKey("HelloWorld"))
		}
	})
}

var list1 = []*registry.Node{
	{
		Address: "list1.ip.1:8080",
	},
	{
		Address: "list1.ip.2:8080",
	},
	{
		Address: "list1.ip.3:8080",
	},
	{
		Address: "list1.ip.4:8080",
	},
}

var list2 = []*registry.Node{
	{
		Address: "list2.ip.1:8080",
	},
}

var list3 = []*registry.Node{
	{
		Address: "list3.ip.2:8080",
	},
	{
		Address: "list3.ip.4:8080",
	},
	{
		Address: "list3.ip.1:8080",
	},
}

var list4 = []*registry.Node{
	{
		Address: "list4.ip.168:8080",
	},
	{
		Address: "list4.ip.167:8080",
	},
	{
		Address: "list4.ip.15:8080",
	},
	{
		Address: "list4.ip.15:8081",
	},
}

var list5 = []*registry.Node{
	{
		Address: "list5.ip.2:8080",
	},
}

func deleteNode(adddress string, list []*registry.Node) []*registry.Node {
	ret := make([]*registry.Node, 0, len(list))
	for _, n := range list {
		if n.Address != adddress {
			ret = append(ret, n)
		}
	}
	return ret
}

func isInList(adddress string, list []*registry.Node) bool {
	for _, n := range list {
		if n.Address == adddress {
			return true
		}
	}
	return false
}
