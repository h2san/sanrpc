package sanhttp

import sp "github.com/hillguo/sanhttp"

type SanHttp struct {
	e *sp.Engine
}

func (s *SanHttp) RegisterService(rcvr interface{}) error {
	// TODO
	// url = serverName/methodName
	return nil
}

var DefaultHTTProtocol = &SanHttp{
	e: sp.Default(),
}
