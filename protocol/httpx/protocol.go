package httpx

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/hillguo/sanrpc/protocol"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func init() {
	DefaultHTTProtocol.router.POST("/:service/:method", DefaultHTTProtocol.routeHander)
	DefaultHTTProtocol.router.GET("/:service/:method", DefaultHTTProtocol.routeHander)

}

var DefaultHTTProtocol = &HTTProtocol {

}

//HTTProtocol 路由实现
type HTTProtocol struct {
	protocol.BaseService
	plugins httpxPlugin
	router  httprouter.Router
	ws      *websocket.Upgrader
}

func (p *HTTProtocol) AddPlugin(plugin interface{}) {
	p.plugins.Add(plugin)
}

func (p *HTTProtocol) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.router.ServeHTTP(w, r)
}

func (p *HTTProtocol) routeHander(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	ctx := context.Background()
	err := p.plugins.DoPreHandleHTTPRequest(ctx, r)
	if err != nil {
		writeErrResponse(w, HTTP_PRE_HANDLER_ERR, err.Error())
		return
	}

	p.handle(ctx, w, r, param)
	p.plugins.DoPostHandleHTTPRequest(ctx)
}
