package defines

import (
	"github.com/labstack/echo"
)

// Resp return struct
type Resp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Body interface{} `json:"body"`
}

// NewRespMsg create default resp
func NewRespMsg(err *echo.HTTPError, body interface{}) *Resp {
	return &Resp{Code: err.Code, Msg: err.Message.(string), Body: body}
}

func CreateRespMsg(code int, msg string, body interface{}) *Resp {
	return &Resp{Code: code, Msg: msg, Body: body}
}

func NewSuccessMsg(body interface{}) *Resp {
	return &Resp{Code: 0, Msg: "success", Body: body}
}
