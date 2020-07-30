package discovery

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	opts := &Options{}
	WithNamespace("ns")(opts)
	assert.Equal(t, opts.Namespace, "ns")
}
