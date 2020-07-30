package selector

import (
	"testing"

	"git.code.oa.com/trpc-go/trpc-go/naming/circuitbreaker"
	"git.code.oa.com/trpc-go/trpc-go/naming/discovery"
	"git.code.oa.com/trpc-go/trpc-go/naming/loadbalance"

	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	opts := &Options{}
	WithKey("key")(opts)
	WithSourceSetName("set")(opts)
	WithDestinationSetName("dstSet")(opts)
	d := &discovery.IPDiscovery{}
	WithDiscovery(d)(opts)
	b := &loadbalance.Random{}
	WithLoadBalance(b)(opts)
	cb := &circuitbreaker.NoopCircuitBreaker{}
	WithCircuitBreaker(cb)(opts)
	WithDisableServiceRouter()(opts)
	WithDestinationEnvName("dst_env")(opts)
	WithNamespace("test_namespace")(opts)
	WithSourceNamespace("src_namespace")(opts)
	WithEnvKey("env_key")(opts)
	WithSourceServiceName("src_svcname")(opts)
	WithSourceEnvName("src_env")(opts)

	assert.Equal(t, opts.SourceSetName, "set")
	assert.Equal(t, opts.Key, "key")
	assert.Equal(t, opts.CircuitBreaker, cb)
	assert.Equal(t, opts.LoadBalance, b)
	assert.Equal(t, opts.Discovery, d)
	assert.True(t, opts.DisableServiceRouter)
	assert.Equal(t, opts.DestinationEnvName, "dst_env")
	assert.Equal(t, opts.DestinationSetName, "dstSet")
	assert.Equal(t, opts.Namespace, "test_namespace")
	assert.Equal(t, opts.SourceNamespace, "src_namespace")
	assert.Equal(t, opts.SourceServiceName, "src_svcname")
	assert.Equal(t, opts.EnvKey, "env_key")
	assert.Equal(t, opts.SourceEnvName, "src_env")
}
