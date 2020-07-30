package loadbalance

import (
	"testing"

	"git.code.oa.com/trpc-go/trpc-go/naming/registry"

	"github.com/stretchr/testify/assert"
)

func TestRandomEmptyList(t *testing.T) {
	b := &Random{}
	_, err := b.Select("", nil)
	assert.Equal(t, err, ErrNoServerAvailable)
}

func TestRandomGet(t *testing.T) {
	b := &Random{}
	node, err := b.Select("", []*registry.Node{testNode})
	assert.Nil(t, err)
	assert.Equal(t, node, testNode)
}
