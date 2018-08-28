package defines

import (
	"github.com/labstack/echo"
)

type Resp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Body interface{} `json:"body"`
}
