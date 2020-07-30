package loadbalance

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	opts := &Options{}
	WithNamespace("ns")(opts)
	WithInterval(time.Second * 2)(opts)
	WithKey("hash key")(opts)
	assert.Equal(t, opts.Namespace, "ns")
	assert.Equal(t, opts.Interval, time.Second*2)
	assert.Equal(t, opts.Key, "hash key")
}
