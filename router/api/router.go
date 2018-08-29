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
	e.POST("/create", CreateUser)
	e.GET("/id/read", ReadUser)
	e.GET("/name/read", ReadUsers)
	e.POST("/update", UpdateUser)
	e.POST("/delete", DeleteUser)
}
