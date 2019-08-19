package sanhttp

import "net/http"
import  "github.com/golang/protobuf/proto"
type BaseCGI struct {
	w http.ResponseWriter
	r *http.Request

	req proto.Message
	resp proto.Message
}



