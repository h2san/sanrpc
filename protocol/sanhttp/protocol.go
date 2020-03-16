// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sanhttp

import (
	"fmt"
	"github.com/hillguo/sanrpc/protocol/sanhttp/ctx"
	"github.com/hillguo/sanrpc/errs"
	"github.com/hillguo/sanrpc/protocol/sanhttp/middleware"
	"github.com/hillguo/sanrpc/protocol/sanhttp/route"
	"log"
	"net"
	"net/http"
	"os"
)

const defaultMultipartMemory = 32 << 20 // 32 MB

var (
	default404Body = "404 page not found"
)

func noRoute(c *ctx.Context) {
	c.String(404, default404Body)
}

// server is the framework's instance, it contains the muxer, middleware and configuration settings.
// Create an instance of server, by using New() or Default()
type HttpServer struct {
	RouterGroup

	route *route.Route

	// Value of 'maxMemory' param that is given to http.Request's ParseMultipartForm
	// method call.
	MaxMultipartMemory int64

	noRoute  ctx.HandlersChain
	noMethod ctx.HandlersChain
}

func New() *HttpServer {
	server := &HttpServer{
		RouterGroup: RouterGroup{
			Handlers: nil,
			basePath: "/",
			root:     true,
		},
		MaxMultipartMemory: defaultMultipartMemory,
		route:              route.DefaultRoute,
		noRoute:            ctx.HandlersChain{noRoute},
	}
	server.RouterGroup.server = server
	return server
}

// Default returns an server instance with the Logger and Recovery middleware already attached.
func Default() *HttpServer {
	server := New()
	server.Use(middleware.Logger(), middleware.Recovery())
	return server
}

// NoRoute adds handlers for NoRoute. It return a 404 code by default.
func (server *HttpServer) NoRoute(handlers ...ctx.HandlerFunc) {
	server.noRoute = handlers
	server.rebuild404Handlers()
}

// NoMethod sets the handlers called when... TODO.
func (server *HttpServer) NoMethod(handlers ...ctx.HandlerFunc) {
	server.noMethod = handlers
	server.rebuild405Handlers()
}

// Use attaches a global middleware to the router. ie. the middleware attached though Use() will be
// included in the handlers chain for every single request. Even 404, 405, static files...
// For example, this is the right place for a logger or error management middleware.
func (server *HttpServer) Use(middleware ...ctx.HandlerFunc) IRoutes {
	server.RouterGroup.Use(middleware...)
	server.rebuild404Handlers()
	server.rebuild405Handlers()
	return server
}

func (server *HttpServer) rebuild404Handlers() {
	server.noRoute = server.combineHandlers(server.noRoute)
}

func (server *HttpServer) rebuild405Handlers() {
	server.noMethod = server.combineHandlers(server.noMethod)
}

func (server *HttpServer) addRoute(method, path string, handlers ctx.HandlersChain) {
	server.route.AddHandler(method, path, handlers)
}

// Run attaches the router to a http.Server and starts listening and serving HTTP requests.
// It is a shortcut for http.ListenAndServe(addr, router)
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func (server *HttpServer) Run(addr ...string) (err error) {
	defer func() {}()

	address := ResolveAddress(addr)
	log.Printf("Listening and serving HTTP on %s\n", address)
	err = http.ListenAndServe(address, server)
	return
}

// RunTLS attaches the router to a http.Server and starts listening and serving HTTPS (secure) requests.
// It is a shortcut for http.ListenAndServeTLS(addr, certFile, keyFile, router)
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func (server *HttpServer) RunTLS(addr, certFile, keyFile string) (err error) {
	log.Fatalf("Listening and serving HTTPS on %s\n", addr)
	defer func() { log.Println(err) }()

	err = http.ListenAndServeTLS(addr, certFile, keyFile, server)
	return
}

// RunUnix attaches the router to a http.Server and starts listening and serving HTTP requests
// through the specified unix socket (ie. a file).
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func (server *HttpServer) RunUnix(file string) (err error) {
	log.Printf("Listening and serving HTTP on unix:/%s", file)
	defer func() { log.Println(err) }()

	os.Remove(file)
	listener, err := net.Listen("unix", file)
	if err != nil {
		return
	}
	defer listener.Close()
	os.Chmod(file, 0777)
	err = http.Serve(listener, server)
	return
}

// RunFd attaches the router to a http.Server and starts listening and serving HTTP requests
// through the specified file descriptor.
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func (server *HttpServer) RunFd(fd int) (err error) {
	log.Printf("Listening and serving HTTP on fd@%d", fd)
	defer func() { log.Println(err) }()

	f := os.NewFile(uintptr(fd), fmt.Sprintf("fd@%d", fd))
	listener, err := net.FileListener(f)
	if err != nil {
		return
	}
	defer listener.Close()
	err = http.Serve(listener, server)
	return
}

// ServeHTTP conforms to the http.Handler interface.
func (server *HttpServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := &ctx.Context{
		Index: -1,
		Err:   &errs.Error{},
	}

	c.Request = req
	c.Writer = w

	server.handleHTTPRequest(c)

}

func (server *HttpServer) handleHTTPRequest(c *ctx.Context) {
	httpMethod := c.Request.Method
	rPath := c.Request.URL.Path

	handlers := server.route.GetHandler(httpMethod, rPath)
	if handlers == nil {
		c.Handlers = server.noRoute
	} else {
		c.Handlers = handlers
	}
	c.Next()
}

func ResolveAddress(addr []string) string {
	switch len(addr) {
	case 0:
		if port := os.Getenv("PORT"); port != "" {
			log.Printf("Environment variable PORT=\"%s\"", port)
			return ":" + port
		}
		log.Println("Environment variable PORT is undefined. Using port :8080 by default")
		return ":8080"
	case 1:
		return addr[0]
	default:
		panic("too many parameters")
	}
}
