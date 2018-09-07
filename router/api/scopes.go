package api

import (
	"github.com/labstack/echo"
	"questionair_backend/defines"
	"questionair_backend/module"
	log "questionair_backend/util/logger"
)

type ReqScopeExcept struct {
	ScopeName string `json:"scope_name" form:"scope_name" query:"scope_name" validate:"required"`
}

func (h *apiHandler) ReadScopes(e echo.Context) error {
	rsp, err := module.ReadScopes()
	if err != nil {
		log.Logger().Errorf("ReadScopes: %v", err)
		return err
	}
	RspData(&e, rsp)
	return nil
}

func (h *apiHandler) ReadScopesExceptOne(e echo.Context) error {
	req := new(ReqScopeExcept)
	if err := e.Bind(req); err != nil {
		log.Logger().Errorf("ReadScopesExceptOne: %+v", err)
		return defines.ComBadParam
	}

	if err := e.Validate(req); err != nil {
		log.Logger().Errorf("ReadScopesExceptOne: %+v", err)
		return defines.ComBadParam
	}
	rsp, err := module.ReadScopesExceptOne(req.ScopeName)
	if err != nil {
		log.Logger().Errorf("ReadScopesExceptOne: %v", err)
		return err
	}
	RspData(&e, rsp)
	return nil
}
