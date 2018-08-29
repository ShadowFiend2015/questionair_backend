package user

import (
	"github.com/labstack/echo"

	"questionair_backend/defines"
	"questionair_backend/module"
	log "questionair_backend/util/logger"
	"questionair_backend/util/token"
)

type ReqAuth struct {
	Account  string `json:"account" form:"account" query:"account" validate:"required"`
	Password string `json:"password" form:"password" query:"password" validate:"required"`
}

func UserLogin(e echo.Context) error {
	req := new(ReqAuth)
	if err := e.Bind(req); err != nil {
		log.Logger().Errorf("UserLogin: %+v", err)
		return defines.ComBadParam
	}
	if err := e.Validate(req); err != nil {
		log.Logger().Errorf("UserLogin: %+v", err)
		return defines.ComBadParam
	}
	rsp, err := module.CheckUser(req.Account, req.Password)
	if err != nil {
		log.Logger().Errorf("UserLogin: %+v", err)
		return err
	}
	if rsp.UserId != 0 {
		info := make(map[string]interface{})
		info["id"] = rsp.UserId
		info["role"] = "user"
		t, err := token.CreateToken(info)
		if err != nil {
			log.Logger().Errorf("UserLogin: create token failed - %v", err)
			return defines.ComInnerError
		}
		rsp.Token = t
	}
	RspData(&e, rsp)
	return nil
}
