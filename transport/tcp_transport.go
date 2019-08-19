package transport

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	log "github.com/hillguo/sanlog"
	"github.com/hillguo/sanrpc/protocol"
	"github.com/hillguo/sanrpc/protocol/sanrpc"
	"io"
	"net"
	"runtime"
	"strings"
	"sync"
	"time"
)

const ReaderBuffsize = 1024

func NewTCPTransport(opt ...TransportOption) ServerTransport {

	opts := &TransportOptions{
		InMsgChanSize:  10,
		OutMsgChanSize: 10,
		MsgProtocol:    &sanrpc.SanRPCProtocol{},
	}

	for _, o := range opt {
		o(opts)
	}
	return &tcpTransport{
		opts: opts,
	}
}

type tcpTransport struct {
	opts *TransportOptions

	mu         sync.RWMutex
	activeConn map[net.Conn]struct{}
	doneChan   chan struct{}
}


func (t *tcpTransport) Close() error {
	return nil
}

func (t *tcpTransport) ListenAndServer(opt ...TransportOption) error {
	for _, o := range opt {
		o(t.opts)
	}

	if t.opts.MsgProtocol == nil {
		t.opts.MsgProtocol = &sanrpc.SanRPCProtocol{}
	}

	switch t.opts.NetWork {
	case "tcp", "tcp4", "tcp6":
		ln, err := net.Listen(t.opts.NetWork, t.opts.Address)
		if err != nil {
			log.Error(err)
			return err
		}
		log.Infof("sanrpc listening network:%s ,address:%s", t.opts.NetWork, t.opts.Address)
		return t.server(ln)
	default:
		return fmt.Errorf("transport: not support network type %s", t.opts.NetWork)
	}
}

func (t *tcpTransport) server(ln net.Listener) error {
	var tempDelay time.Duration

	t.mu.Lock()
	if t.activeConn == nil {
		t.activeConn = make(map[net.Conn]struct{})
	}
	t.mu.Unlock()

	for {
		conn, e := ln.Accept()
		if e != nil {
			if ne, ok := e.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}

				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}

				log.Errorf("sanrpc: Accept error: %v; retrying in %v", e, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return e
		}
		tempDelay = 0

		if tc, ok := conn.(*net.TCPConn); ok {
			tc.SetKeepAlive(true)
			tc.SetKeepAlivePeriod(3 * time.Minute)
			tc.SetLinger(10)
		}

		t.mu.Lock()
		t.activeConn[conn] = struct{}{}
		t.mu.Unlock()

		log.Infof("rpc: receive a client conn, remote addr: %+v", conn.RemoteAddr())

		go t.serveConn(conn)
	}
}

