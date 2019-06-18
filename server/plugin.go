package server

import (
	"context"
	"github.com/hillguo/sanrpc/errors"
	"github.com/hillguo/sanrpc/protocol"
	"net"
	"net/http"
)

//PluginContainer represents a plugin container that defines all methods to manage plugins.
//And it also defines all extension points.
type PluginContainerer interface {
	Add(plugin Plugin)
	Remove(plugin Plugin)
	All() []Plugin

	DoRegister(name string, rcvr interface{}, metadata string) error
	DoRegisterFunction(name string, fn interface{}, metadata string) error
	DoUnregister(name string) error

	DoPostConnAccept(conn net.Conn) (net.Conn, bool)
	DoPostConnClose(conn net.Conn) bool

	DoPreReadRequest(ctx context.Context) error
	DoPostReadRequest(ctx context.Context, req protocol.Message, e error) error

	DoPreHandleRequest(ctx context.Context, req protocol.Message) error

	DoPreWriteResponse(context.Context, protocol.Message, protocol.Message) error
	DoPostWriteResponse(context.Context, protocol.Message, protocol.Message, error) error

	DoPreWriteRequest(ctx context.Context) error
	DoPostWriteRequest(ctx context.Context, r protocol.Message, e error) error

	// http
	DoPreHandlerHTTPRequest(w http.ResponseWriter, r *http.Request) error
	DoPostHandlerHTTPRequest(w http.ResponseWriter, r *http.Request) error
}

// Plugin is the server plugin interface.
type Plugin interface {
}

type (
	// RegisterPlugin is .
	RegisterPlugin interface {
		Register(name string, rcvr interface{}, metadata string) error
		Unregister(name string) error
	}

	// RegisterFunctionPlugin is .
	RegisterFunctionPlugin interface {
		RegisterFunction(name string, fn interface{}, metadata string) error
	}

	// PostConnAcceptPlugin represents connection accept plugin.
	// if returns false, it means subsequent IPostConnAcceptPlugins should not contiune to handle this conn
	// and this conn has been closed.
	PostConnAcceptPlugin interface {
		HandleConnAccept(net.Conn) (net.Conn, bool)
	}

	// PostConnClosePlugin represents client connection close plugin.
	PostConnClosePlugin interface {
		HandleConnClose(net.Conn) bool
	}

	//PreReadRequestPlugin represents .
	PreReadRequestPlugin interface {
		PreReadRequest(ctx context.Context) error
	}

	//PostReadRequestPlugin represents .
	PostReadRequestPlugin interface {
		PostReadRequest(ctx context.Context, r protocol.Message, e error) error
	}

	//PreHandleRequestPlugin represents .
	PreHandleRequestPlugin interface {
		PreHandleRequest(ctx context.Context, r protocol.Message) error
	}

	//PreWriteResponsePlugin represents .
	PreWriteResponsePlugin interface {
		PreWriteResponse(context.Context, protocol.Message, protocol.Message) error
	}

	//PostWriteResponsePlugin represents .
	PostWriteResponsePlugin interface {
		PostWriteResponse(context.Context, protocol.Message, protocol.Message, error) error
	}

	//PreWriteRequestPlugin represents .
	PreWriteRequestPlugin interface {
		PreWriteRequest(ctx context.Context) error
	}

	//PostWriteRequestPlugin represents .
	PostWriteRequestPlugin interface {
		PostWriteRequest(ctx context.Context, r protocol.Message, e error) error
	}

	// http
	PreHandlerHTTPRequestPlugin interface {
		PreHandlerHTTPRequest(w http.ResponseWriter, r *http.Request) error
	}

	PostHandlerHTTPRequestPlugin interface {
		PostHandlerHTTPRequest(w http.ResponseWriter, r *http.Request) error
	}
)

// pluginContainer implements PluginContainer interface.
type PluginContainer struct {
	plugins []Plugin
}

// Add adds a plugin.
func (p *PluginContainer) Add(plugin Plugin) {
	p.plugins = append(p.plugins, plugin)
}

