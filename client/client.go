package client

import (
	"context"
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
	opts *Options
}

func NewClient(opt ...Option) *Client{
	opts:= &Options{
		Network: "tcp",
		Discovery: "ip",
		Selector: "random",
		ConnectTimeout:0,
		ReadTimeout: 0,
		WriteTimeout:0,
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
		log.Error("select node error", err)
		return err
	}
	node.Network = c.opts.Network
	log.Debugf("select node success. node: %+v", node)

	conn ,err := c.Connect(node)
	if err != nil {
		return err
	}
	defer func() {
		conn.Close()
	}()

	reqmsg := &sanrpc.MessageProtocol{}
	reqmsg.Header = &sanrpc.HeaderMsg{
		Version: 0,
		CallType: uint32(sanrpc.SanrpcMsgType_SANRPC_REQUEST_MSG),
		ServiceName:c.opts.ServiceName,
		MethodName: c.opts.MethodName,
		EncodeType: uint32(codec.ProtoBuffer),
		CompressType: uint32(codec.CompressNone),
		MetaData: nil,
	}
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
	d, err := msgProtocol.EncodeMessage(reqmsg)
	if err != nil {
		return err
	}
	if c.opts.WriteTimeout != 0 {
		_ = conn.SetWriteDeadline(time.Now().Add(time.Duration(c.opts.WriteTimeout)))
	}
	_, err = conn.Write(d)
	if err != nil {
		return err
	}
	if c.opts.ReadTimeout != 0 {
		_ = conn.SetReadDeadline(time.Now().Add(time.Duration(c.opts.ReadTimeout)))
	}
	respmsg, err := msgProtocol.DecodeMessage(conn)
	if err != nil {
		log.Debug("DecodeMessage", err)
		return err
	}
	res, _ := respmsg.(*sanrpc.MessageProtocol)
	log.Debugf("resp msg %+v", res)
	if res.Header == nil {
		return errors.New("resp header nil")
	}
	cc = codec.Codecs[codec.SerializeType(res.Header.EncodeType)]
	if cc == nil {
		return errors.New("resp codec not support")
	}
	err = cc.Decode(res.Data, resp)
	if err != nil {
		log.Error("data decode fail")
		return errors.New("data decode fail")
	}
	return nil
}

func (c *Client) Connect(node *node.Node) ( net.Conn ,error){
	var conn net.Conn
	var err error
	conn, err = net.DialTimeout(node.Network, node.Address, time.Duration(c.opts.ConnectTimeout))
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
		log.Debugf("can't find assign discovery [%s]. use default [%s]",c.opts.Discovery, dis.Name())
	}

	address := c.opts.Address
	nodes,err := dis.List(address)
	if err != nil {
		return nil,err
	}

	sec:= selector.GetSelector(c.opts.Selector)
	if sec == nil {
		sec=selector.DefaultSelector
		log.Debugf("can't find assign selector [%s]. use default [%s]",c.opts.Selector, sec.Name())
	}
	node,err := sec.Select(nodes)
	if err != nil {
		return nil,err
	}
	return node, nil
}
