package filter_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"git.code.oa.com/trpc-go/trpc-go/filter"
)

func echoHandle(ctx context.Context, req interface{}, rsp interface{}) error {
	preq := req.(*string)
	prsp := rsp.(*string)
	*prsp = *preq
	rsp = prsp

	return nil
}

func TestNoopFilter(t *testing.T) {

	req := "echo"
	rsp := ""
	err := filter.NoopFilter(context.Background(), &req, &rsp, echoHandle)
	assert.Nil(t, err)
	assert.Equal(t, rsp, req)
}

func TestFilterChain_Handle(t *testing.T) {
	req := "echo"
	rsp := ""
	//noopFilter
	{
		fc := filter.Chain{}
		err := fc.Handle(context.Background(), &req, &rsp, echoHandle)
		assert.Nil(t, err)
		assert.Equal(t, rsp, req)
	}

	//oneFilter
	{
		fc := filter.Chain{filter.NoopFilter}
		err := fc.Handle(context.Background(), &req, &rsp, echoHandle)
		assert.Nil(t, err)
		assert.Equal(t, rsp, req)
	}

	// multiFilter
	{
		fc := filter.Chain{filter.NoopFilter, filter.NoopFilter, filter.NoopFilter}
		err := fc.Handle(context.Background(), &req, &rsp, echoHandle)
		assert.Nil(t, err)
		assert.Equal(t, rsp, req)
	}
}

func TestGetClient(t *testing.T) {
	filter.Register("noop", filter.NoopFilter, filter.NoopFilter)
	f := filter.GetClient("noop")
	assert.NotNil(t, f)
}

func TestGetServer(t *testing.T) {
	filter.Register("noop", filter.NoopFilter, filter.NoopFilter)
	f := filter.GetServer("noop")
	assert.NotNil(t, f)
}
