package middleware

import (
"fmt"
	"github.com/hillguo/sanrpc/protocol/sanhttp/ctx"
	log "github.com/hillguo/sanlog"
"time"
)

// Logger is logger  middleware
func Logger() ctx.HandlerFunc {
	return func(c *ctx.Context) {
		now := time.Now()
		req := c.Request

		c.Next()
		err := c.Err
		dt := time.Since(now)

		msg := fmt.Sprintf("method:%s, ip:%s, path:%s, params:%s, ret:%d, msg:%s, user cost time:%v",
			req.Method, c.ClientIP(), req.URL.Path, req.Form.Encode(),err.Code, err.Msg, dt.String())

		if err.Code != 0 {
			log.Warn(msg)
		} else {
			log.Info(msg)
		}
	}
}

