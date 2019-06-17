package httpx

import (
	"github.com/gorilla/websocket"
	"github.com/h2san/sanrpc/protocol"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

//HTTProtocol 路由实现
type HTTProtocol struct {
	protocol.BaseService
	plugins httpxPlugin
	router httprouter.Router
	ws *websocket.Upgrader
}

func (p *HTTProtocol) AddPlugin(plugin interface{}){
	p.plugins.Add(plugin)
}

func (p *HTTProtocol) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.router.ServeHTTP(w,r)
}

func (p *HTTProtocol) routeHander(w http.ResponseWriter, r *http.Request,param httprouter.Params) {
	err := p.plugins.DoPreHandleHTTPRequest(w,r,param)
	if err!=nil{
		writeErrResponse(w,HTTP_PRE_HANDLER_ERR,err.Error())
		return
	}

	p.handle(w,r,param)

	err = p.plugins.DoPreHandleHTTPRequest(w,r,param)
	if err!=nil{
		writeErrResponse(w,HTTP_POST_HANDLER_ERR,err.Error())
		return
	}
}

func (p *HTTProtocol) RegisterService(service interface{}) error{
	p.router.GET("/:service/:method", p.routeHander)
	p.router.POST("/:service/:method",p.routeHander)
	p.router.GET("/:protocol/:service/:method", p.routeHander)
	p.router.POST("/:protocol/:service/:method", p.routeHander)
	return p.BaseService.RegisterService(service)
}