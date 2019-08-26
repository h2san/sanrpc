package sanhttp

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

type Msg struct {
	Req  *http.Request
	Resp http.ResponseWriter
}

type ReqMsg struct {
	Service string
	Method  string
	Type    string
	Data    string
}

type RespMsg struct {
	Code int32
	Msg  string
	Data string
}

func (msg *ReqMsg) Decode(r io.Reader) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, msg)
	if err != nil {
		return err
	}
	return nil
}

func (msg *RespMsg) Write(w http.ResponseWriter) error {
	data ,err := json.Marshal(msg)
	if err != nil {
		return  err
	}
	_, err = w.Write(data)
	if err != nil {
		return err
	}
	return nil
}
