package client

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/hillguo/sanrpc/client/selector"
	"github.com/hillguo/sanrpc/client/servicediscovery"
	"github.com/hillguo/sanrpc/share"
)

var (
	// ErrXClientShutdown xclient is shutdown.
	ErrXClientShutdown = errors.New("xClient is shut down")
	// ErrXClientNoServer selector can't found one server.
	ErrXClientNoServer = errors.New("can not found any server")
	// ErrServerUnavailable selected server is unavailable.
	ErrServerUnavailable = errors.New("selected server is unavilable")
)

type XClient interface {
	Go(ctx context.Context, serviceMethod string, args interface{}, reply interface{}, done chan *Call) (*Call, error)
	Call(ctx context.Context, serviceMethod string, args interface{}, reply interface{}) error
	Close() error
}

type xClient struct {
	failMode     FailMode
	selectMode   selector.SelectMode
	cachedClient map[string]RPCClient
	servicePath  string
	option       Option

	mu         sync.RWMutex
	servers    map[string]string
	discovery  servicediscovery.ServiceDiscovery
	selector   selector.Selector
	isShutdown bool
	auth       string
	Plugins    PluginContainer
	ch         chan []*servicediscovery.KVPair
}

type FailMode int

const (
	Failover FailMode = iota
	Failfast
	Failtry
	Failbackup
)

func NewXClient(servicePath string, failMode FailMode, selectMode selector.SelectMode,
	discovery servicediscovery.ServiceDiscovery, option Option) XClient {
	client := &xClient{
		failMode:     failMode,
		selectMode:   selectMode,
		discovery:    discovery,
		servicePath:  servicePath,
		cachedClient: make(map[string]RPCClient),
		option:       option,
	}
	servers := make(map[string]string)
	pairs := discovery.GetServices()
	for _, p := range pairs {
		servers[p.Key] = p.Value
	}
	client.servers = servers

	client.Plugins = &pluginContainer{}
	ch := client.discovery.WatchService()
	if ch != nil {
		client.ch = ch
		go client.watch(ch)
	}
	return client
}

func (c *xClient) watch(ch chan []*servicediscovery.KVPair) {
	for pairs := range ch {
		servers := make(map[string]string)
		for _, p := range pairs {
			servers[p.Key] = p.Value
		}
		c.mu.Lock()
		c.servers = servers
		if c.selector != nil {
			c.selector.UpdateServer(servers)
		}
		c.mu.Unlock()
	}
}

// SetSelector sets customized selector by users.
func (c *xClient) SetSelector(s selector.Selector) {
	c.mu.RLock()
	s.UpdateServer(c.servers)
	c.mu.RUnlock()

	c.selector = s
}

// SetPlugins sets client's plugins.
func (c *xClient) SetPlugins(plugins PluginContainer) {
	c.Plugins = plugins
}

func (c *xClient) GetPlugins() PluginContainer {
	return c.Plugins
}

// Auth sets s token for Authentication.
func (c *xClient) Auth(auth string) {
	c.auth = auth
}

// selects a client from candidates base on c.selectMode
func (c *xClient) selectClient(ctx context.Context, servicePath, serviceMethod string, args interface{}) (string, RPCClient, error) {
	c.mu.Lock()
	k := c.selector.Select(ctx, servicePath, serviceMethod, args)
	c.mu.Unlock()
	if k == "" {
		return "", nil, ErrXClientNoServer
	}
	client, err := c.getCachedClient(k)
	return k, client, err
}

func (c *xClient) getCachedClient(k string) (RPCClient, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	client := c.cachedClient[k]
	if client != nil {
		if !client.IsClosing() && !client.IsShutdown() {
			return client, nil
		}
		delete(c.cachedClient, k)
		client.Close()
	}
	if client == nil || client.IsShutdown() {
		network, addr := splitNetworkAndAddress(k)
		client = &Client{
			option:  c.option,
			Plugins: c.Plugins,
		}
		err := client.Connect(network, addr)
		if err != nil {
			return nil, err
		}
		c.cachedClient[k] = client
	}
	return client, nil
}

