package circuitbreaker

import (
	"testing"
	"time"

	"git.code.oa.com/trpc-go/trpc-go/naming/registry"

	"github.com/stretchr/testify/assert"
)

type testCircuitBreaker struct{}

func (cb *testCircuitBreaker) Available(node *registry.Node) bool {
	return true
}

func (cb *testCircuitBreaker) Report(node *registry.Node, cost time.Duration, err error) error {
	return nil
}

func TestCircuitBreakerRegister(t *testing.T) {
	Register("cb", &testCircuitBreaker{})
	unregisterForTesting("cb")
}

func TestCircuitBreakerGet(t *testing.T) {
	Register("cb", &testCircuitBreaker{})
	assert.NotNil(t, Get("cb"))
	unregisterForTesting("cb")
	assert.Nil(t, Get("not_exist"))
}

func TestNoopCircuitBreaker(t *testing.T) {
	noop := &NoopCircuitBreaker{}
	assert.True(t, noop.Available(nil))
	assert.Nil(t, noop.Report(nil, 0, nil))
}
