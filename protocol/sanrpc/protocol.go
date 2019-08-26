package sanrpc

import (
	"context"
	"encoding/binary"
	"errors"
	"github.com/golang/protobuf/proto"
	log "github.com/hillguo/sanlog"
	"github.com/hillguo/sanrpc/codec"
	"github.com/hillguo/sanrpc/errs"
	"github.com/hillguo/sanrpc/protocol"
	"io"
	"reflect"
	"strings"
)

const headLen int = 8 // 4 bytes magic_value + 4 bytes pb length

var (
	ErrReadMsgHeaderInvalid = errors.New("read msg header num invalid")
	ErrReadMsgMagicInvalid    = errors.New("read msg header magic invalid")
	ErrReadMsgBodyInvalid   = errors.New("read msg body num invalid")
	ErrServerMarshalFail    = errors.New("server marshal response interface invalid")
	ErrServerUnmarshalFail    = errors.New("server unmarshal request interface invalid")
	ErrClientMarshalFail    = errors.New("client marshal request interface invalid")
	ErrClientUnmarshalFail  = errors.New("client unmarshal response interface invalid")
	ErrMsgAssertInvalid = errors.New("msg assert fail")
)

type SanRPCProtocol struct {
	protocol.BaseService
}

func (p *SanRPCProtocol) DecodeMessage(r io.Reader) (protocol.Message, error) {

	head := make([]byte, headLen)
	n ,err := io.ReadFull(r,head)
	if err != nil {
		log.Debug(err)
		return nil, err
	}
	if n != headLen {
		return nil, ErrReadMsgHeaderInvalid
	}
	magic := binary.BigEndian.Uint32(head[:4])
	bodylen := binary.BigEndian.Uint32(head[4:])

	log.Infof("read msg magic:%d, body len:%d", magic, bodylen)
	if magic != uint32(SanrpcMagic_SANRPC_MAGIC_VALUE) {
		return nil, ErrReadMsgMagicInvalid
	}
	if bodylen == 0 {
		return &MessageProtocol{}, nil
	}
	msg := make([]byte, bodylen)
	n,err = io.ReadFull(r, msg)
	if  err != nil {
		return nil, err
	}
	if n != int(bodylen) {
		return nil, ErrReadMsgBodyInvalid
	}
	log.Info("read a full sanrpc message")

	msgprotocol := &MessageProtocol{}
	if err :=proto.Unmarshal(msg, msgprotocol); err != nil {
		return nil, ErrServerUnmarshalFail
	}
	log.Info("unmarshal sanrpc message success")
	return msgprotocol, nil
}

func (p *SanRPCProtocol) EncodeMessage(msg protocol.Message) ([]byte,error) {

	m, ok := msg.(*MessageProtocol)
	log.Debug(msg)
	if !ok {
		log.Errorf("rpc: encoding msg error %+v", msg)
		return nil , ErrMsgAssertInvalid
	}
	data, err := proto.Marshal(m)
	if err != nil {
		return nil, ErrServerMarshalFail
	}

	d := make([]byte, headLen + len(data))

	binary.BigEndian.PutUint32(d[:4], uint32(SanrpcMagic_SANRPC_MAGIC_VALUE))
	binary.BigEndian.PutUint32(d[4:8], uint32(len(data)))
	copy(d[8:],data)
	return d, nil
}

func (p *SanRPCProtocol) DisspatchMessage(req *MessageProtocol, resp *MessageProtocol)  error{
	serviceName := strings.ToLower(req.Header.ServiceName)
	methodName := strings.ToLower(req.Header.MethodName)

	p.ServiceMapMu.RLock()
	service := p.ServiceMap[serviceName]
	p.ServiceMapMu.RUnlock()

	if service == nil {
		return errs.ErrServerNoService
	}
	mtype := service.GetMethod(methodName)
	if mtype == nil {
		return errs.ErrServerNoMethod
	}
	var argv interface{}
	if mtype.ArgType.Kind() != reflect.Ptr {
		argv = reflect.New(mtype.ArgType).Elem().Interface()
	} else {
		argv = reflect.New(mtype.ArgType.Elem()).Interface()
	}


	cc := codec.Codecs[codec.SerializeType(req.Header.EncodeType)]
	if cc == nil {
		return errs.ErrServerNoSupportEncodeType
	}
	err := cc.Decode(req.Data, argv)
	if err != nil {
		return errs.ErrServerDecodeDataErr
	}

	replyv := reflect.New(mtype.ReplyType.Elem()).Interface()

	ctx:=context.Background()


	err = service.Call(ctx, mtype, reflect.ValueOf(argv), reflect.ValueOf(replyv))


	if err != nil {
		return err
	}

	data, err := cc.Encode(replyv)
	log.Debug(data)
	if err != nil {
		return errs.ErrServerEncodeDataErr
	}
	resp.Data = data
	return nil
}

func (p *SanRPCProtocol) HandleMessage(ctx context.Context, r protocol.Message) (protocol.Message, error) {
	req, ok := r.(*MessageProtocol)
	if !ok {
		return nil, ErrMsgAssertInvalid
	}

	resp := &MessageProtocol{
		Err:&ErrMsg{

		},
	}

	log.Infof("req msg: %v", req)
	err := p.DisspatchMessage(req, resp)
	log.Error(err)
	if err != nil {
		if e, ok := err.(*errs.Error); ok {
			resp.Err.Type = e.Type
			resp.Err.Code = e.Code
			resp.Err.Msg = e.Msg
			return resp, nil
		}
		resp.Err.Type = 0
		resp.Err.Code = 999
		resp.Err.Msg = err.Error()
		return resp, nil
	}
	log.Infof("resp msg: %v", resp)
	resp.Err = &ErrMsg{
		Code:0,
		Type:1,
		Msg:"success",
	}
	return resp, nil
}