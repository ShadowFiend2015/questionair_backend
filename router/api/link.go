package api

import (
	"github.com/labstack/echo"
	"questionair_backend/defines"
	"questionair_backend/module"
	log "questionair_backend/util/logger"
)

type ReqLink struct {
	ScopeName1   string `json:"scope_name_1" form:"scope_name_1" query:"scope_name_1" validate:"required"`
	ScopeName2   string `json:"scope_name_2" form:"scope_name_2" query:"scope_name_2" validate:"required"`
	ElementName1 string `json:"element_name_1" form:"element_name_1" query:"element_name_1" validate:"required"`
	ElementName2 string `json:"element_name_2" form:"element_name_2" query:"element_name_2" validate:"required"`
	HostScope    string `json:"host_scope" form:"host_scope" query:"host_scope" validate:"required"`
}

type ReqConfirmLink struct {
	Id        int64  `json:"id" form:"id" query:"id" validate:"required"`
	HostScope string `json:"host_scope" form:"host_scope" query:"host_scope" validate:"required"`
	Agree     bool   `json:"agree" form:"agree" query:"agree"`
}

func (h *apiHandler) CreateLink(e echo.Context) error {
	req := new(ReqLink)
	if err := e.Bind(req); err != nil {
		log.Logger().Errorf("CreateLink: %+v", err)
		return defines.ComBadParam
	}
	req.HostScope = req.ScopeName1
	if err := e.Validate(req); err != nil {
		log.Logger().Errorf("CreateLink: %+v", err)
		return defines.ComBadParam
	}

	rsp, err := module.CreateLink(req.ScopeName1, req.ScopeName2, req.ElementName1, req.ElementName2, req.HostScope)
	if err != nil {
		log.Logger().Errorf("CreateLink: %v", err)
		return err
	}
	RspData(&e, rsp)
	return nil
}

func (h *apiHandler) ReadLinksByScope(e echo.Context) error {
	req := new(ReqScopeName)
	if err := e.Bind(req); err != nil {
		log.Logger().Errorf("ReadLinksByScope: %+v", err)
		return defines.ComBadParam
	}
	if err := e.Validate(req); err != nil {
		log.Logger().Errorf("ReadLinksByScope: %+v", err)
		return defines.ComBadParam
	}

	rsp, err := module.ReadLinksByScope(req.ScopeName)
	if err != nil {
		log.Logger().Errorf("ReadLinksByScope: %v", err)
		return err
	}
	RspData(&e, rsp)
	return nil
}

func (h *apiHandler) ConfirmLink(e echo.Context) error {
	req := new(ReqConfirmLink)
	if err := e.Bind(req); err != nil {
		log.Logger().Errorf("ConfirmLink: %+v", err)
		return defines.ComBadParam
	}
	if err := e.Validate(req); err != nil {
		log.Logger().Errorf("ConfirmLink: %+v", err)
		return defines.ComBadParam
	}

	confirm := 0
	if req.Agree == true {
		confirm = 1
	}
	if err := module.UpdateLink(req.Id, req.HostScope, confirm); err != nil {
		log.Logger().Errorf("ConfirmLink: %v", err)
		return err
	}
	RspData(&e, nil)
	return nil
}
