package httpx

const (
	HTTPX_REQ_UNMARSHAL_ERR = 1000001
	HTTPX_RESP_MARSHAL_ERR  = 1000002
	HTTP_REQ_HANDLE_ERR     = 1000003

	// pre handler http errro
	HTTP_PRE_HANDLER_ERR  = 1000004
	HTTP_POST_HANDLER_ERR = 1000004
)

type error302 struct {
}

type error404 struct {
}

func (e error302) Error() string {
	return "302"
}

func (e error404) Error() string {
	return "404"
}
