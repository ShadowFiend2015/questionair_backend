package api

import (
	"net/http"

	"github.com/labstack/echo"

	"questionair_backend/defines"
	md "questionair_backend/middleware"
	"questionair_backend/util/token"
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

	user.GET("/scope/read", h.ReadScopes)
	user.GET("/scope/other/read", h.ReadScopesExceptOne)

	user.POST("/link/create", h.CreateLink)
	user.GET("/link/scope/read", h.ReadLinksByScope)
	user.POST("/link/confirm", h.ConfirmLink)
}
