package httpx

import (
	"context"
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
	router  httprouter.Router
}


func (p *HTTProtocol) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.router.ServeHTTP(w, r)
}

func (p *HTTProtocol) routeHander(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	ctx := context.Background()
	p.handle(ctx, w, r, param)
}
