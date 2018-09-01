package user

import (
	"net/http"
	"questionair_backend/defines"

	"github.com/labstack/echo"
)

func RspData(e *echo.Context, data interface{}) {
	b := defines.NewSuccessMsg(data)
	(*e).JSON(http.StatusOK, b)
}

// Routers user routers
func Routers(e *echo.Group) {
}
