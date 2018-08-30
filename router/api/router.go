package api

import (
	"net/http"
	"questionair_backend/defines"

	"github.com/labstack/echo"
)

func RspData(e *echo.Context, data interface{}) {
	b := defines.NewSuccessMsg(data)
	(*e).JSON(http.StatusOK, b)
}

type apiHandler struct{}

// Routers api routers
func Routers(e *echo.Group) {

	h := new(apiHandler)

	user := e.Group("", md.UserVertify(token.VerifyUser))

	user.POST("/check", h.Check)
}
