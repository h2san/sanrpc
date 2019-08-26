package protocol

import (
	"context"
	"errors"
	"fmt"
	log "github.com/hillguo/sanlog"
	"reflect"
	"strings"
	"sync"
)

type BaseService struct {
	ServiceMapMu sync.RWMutex
	ServiceMap   map[string]*service
}

type methodType struct {
	sync.Mutex // protects counters
	method     reflect.Method
	ArgType    reflect.Type
	ReplyType  reflect.Type
}

type service struct {
	name   string                 // name of service
	rcvr   reflect.Value          // receiver of methods for the service
	typ    reflect.Type           // type of the receiver
	method map[string]*methodType // registered methods
}

func (s *service) GetMethod(methodName string) *methodType {
	return s.method[methodName]
}

func (p *BaseService) RegisterService(rcvr interface{}) error {
	_, err := p.register(rcvr)
	return err
}

func (p *BaseService) register(rcvr interface{}) (string, error) {
	p.ServiceMapMu.Lock()
	defer p.ServiceMapMu.Unlock()
	if p.ServiceMap == nil {
		p.ServiceMap = make(map[string]*service)
	}

	service := new(service)
	service.typ = reflect.TypeOf(rcvr)
	service.rcvr = reflect.ValueOf(rcvr)
	sname := reflect.Indirect(service.rcvr).Type().Name() // Type

	if service.typ.Kind() != reflect.Ptr {
		msg := "register service shoule be point struct"
		log.Error(msg)
		return sname, errors.New(msg)
	}

	service.name = strings.ToLower(sname)

	// Install the methods
	service.method = suitableMethods(service.typ)

	if len(service.method) == 0 {
		errorStr := "sanrpc.Register: type " + sname + " has no exported methods of suitable type"
		log.Error(errorStr)
		return sname, errors.New(errorStr)
	}
	p.ServiceMap[service.name] = service

	log.Infof("register service: %v success ", service.name)
	for key, _ := range service.method {
		log.Infof("register method: %v success ", key)
	}

	return sname, nil
}

// suitableMethods returns suitable Rpc methods of typ, it will report
// error using log if reportErr is true.
func suitableMethods(typ reflect.Type) map[string]*methodType {
	methods := make(map[string]*methodType)
	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		mtype := method.Type
		mname := strings.ToLower(method.Name)
		// Method must be exported.
		if method.PkgPath != "" {
			continue
		}
		// Method needs four ins: receiver, context.Context, *args, *reply.
		if mtype.NumIn() != 4 {
			log.Info("method", mname, " has wrong number of ins:", mtype.NumIn())
			continue
		}
		// First arg must be context.Context
		ctxType := mtype.In(1)
		if !ctxType.Implements(reflect.TypeOf((*context.Context)(nil)).Elem()) {
			log.Info("method", mname, " must use context.Context as the first parameter")
			continue
		}

		// Second arg need not be a pointer.
		argType := mtype.In(2)

		// Third arg must be a pointer.
		replyType := mtype.In(3)
		if replyType.Kind() != reflect.Ptr {
			log.Info("method", mname, " reply type not a pointer:", replyType)
			continue
		}

		// Method needs one out.
		if mtype.NumOut() != 1 {
			log.Info("method", mname, " has wrong number of outs:", mtype.NumOut())
			continue
		}
		// The return type of the method must be error.
		if returnType := mtype.Out(0); returnType != reflect.TypeOf((*error)(nil)).Elem() {
			log.Info("method", mname, " returns ", returnType.String(), " not error")
			continue
		}
		methods[mname] = &methodType{method: method, ArgType: argType, ReplyType: replyType}
	}
	return methods
}

func (s *service) Call(ctx context.Context, mtype *methodType, argv, replyv reflect.Value) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[service internal error]: %v, method: %s, argv: %+v",
				r, mtype.method.Name, argv.Interface())
			log.Error(err)
		}
	}()

	function := mtype.method.Func
	// Invoke the method, providing a new value for the reply.
	returnValues := function.Call([]reflect.Value{s.rcvr, reflect.ValueOf(ctx), argv, replyv})
	// The return value for the method is an error.
	errInter := returnValues[0].Interface()
	if errInter != nil {
		return errInter.(error)
	}

	return nil
}
