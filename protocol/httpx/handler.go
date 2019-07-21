package httpx

import (
	"context"
	"encoding/json"
	log "github.com/hillguo/sanlog"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

const (
	WebSocketKey   = "Context_WebSocketKey"
	HTTPRequestKey = "Context_HTTPRequestKey"
)

func (p *HTTProtocol) handle(ctx context.Context, w http.ResponseWriter, req *http.Request, param httprouter.Params) {
	serviceName := param.ByName("service")
	methodName := param.ByName("method")

	if serviceName == "" || methodName == "" {
		http.NotFound(w, req)
		log.Warn(req.URL, "serviceName or methodName null")
		return
	}

	service := p.ServiceMap[serviceName]

	if service == nil {
		log.Warn("req url: ", req.URL, " not found service")
		http.NotFound(w, req)
		return
	}
	mtype := service.GetMethod(methodName)
	if mtype == nil {
		http.NotFound(w, req)
		log.Warn("req url: ", req.URL, "not found methodName")
		return
	}

	argv := reflect.New(mtype.ArgType.Elem()).Interface()
	replyv := reflect.New(mtype.ReplyType.Elem()).Interface()

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
			log.Error("websocket Upgrade error: ", err.Error())
			return
		}
		ctx = context.WithValue(ctx, WebSocketKey, conn)
		break
	default:
		data, err := ioutil.ReadAll(req.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Error("ioutil.ReadAll error:", err)
			return
		}
		err = json.Unmarshal(data, argv)
		if err != nil {
			writeErrResponse(w, HTTPX_REQ_UNMARSHAL_ERR, err.Error())
			log.Error("json Unmarshal error:", string(data), err)
			return
		}
		ctx = context.WithValue(ctx, HTTPRequestKey, req)
	}

	var err error
	if mtype.ArgType.Kind() != reflect.Ptr {
		err = service.Call(ctx, mtype, reflect.ValueOf(argv).Elem(), reflect.ValueOf(replyv))
	} else {
		err = service.Call(ctx, mtype, reflect.ValueOf(argv), reflect.ValueOf(replyv))
	}

	if err != nil {
		writeErrResponse(w, HTTP_REQ_HANDLE_ERR, err.Error())
		log.Error("service.Call handle error:", argv, err)
		return
	}

	writeResponse(w, 0, "success", replyv)
	log.Info("return 0 ", "req:", argv, "resp:", replyv)

}
