package httpx

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/gorilla/websocket"
	"github.com/hillguo/sanrpc/log"
	"github.com/julienschmidt/httprouter"
)

const (
	WebSocketKey   = "Context_WebSocketKey"
	HTTPRequestKey = "Context_HTTPRequestKey"
)

func (p *HTTProtocol) handle(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
	serviceName := param.ByName("service")
	methodName := param.ByName("method")
	log.Debug("service:", serviceName, " , ", method:", methodName)
	if serviceName == "" || methodName == "" {
		http.NotFound(w, req)
		return
	}

	service := p.ServiceMap[serviceName]

	if service == nil {
		http.NotFound(w, req)
		return
	}
	mtype := service.GetMethod(methodName)
	if mtype == nil {

		http.NotFound(w, req)
		return
	}

	argv := reflect.New(mtype.ArgType.Elem()).Interface()
	replyv := reflect.New(mtype.ReplyType.Elem()).Interface()

	ctx := context.Background()

	switch param.ByName("protocol") {
	case "ws":
		if p.ws == nil {
			p.ws = &websocket.Upgrader{
				ReadBufferSize:  4096,
				WriteBufferSize: 1024,
			}
		}
		conn, err := p.ws.Upgrade(w, req, nil)
		if err != nil {
			writeErrResponse(w, HTTP_REQ_HANDLE_ERR, err.Error())
			return
		}
		ctx = context.WithValue(ctx, WebSocketKey, conn)
		break
	default:
		data, err := ioutil.ReadAll(req.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(data, argv)
		if err != nil {
			writeErrResponse(w, HTTPX_REQ_UNMARSHAL_ERR, err.Error())
			return
		}
		ctx = context.WithValue(ctx,HTTPRequestKey,req)
		w.Header().Set("Content-Type","application/json")
	}

	var err error
	if mtype.ArgType.Kind() != reflect.Ptr {
		err = service.Call(ctx, mtype, reflect.ValueOf(argv).Elem(), reflect.ValueOf(replyv))
	} else {
		err = service.Call(ctx, mtype, reflect.ValueOf(argv), reflect.ValueOf(replyv))
	}

	if err != nil {
		if _, ok := err.(error302); ok {
			return
		}
		if _, ok := err.(error404); ok {
			return
		}
		writeErrResponse(w, HTTP_REQ_HANDLE_ERR, err.Error())
		return
	}

	if w.Header().Get("Content-Type") == "application/json"{
		writeResponse(w,0,"success", replyv)
	}

}
