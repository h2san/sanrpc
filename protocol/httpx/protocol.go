package httpx

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/hillguo/sanrpc/protocol"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"fmt"
)

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
	p.router.GET("/:service/:method", p.routeHander)
	p.router.POST("/:service/:method", p.routeHander)
	fmt.Print(r)
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
