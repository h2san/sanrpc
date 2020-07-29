package client

import (
	"bufio"
	"context"
	"errors"
	"io"
	"net"
	"sync"
	"time"

	"github.com/hillguo/sanrpc/errs"

	log "github.com/hillguo/sanlog"
	"github.com/hillguo/sanrpc/codec"
	"github.com/hillguo/sanrpc/protocol/sanrpc"
)

// ErrShutdown connection is closed.
var (
	ErrShutdown         = errors.New("connection is shut down")
	ErrUnsupportedCodec = errors.New("unsupported codec")
)

// ServiceError is an error from server.
type ServiceError string

func (e ServiceError) Error() string {
	return string(e)
}

type seqKey struct{}

// RPCClient is interface that defines one client to call one server.
type RPCClient interface {
	Connect(network, address string) error
	Go(ctx context.Context, servicePath, serviceMethod string, args interface{}, reply interface{}, done chan *Call) *Call
	Call(ctx context.Context, servicePath, serviceMethod string, args interface{}, reply interface{}) error
	Close() error
	IsClosing() bool
	IsShutdown() bool
}
type Client struct {
	option Option
	Conn   net.Conn
	r      *bufio.Reader

	mutex    sync.Mutex
	seq      uint64
	pending  map[uint64]*Call
	closing  bool
	shutdown bool
}

func NewClient(option Option) *Client {
	return &Client{
		option: option,
	}
}

func (client *Client) IsClosing() bool {
	return client.closing
}

func (client *Client) IsShutdown() bool {
	return client.shutdown
}
func (client *Client) Go(ctx context.Context, servicePath, serviceMethod string, args interface{}, reply interface{}, done chan *Call) *Call {
	call := new(Call)
	call.ServicePath = servicePath
	call.ServiceMethod = serviceMethod
	meta := ctx.Value(ReqMetaDataKey)
	if meta != nil {
		call.Metadata = meta.(map[string]string)
	}

	call.Args = args
	call.Reply = reply
	call.CompressType = codec.CompressNone
	call.SerializeType = codec.ProtoBuffer

	if done == nil {
		done = make(chan *Call, 10)
	} else {
		if cap(done) == 0 {
			log.Panic("rpc: done channel is unbuffered")
		}
	}
	call.Done = done
	client.send(ctx, call)
	return call
}

//Call 同步调用
func (client *Client) Call(ctx context.Context, servicePath, serviceMethod string,
	args interface{}, reply interface{}) error {
	if client.Conn == nil {
		return errors.New("conn not establish")
	}

	seq := new(uint64)
	ctx = context.WithValue(ctx, seqKey{}, seq)
	log.Debugf("req %+v", args)
	DoneChan := client.Go(ctx, servicePath, serviceMethod, args, reply, make(chan *Call, 1)).Done
	var err error
	select {
	case <-ctx.Done():
		client.mutex.Lock()
		call := client.pending[*seq]
		delete(client.pending, *seq)
		client.mutex.Unlock()
		if call != nil {
			call.Error = ctx.Err()
			call.done()
		}
		err = ctx.Err()
	case call := <-DoneChan:
		err = call.Error
		meta := ctx.Value(ResMetaDataKey)
		if meta != nil && len(call.ResMetadata) > 0 {
			resMata := meta.(map[string]string)
			for k, v := range call.ResMetadata {
				resMata[k] = v
			}
		}
	}
	log.Debugf("resp %+v", reply)
	return err
}

