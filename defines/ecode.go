package defines

import "github.com/labstack/echo"

var (
	ComInnerError   = errorPair(90001, "内部错误")
	ComNotExist     = errorPair(90002, "接口不存在")
	ComUnAuthorized = errorPair(90003, "未鉴权")
	ComAuthFailed   = errorPair(90004, "鉴权失败")
	ComBadParam     = errorPair(90005, "请求参数错误")
	ComNoRight      = errorPair(90006, "没有操作权限")
)

var errorMap = make(map[int]string)

func errorPair(code int, desc string) *echo.HTTPError {
	if v, ok := errorMap[code]; ok {
		panic("error code exist, desc: " + v)
	} else {
		errorMap[code] = desc
		return &echo.HTTPError{Code: code, Message: desc}
	}
}