func (t *tcpTransport) serveConn(conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			ss := runtime.Stack(buf, false)
			if ss > size {
				ss = size
			}
			buf = buf[:ss]
			log.Errorf("serving %s panic error: %s, stack:\n %s", conn.RemoteAddr(), err, buf)
		}
		t.mu.Lock()
		delete(t.activeConn, conn)
		t.mu.Unlock()
		conn.Close()

	}()

	if tlsConn, ok := conn.(*tls.Conn); ok {
		if d := t.opts.ReadTimeout; d != 0 {
			conn.SetReadDeadline(time.Now().Add(time.Duration(d)))
		}
		if d := t.opts.WriteTimeout; d != 0 {
			conn.SetWriteDeadline(time.Now().Add(time.Duration(d)))
		}
		if err := tlsConn.Handshake(); err != nil {
			log.Errorf("sanrpc: TLS handshake error from %s: %v", conn.RemoteAddr(), err)
			return
		}
		log.Infof("sanrpc: TLS handshake success")
	}

	ctx := context.Background()
	ctx, cancelCtx := context.WithCancel(ctx)

	in := make(chan protocol.Message, t.opts.InMsgChanSize)
	out := make(chan protocol.Message, t.opts.OutMsgChanSize)

	r := bufio.NewReaderSize(conn, ReaderBuffsize)

	var wg sync.WaitGroup
	wg.Add(3)

	// 1. read request msg
	{
		go func(ctx context.Context, in chan<- protocol.Message) {
			defer func() {
				if err := recover(); err != nil {
					const size = 64 << 10
					buf := make([]byte, size)
					ss := runtime.Stack(buf, false)
					if ss > size {
						ss = size
					}
					buf = buf[:ss]
					log.Errorf("serving %s panic error: %s, stack:\n %s", conn.RemoteAddr(), err, buf)
				}
				wg.Done()
				cancelCtx()
			}()

			for {

				t0 := time.Now()
				if t.opts.ReadTimeout != 0 {
					conn.SetReadDeadline(t0.Add(time.Duration(t.opts.ReadTimeout)))
				}
				req, err := t.opts.MsgProtocol.DecodeMessage(r)
				if err != nil {
					if err == io.EOF {
						log.Infof("client has closed this connection: %s", conn.RemoteAddr().String())
					} else if strings.Contains(err.Error(), "use of closed network connection") {
						log.Infof("sanrpc: connection %s is closed", conn.RemoteAddr().String())
					} else {
						log.Errorf("sanrpc: failed to read request: %v", err)
					}
					log.Infof("connection: %s read routine context done", conn.RemoteAddr().String())
					return
				}

				log.Infof("read a message from conn %v", conn.RemoteAddr())

				select {
				case <-ctx.Done():
					log.Infof("connection: %s read routine context done %v", conn.RemoteAddr().String(), ctx.Err())
					return
				case in <- req:
					log.Infof("put message into in queue")
				}
			}
		}(ctx, in)
	}

	// 2. handler request msg
	{
		go func(ctx context.Context, in <-chan protocol.Message, out chan<- protocol.Message) {
			defer func() {
				if err := recover(); err != nil {
					const size = 64 << 10
					buf := make([]byte, size)
					ss := runtime.Stack(buf, false)
					if ss > size {
						ss = size
					}
					buf = buf[:ss]
					log.Errorf("serving %s panic error: %s, stack:\n %s", conn.RemoteAddr(), err, buf)
				}
				wg.Done()
				cancelCtx()
			}()
			for {
				select {
				case <-ctx.Done():
					log.Infof("connection: %s handle routine context done %v", conn.RemoteAddr().String(), ctx.Err())
					return
				case req := <-in:
					go func() {
						defer func() {
							if err := recover(); err != nil {
								const size = 64 << 10
								buf := make([]byte, size)
								ss := runtime.Stack(buf, false)
								if ss > size {
									ss = size
								}
								buf = buf[:ss]
								log.Errorf("serving %s panic error: %s, stack:\n %s", conn.RemoteAddr(), err, buf)
								cancelCtx()
							}
						}()
						log.Infof("get a message from in queue to handler")
						resp, err := t.opts.MsgProtocol.HandleMessage(ctx, req)
						if err != nil {
							log.Warnf("rpc: failed to handle request: %v", err)
						}
						out <- resp
						log.Infof("handler message over, put it to out queue")
					}()
				}
			}
		}(ctx, in, out)
	}

	// 3. write response msg
	{
		go func(ctx context.Context, ch <-chan protocol.Message) {
			defer func() {
				if err := recover(); err != nil {
					const size = 64 << 10
					buf := make([]byte, size)
					ss := runtime.Stack(buf, false)
					if ss > size {
						ss = size
					}
					buf = buf[:ss]
					log.Errorf("serving %s panic error: %s, stack:\n %s", conn.RemoteAddr(), err, buf)
				}
				wg.Done()
				cancelCtx()
			}()

			for {
				select {
				case <-ctx.Done():
					log.Infof("connection: %s write routine context done %v", conn.RemoteAddr().String(), ctx.Err())
					return
				case resp := <-out:
					log.Infof("read a resp message form out queue")
					data,err := t.opts.MsgProtocol.EncodeMessage(resp)
					if err != nil{
						log.Error(err)
						return
					}
					log.Infof("rpc: encode resp , writr into conn")
					_, err = conn.Write(data)
					if err != nil {
						log.Error("connection: %s write routine context done %v", conn.RemoteAddr().String(), err)
						return
					}
				}
			}
		}(ctx, out)
	}
	wg.Wait()
	log.Infof("connection %s destroyed", conn.RemoteAddr().String())
}
