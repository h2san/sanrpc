package registry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testRegistry struct{}

func (r *testRegistry) Register(service string, opt ...Option) error {
	return nil
}
func (r *testRegistry) Deregister(service string) error {
	return nil
}

func TestRegistryRegister(t *testing.T) {
	Register("test-registry", &testRegistry{})
	unregisterForTesting("test-registry")
}

func TestRegistryGet(t *testing.T) {
	Register("test-registry", &testRegistry{})
	r := Get("test-registry")
	assert.Nil(t, r.Register("service1", nil))
	assert.Nil(t, r.Deregister("service1"))
	unregisterForTesting("test-registry")
}

func TestNoopRegister(t *testing.T) {
	noop := &NoopRegistry{}
	assert.Equal(t, noop.Register("test", nil), ErrNotImplement)
	assert.Equal(t, noop.Deregister("test"), ErrNotImplement)
}
