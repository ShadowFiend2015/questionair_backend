package router

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	validator "gopkg.in/go-playground/validator.v9"

	"questionair_backend/conf"
	"questionair_backend/defines"
	"questionair_backend/module"
	"questionair_backend/router/api"
	"questionair_backend/router/user"
	log "questionair_backend/util/logger"
)

func RspData(e *echo.Context, data interface{}) {
	b := defines.NewSuccessMsg(data)
	(*e).JSON(http.StatusOK, b)
}

func errorHandler(err error, c echo.Context) {
	he := err.(*echo.HTTPError)
	code := he.Code

	log.Logger().Print(he.Error())
	msg := he.Message.(string)
	if code == 500 {
		rsp := defines.NewRespMsg(defines.ComInnerError, nil)
		c.JSON(200, rsp)
	} else if code == 404 {
		rsp := defines.NewRespMsg(defines.ComNotExist, nil)
		c.JSON(200, rsp)
	} else if code == 401 {
		rsp := defines.NewRespMsg(defines.ComAuthFailed, nil)
		c.JSON(200, rsp)
	} else if code == 400 {
		rsp := defines.NewRespMsg(defines.ComUnAuthorized, nil)
		c.JSON(200, rsp)
	} else if code == 405 {
		rsp := defines.NewRespMsg(defines.ComNotExist, nil)
		c.JSON(200, rsp)
	} else if code >= 90000 {
		rsp := defines.CreateRespMsg(code, msg, nil)
		c.JSON(200, rsp)
	}

	log.Logger().Infof("error code = %d", code)
}

type Vallidator struct {
	validator *validator.Validate
}

func (v *Vallidator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

func RunServer(runmode string) {
	if err := conf.InitConfig(runmode); err != nil {
		log.Logger().Error("read config file err")
		panic(err)
	}

	if err := module.InitSql(); err != nil {
		log.Logger().Error("init mysql err")
		panic(err)
	}
	if err := module.InitScopeMap(); err != nil {
		log.Logger().Error("init scope_map err")
		panic(err)
	}

	srv := echo.New()
	srv.Use(middleware.BodyLimit("20M"))

	srv.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           "${time_custom}  ${remote_ip}  ${method}  ${path}  ${status}  ${latency_human}\n",
		CustomTimeFormat: "2006/01/02 15:04:05",
	}))
	srv.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST},
	}))

	srv.HTTPErrorHandler = errorHandler
	srv.Validator = &Vallidator{validator: validator.New()}

	srv.POST("/user/login", user.UserLogin)

	at := srv.Group("", middleware.JWTWithConfig(middleware.JWTConfig{
		ContextKey:  "usr",
		SigningKey:  []byte(conf.Conf.Token.Salt),
		TokenLookup: "header:token",
	}))

	apiGroup := at.Group("/api")
	api.Routers(apiGroup)

	go func() {
		if err := srv.Start(conf.Conf.Server.Addr); err != nil {
			srv.Logger.Panicf("Shutting down the server with error:%v", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		srv.Logger.Fatal(err)
	}
}
