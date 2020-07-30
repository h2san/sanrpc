package selector

import (
	"testing"

	"git.code.oa.com/trpc-go/trpc-go/naming/registry"

	"github.com/stretchr/testify/assert"
)

func TestTrpcSelectorSelect(t *testing.T) {
	selector := &TrpcSelector{}
	n, err := selector.Select("10.100.72.229.12367")
	assert.Nil(t, err)
	assert.Equal(t, n.Address, "10.100.72.229.12367")
}

func TestTrpcSelectorReport(t *testing.T) {
	selector := &TrpcSelector{}
	n, err := selector.Select("10.100.72.229.12367")

	assert.Nil(t, err)
	assert.Equal(t, n.Address, "10.100.72.229.12367")

	assert.Nil(t, selector.Report(n, 0, nil))
}

func TestTrpcSelectorReportErr(t *testing.T) {
	selector := &TrpcSelector{}
	assert.Equal(t, selector.Report(nil, 0, nil), ErrReportNodeEmpty)
	assert.Equal(t, selector.Report(&registry.Node{}, 0, nil), ErrReportMetaDataEmpty)
	assert.Equal(t, selector.Report(&registry.Node{
		Metadata: make(map[string]interface{}),
	}, 0, nil), ErrReportNoCircuitBreaker)
	assert.Equal(t, selector.Report(&registry.Node{
		Metadata: map[string]interface{}{
			"circuitbreaker": "circuitbreaker",
		},
	}, 0, nil), ErrReportInvalidCircuitBreaker)
}
