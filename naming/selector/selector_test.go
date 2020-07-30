package selector

import (
	"testing"
	"time"

	"git.code.oa.com/trpc-go/trpc-go/naming/registry"

	"github.com/stretchr/testify/assert"
)

var testNode *registry.Node = &registry.Node{
	ServiceName: "testservice",
	Address:     "testservice.ip.1:16721",
	Network:     "tcp",
}

type testSelector struct {
}

func (ts *testSelector) Select(serviceName string, opt ...Option) (*registry.Node, error) {
	return testNode, nil
}

func (ts *testSelector) Report(node *registry.Node, cost time.Duration, success error) error {
	return nil
}

func TestSelectorRegister(t *testing.T) {
	Register("test-selector", &testSelector{})
	unregisterForTesting("test-selector")
}

func TestSelectorGet(t *testing.T) {
	Register("test-selector", &testSelector{})
	s := Get("test-selector")
	assert.NotNil(t, s)
	unregisterForTesting("test-selector")
	assert.Nil(t, Get("not_exist"))
}