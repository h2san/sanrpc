package discovery

import (
	"testing"

	"git.code.oa.com/trpc-go/trpc-go/naming/registry"

	"github.com/stretchr/testify/assert"
)

var testNode *registry.Node = &registry.Node{
	ServiceName: "testservice",
	Address:     "testservice.ip.1:16721",
	Network:     "tcp",
}

type testDiscovery struct{}

func (d *testDiscovery) List(serviceName string, opt ...Option) ([]*registry.Node, error) {
	return []*registry.Node{testNode}, nil
}

func TestDiscoveryRegister(t *testing.T) {
	Register("test-discovery", &testDiscovery{})
	unregisterForTesting("test-discovery")
}

func TestDiscoveryGet(t *testing.T) {
	Register("test-discovery", &testDiscovery{})
	assert.NotNil(t, Get("test-discovery"))
	unregisterForTesting("test-discovery")
	assert.Nil(t, Get("not_exist"))
}

func TestDiscoveryList(t *testing.T) {
	Register("test-discovery", &testDiscovery{})
	d := Get("test-discovery")
	list, err := d.List("test-service", nil)
	assert.Nil(t, err)
	assert.Equal(t, list[0], testNode)
	unregisterForTesting("test-discovery")
}
