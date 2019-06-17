package log

import "github.com/hillguo/sanlog"

func NewSanRpcLoger() *sanlog.Logger {
	l := sanlog.GetLogger("sanrpc")
	writer := sanlog.NewDateWriter("log", "sanrpc", sanlog.DAY, 30)
	l.SetWriter(writer)
	return l
}
