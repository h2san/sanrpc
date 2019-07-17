package httpx

import (
	"context"
	"net/http"
)

type Plugin interface {
}

type (
	PreHandleHTTPRequestPlugin interface {
		PreHandleHTTPRequest(ctx context.Context, req *http.Request) error
	}

	PostHandleHTTPRequestPlugin interface {
		PostHandleHTTPRequest(ctx context.Context)
	}
)

type httpxPlugin struct {
	plugins []Plugin
}

func (p *httpxPlugin) DoPreHandleHTTPRequest(ctx context.Context, req *http.Request) error {
	for i := range p.plugins {
		if plugin, ok := p.plugins[i].(PreHandleHTTPRequestPlugin); ok {
			err := plugin.PreHandleHTTPRequest(ctx, req)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *httpxPlugin) DoPostHandleHTTPRequest(ctx context.Context) {
	for i := range p.plugins {
		if plugin, ok := p.plugins[i].(PostHandleHTTPRequestPlugin); ok {
			plugin.PostHandleHTTPRequest(ctx)
		}
	}
}

// Add adds a plugin.
func (p *httpxPlugin) Add(plugin Plugin) {
	p.plugins = append(p.plugins, plugin)
}

// Remove removes a plugin by it's name.
func (p *httpxPlugin) Remove(plugin Plugin) {
	if p.plugins == nil {
		return
	}

	var plugins []Plugin
	for _, p := range p.plugins {
		if p != plugin {
			plugins = append(plugins, p)
		}
	}

	p.plugins = plugins
}
