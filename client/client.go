package client

import (
	"context"
	"crypto/tls"
	"errors"
	log "github.com/hillguo/sanlog"
	"github.com/hillguo/sanrpc/client/discovery"
	"github.com/hillguo/sanrpc/client/node"
	"github.com/hillguo/sanrpc/codec"
	"github.com/hillguo/sanrpc/protocol/sanrpc"
	"net"
	"time"
)
import "github.com/hillguo/sanrpc/client/selector"

type Client struct {
	opts Options
}

func NewClient(opt ...Option) *Client{
	opts:= &Options{
		ServiceName: "11",
		ConnectTimeout:0,
		ReadTimeout: 0,
		ConnectTimeout:0,
	}

	for _,o :=range opt {
		o(opts)
	}
	return &Client{
		opts: opts,
	}
}

func (c *Client) Invoke(ctx context.Context, req interface{}, resp interface{}) error {
	// 1. select
	node, err := c.selectNode()
	if err != nil {
		return err
	}

	conn ,err := c.Connect(node)
	if err != nil {
		return err
	}

	reqmsg := &sanrpc.MessageProtocol{}
	reqmsg.Header = &sanrpc.HeaderMsg{}
	reqmsg.Header.CallType = uint32(sanrpc.SanrpcMsgType_SANRPC_REQUEST_MSG)
	reqmsg.Header.Seq = 1
	reqmsg.Header.EncodeType = uint32(codec.ProtoBuffer)
	reqmsg.Header.MetaData = nil
	// msg.Header.ServiceName = call.ServicePath
	// msg.Header.MethodName = call.ServiceMethod
	reqmsg.Header.CompressType= uint32(codec.CompressNone)


	cc := codec.Codecs[codec.ProtoBuffer]
	if cc == nil {
		return errors.New("no codec")
	}
	data, err := cc.Encode(req)
	if err != nil {
		return err
	}
	reqmsg.Data = data
	log.Debugf("req msg %+v", req)
	msgProtocol := &sanrpc.SanRPCProtocol{}
	d, err := msgProtocol.EncodeMessage(req)
	if c.opts.WriteTimeout != 0 {
		_ = conn.SetWriteDeadline(time.Now().Add(c.opts.WriteTimeout))
	}
	_, err = conn.Write(d)
	if err != nil {
		return err
	}
	if c.opts.ReadTimeout != 0 {
		_ = conn.SetReadDeadline(time.Now().Add(c.opts.ReadTimeout))
	}
	respmsg, err := msgProtocol.DecodeMessage(conn)
	if err != nil {
		log.Debug("DecodeMessage", err)
		return err
	}
	res, _ := respmsg.(*sanrpc.MessageProtocol)
	log.Debugf("resp msg %+v", res)
	if res.Header == nil {
		res.Header = &sanrpc.HeaderMsg{}
	}










}

func (c *Client) Connect(node *node.Node) ( net.Conn ,error){
	var conn net.Conn
	var err error
	conn, err = net.DialTimeout(node.Network, node.Address, c.opts.ConnectTimeout)
	if err != nil {
		log.Warnf("failed to dial server: %v", err)
		return nil, err
	}
	if tc, ok := conn.(*net.TCPConn); ok {
		tc.SetKeepAlive(true)
		tc.SetKeepAlivePeriod(3 * time.Minute)
	}
	return conn, nil
}

func (c *Client) selectNode() (*node.Node,error) {

	dis := discovery.Get(c.opts.Discovery)
	if dis == nil {
		dis = discovery.DefaultDiscovery
	}

	serviceName := c.opts.ServiceName
	nodes,err := dis.List(serviceName)
	if err != nil {
		return nil,err
	}

	sec:= selector.GetSelector(c.opts.Selector)
	if sec == nil {
		sec=selector.DefaultSelector
	}
	node,err := sec.Select(nodes)
	if err != nil {
		return nil,err
	}
	return node, nil
}
