package httpx

import (
	"encoding/json"
	"net/http"
)

type ErrResp struct {
	Code int
	ErrMsg string
}

type CommResp struct {
	Code int
	ErrMsg string
	Data interface{}
}



func writeErrResponse(w http.ResponseWriter, code int, errMsg string) {

	resp := ErrResp{
		Code:code,
		ErrMsg:errMsg,
		}

	d, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Context-Type", "application/json")
	w.Write(d)
}

func writeResponse(w http.ResponseWriter, code int, errMsg string, data interface{} ) {

	resp := CommResp{
		Code:code,
		ErrMsg:errMsg,
		Data:data,
	}

	d, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Context-Type", "application/json")
	w.Write(d)
}