// Remove removes a plugin by it's name.
func (p *PluginContainer) Remove(plugin Plugin) {
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

func (p *PluginContainer) All() []Plugin {
	return p.plugins
}

// DoRegister invokes DoRegister plugin.
func (p *PluginContainer) DoRegister(name string, rcvr interface{}, metadata string) error {
	var es []error
	for _, rp := range p.plugins {
		if plugin, ok := rp.(RegisterPlugin); ok {
			err := plugin.Register(name, rcvr, metadata)
			if err != nil {
				es = append(es, err)
			}
		}
	}

	if len(es) > 0 {
		return errors.NewMultiError(es)
	}
	return nil
}

// DoRegisterFunction invokes DoRegisterFunction plugin.
func (p *PluginContainer) DoRegisterFunction(name string, fn interface{}, metadata string) error {
	var es []error
	for _, rp := range p.plugins {
		if plugin, ok := rp.(RegisterFunctionPlugin); ok {
			err := plugin.RegisterFunction(name, fn, metadata)
			if err != nil {
				es = append(es, err)
			}
		}
	}

	if len(es) > 0 {
		return errors.NewMultiError(es)
	}
	return nil
}

// DoUnregister invokes RegisterPlugin.
func (p *PluginContainer) DoUnregister(name string) error {
	var es []error
	for _, rp := range p.plugins {
		if plugin, ok := rp.(RegisterPlugin); ok {
			err := plugin.Unregister(name)
			if err != nil {
				es = append(es, err)
			}
		}
	}

	if len(es) > 0 {
		return errors.NewMultiError(es)
	}
	return nil
}

//DoPostConnAccept handles accepted conn
func (p *PluginContainer) DoPostConnAccept(conn net.Conn) (net.Conn, bool) {
	var flag bool
	for i := range p.plugins {
		if plugin, ok := p.plugins[i].(PostConnAcceptPlugin); ok {
			conn, flag = plugin.HandleConnAccept(conn)
			if !flag { //interrupt
				conn.Close()
				return conn, false
			}
		}
	}
	return conn, true
}

//DoPostConnClose handles closed conn
func (p *PluginContainer) DoPostConnClose(conn net.Conn) bool {
	var flag bool
	for i := range p.plugins {
		if plugin, ok := p.plugins[i].(PostConnClosePlugin); ok {
			flag = plugin.HandleConnClose(conn)
			if !flag {
				return false
			}
		}
	}
	return true
}

// DoPreReadRequest invokes PreReadRequest plugin.
func (p *PluginContainer) DoPreReadRequest(ctx context.Context) error {
	for i := range p.plugins {
		if plugin, ok := p.plugins[i].(PreReadRequestPlugin); ok {
			err := plugin.PreReadRequest(ctx)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// DoPostReadRequest invokes PostReadRequest plugin.
func (p *PluginContainer) DoPostReadRequest(ctx context.Context, r protocol.Message, e error) error {
	for i := range p.plugins {
		if plugin, ok := p.plugins[i].(PostReadRequestPlugin); ok {
			err := plugin.PostReadRequest(ctx, r, e)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// DoPreHandleRequest invokes PreHandleRequest plugin.
func (p *PluginContainer) DoPreHandleRequest(ctx context.Context, r protocol.Message) error {
	for i := range p.plugins {
		if plugin, ok := p.plugins[i].(PreHandleRequestPlugin); ok {
			err := plugin.PreHandleRequest(ctx, r)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// DoPreWriteResponse invokes PreWriteResponse plugin.
func (p *PluginContainer) DoPreWriteResponse(ctx context.Context, req protocol.Message, res protocol.Message) error {
	for i := range p.plugins {
		if plugin, ok := p.plugins[i].(PreWriteResponsePlugin); ok {
			err := plugin.PreWriteResponse(ctx, req, res)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// DoPostWriteResponse invokes PostWriteResponse plugin.
func (p *PluginContainer) DoPostWriteResponse(ctx context.Context, req protocol.Message, resp protocol.Message, e error) error {
	for i := range p.plugins {
		if plugin, ok := p.plugins[i].(PostWriteResponsePlugin); ok {
			err := plugin.PostWriteResponse(ctx, req, resp, e)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// DoPreWriteRequest invokes PreWriteRequest plugin.
func (p *PluginContainer) DoPreWriteRequest(ctx context.Context) error {
	for i := range p.plugins {
		if plugin, ok := p.plugins[i].(PreWriteRequestPlugin); ok {
			err := plugin.PreWriteRequest(ctx)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// DoPostWriteRequest invokes PostWriteRequest plugin.
func (p *PluginContainer) DoPostWriteRequest(ctx context.Context, r protocol.Message, e error) error {
	for i := range p.plugins {
		if plugin, ok := p.plugins[i].(PostWriteRequestPlugin); ok {
			err := plugin.PostWriteRequest(ctx, r, e)
			if err != nil {
				return err
			}
		}
	}

	return nil
}


func (p *PluginContainer) DoPreHandlerHTTPRequest(w http.ResponseWriter, r *http.Request) error {
	for i := range p.plugins {
		if plugin, ok := p.plugins[i].(PreHandlerHTTPRequestPlugin); ok {
			err := plugin.PreHandlerHTTPRequest(w, r)
			if err != nil {
				return err
			}
		}
	}

	return nil

}
func (p *PluginContainer) DoPostHandlerHTTPRequest(w http.ResponseWriter, r *http.Request) error {
	for i := range p.plugins {
		if plugin, ok := p.plugins[i].(PostHandlerHTTPRequestPlugin); ok {
			err := plugin.PostHandlerHTTPRequest(w, r)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