func (client *Client) send(ctx context.Context, call *Call) {
	client.mutex.Lock()
	if client.shutdown || client.closing {
		call.Error = ErrShutdown
		client.mutex.Unlock()
		call.done()
		return
	}

	cc := codec.Codecs[call.SerializeType]
	if cc == nil {
		call.Error = ErrUnsupportedCodec
		client.mutex.Unlock()
		call.done()
		return
	}
	if client.pending == nil {
		client.pending = make(map[uint64]*Call)
	}
	seq := client.seq
	client.seq++
	client.pending[seq] = call
	client.mutex.Unlock()

	if cseq, ok := ctx.Value(seqKey{}).(*uint64); ok {
		*cseq = seq
	}

	req := &sanrpc.MessageProtocol{}
	req.Header = &sanrpc.HeaderMsg{}
	req.Header.CallType = uint32(sanrpc.SanrpcMsgType_SANRPC_REQUEST_MSG)
	req.Header.Seq = seq

	req.Header.EncodeType = uint32(call.SerializeType)
	if call.Metadata != nil {
		req.Header.MetaData = call.Metadata
	}
	req.Header.ServiceName = call.ServicePath
	req.Header.MethodName = call.ServiceMethod

	data, err := cc.Encode(call.Args)
	if err != nil {
		call.Error = err
		call.done()
		return
	}
	if len(data) > 1024 && call.CompressType != codec.CompressNone {
		req.Header.CompressType = uint32(call.CompressType)
	}
	req.Data = data

	log.Debugf("req msg %+v", req)
	msg := sanrpc.SanRPCProtocol{}
	d, err := msg.EncodeMessage(req)

	if client.option.WriteTimeout != 0 {
		_ = client.Conn.SetWriteDeadline(time.Now().Add(client.option.WriteTimeout))
	}
	if client.Conn == nil {
		client.mutex.Lock()
		call = client.pending[seq]
		delete(client.pending, seq)
		client.mutex.Unlock()
		call.Error = errors.New("conn not establish")
		call.done()
		return
	}
	_, err = client.Conn.Write(d)
	if err != nil {
		client.mutex.Lock()
		call = client.pending[seq]
		delete(client.pending, seq)
		client.mutex.Unlock()
		if call != nil {
			call.Error = err
			call.done()
		}
		return
	}
}

func (client *Client) input() {
	var err error
	var msgProtocol = &sanrpc.SanRPCProtocol{}

	for err == nil {
		if client.option.ReadTimeout != 0 {
			_ = client.Conn.SetReadDeadline(time.Now().Add(client.option.ReadTimeout))
		}
		msg, err := msgProtocol.DecodeMessage(client.Conn)
		if err != nil {
			log.Debug("DecodeMessage", err)
			break
		}
		res, _ := msg.(*sanrpc.MessageProtocol)
		log.Debugf("resp msg %+v", res)
		if res.Header == nil {
			res.Header = &sanrpc.HeaderMsg{}
		}
		seq := res.Header.Seq
		var call *Call

		client.mutex.Lock()
		call = client.pending[seq]
		delete(client.pending, seq)
		client.mutex.Unlock()

		if res.Err == nil {
			res.Err.Code = 0
		}
		switch {
		case res.Err.Code != 0:
			if len(res.Header.MetaData) > 0 {
				meta := make(map[string]string, len(res.Header.MetaData))
				for k, v := range res.Header.MetaData {
					meta[k] = v
				}
				call.ResMetadata = meta
			}
			call.Error = &errs.Error{
				Type: res.Err.Type,
				Code: res.Err.Code,
				Msg:  res.Err.Msg,
			}
			call.done()
		default:
			data := res.Data
			if len(data) > 0 {
				cc := codec.Codecs[call.SerializeType]
				if cc == nil {
					call.Error = ServiceError(ErrUnsupportedCodec.Error())
				} else {
					err = cc.Decode(data, call.Reply)
					if err != nil {
						log.Error("data decode fail")
						call.Error = ServiceError(err.Error())
					}
				}
			}
			if len(res.Header.MetaData) > 0 {
				meta := make(map[string]string, len(res.Header.MetaData))
				for k, v := range res.Header.MetaData {
					meta[k] = v
				}
				call.ResMetadata = res.Header.MetaData
			}
			call.done()
		}
		res.Reset()
	}

	client.mutex.Lock()
	_ = client.Conn.Close()
	client.shutdown = true
	closing := client.closing
	if err == io.EOF {
		err = ErrShutdown
	} else {
		//err = io.ErrUnexpectedEOF
	}
	for _, call := range client.pending {
		call.Error = err
		call.done()
	}
	client.mutex.Unlock()
	if err != nil && err != io.EOF && !closing {
		log.Error("sanrpc: client protocol error:", err)
	}
}

// Close calls the underlying connection's Close method. If the connection is already
// shutting down, ErrShutdown is returned.
func (client *Client) Close() error {
	client.mutex.Lock()

	for seq, call := range client.pending {
		delete(client.pending, seq)
		if call != nil {
			call.Error = ErrShutdown
			call.done()
		}
	}

	var err error

	if client.closing || client.shutdown {
		client.mutex.Unlock()
		return ErrShutdown
	}
	err = client.Conn.Close()

	client.closing = true
	client.mutex.Unlock()
	return err
}