func (c *xClient) removeClient(k string, client RPCClient) {
	c.mu.Lock()
	cl := c.cachedClient[k]
	if cl == client {
		delete(c.cachedClient, k)
	}
	c.mu.Unlock()

	if client != nil {
		client.Close()
	}
}

func splitNetworkAndAddress(server string) (string, string) {
	ss := strings.SplitN(server, "@", 2)
	if len(ss) == 1 {
		return "tcp", server
	}

	return ss[0], ss[1]
}

func (c *xClient) Go(ctx context.Context, serviceMethod string, args interface{}, reply interface{}, done chan *Call) (*Call, error) {
	if c.isShutdown {
		return nil, ErrXClientShutdown
	}

	if c.auth != "" {
		metadata := ctx.Value(share.ReqMetaDataKey)
		if metadata == nil {
			metadata = map[string]string{}
			ctx = context.WithValue(ctx, share.ReqMetaDataKey, metadata)
		}
		m := metadata.(map[string]string)
		m[share.AuthKey] = c.auth
	}

	_, client, err := c.selectClient(ctx, c.servicePath, serviceMethod, args)
	if err != nil {
		return nil, err
	}
	return client.Go(ctx, c.servicePath, serviceMethod, args, reply, done), nil
}

func (c *xClient) Call(ctx context.Context, serviceMethod string, args interface{}, reply interface{}) error {
	if c.isShutdown {
		return ErrXClientShutdown
	}

	if c.auth != "" {
		metadata := ctx.Value(share.ReqMetaDataKey)
		if metadata == nil {
			metadata = map[string]string{}
			ctx = context.WithValue(ctx, share.ReqMetaDataKey, metadata)
		}
		m := metadata.(map[string]string)
		m[share.AuthKey] = c.auth
	}

	var err error
	k, client, err := c.selectClient(ctx, c.servicePath, serviceMethod, args)
	if err != nil {
		if c.failMode == Failfast {
			return err
		}
	}

	var e error
	switch c.failMode {
	case Failtry:
		retries := c.option.Retries
		for retries > 0 {
			retries--

			if client != nil {
				err = client.Call(ctx, c.servicePath, serviceMethod, args, reply)
				if err == nil {
					return nil
				}
				if _, ok := err.(ServiceError); ok {
					return err
				}
			}

			c.removeClient(k, client)
			client, e = c.getCachedClient(k)
		}
		if err == nil {
			err = e
		}
		return err
	case Failover:
		retries := c.option.Retries
		for retries > 0 {
			retries--

			if client != nil {
				err = client.Call(ctx, c.servicePath, serviceMethod, args, reply)
				if err == nil {
					return nil
				}
				if _, ok := err.(ServiceError); ok {
					return err
				}
			}

			c.removeClient(k, client)
			//select another server
			k, client, e = c.selectClient(ctx, c.servicePath, serviceMethod, args)
		}

		if err == nil {
			err = e
		}
		return err
	default: //Failfast
		err = client.Call(ctx, c.servicePath, serviceMethod, args, reply)
		if err != nil {
			if _, ok := err.(ServiceError); !ok {
				c.removeClient(k, client)
			}
		}

		return err
	}
}

// Close closes this client and its underlying connnections to services.
func (c *xClient) Close() error {
	c.isShutdown = true

	var errs []error
	c.mu.Lock()
	for k, v := range c.cachedClient {
		e := v.Close()
		if e != nil {
			errs = append(errs, e)
		}

		delete(c.cachedClient, k)

	}
	c.mu.Unlock()

	go func() {
		defer func() {
			if r := recover(); r != nil {

			}
		}()

		c.discovery.RemoveWatcher(c.ch)
		close(c.ch)
	}()

	if len(errs) > 0 {
		return fmt.Errorf("%v", errs)
	}
	return nil
}
