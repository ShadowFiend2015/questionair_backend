package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/labstack/echo"

	"questionair_backend/defines"
	"questionair_backend/module"
	log "questionair_backend/util/logger"
	"questionair_backend/util/token"
)

func (h *apiHandler) DownloadLinkConfirmed(e echo.Context) error {
	userId, err := token.GetTokenId(e)
	if err != nil {
		log.Logger().Errorf("DownloadLinkConfirmed: %+v", err)
		return err
	}
	if err := module.CheckUserDownload(userId); err != nil {
		log.Logger().Errorf("DownloadLinkConfirmed: check user's right failed - %v", err)
		return err
	}
	elements, err := module.ReadLinksConfirmed()
	if err != nil {
		log.Logger().Errorf("DownloadLinkConfirmed: %v", err)
		return err
	}
	rsp, err := json.Marshal(elements)
	if err != nil {
		log.Logger().Errorf("DownloadLinkConfirmed: json marshal error - %v", err)
		return err
	}

	file, err := ioutil.TempFile("", fmt.Sprintf("links.json"))
	if err != nil {
		log.Logger().Errorf("DownloadLinkConfirmed: failed to create temp json - %v", err)
		return defines.ComInnerError
	}
	defer os.Remove(file.Name())

	if _, err := file.Write(rsp); err != nil {
		log.Logger().Errorf("DownloadLinkConfirmed: failed to write temp json - %v", err)
		return defines.ComInnerError
	}
	return e.Attachment(file.Name(), "links.json")
}

func (h *apiHandler) DownloadElementsByConfirmedLink(e echo.Context) error {
	userId, err := token.GetTokenId(e)
	if err != nil {
		log.Logger().Errorf("DownloadElementsByConfirmedLink: %+v", err)
		return err
	}
	if err := module.CheckUserDownload(userId); err != nil {
		log.Logger().Errorf("DownloadElementsByConfirmedLink: check user's right failed - %v", err)
		return err
	}

	elements, err := module.ReadElementsByConfirmedLink()
	if err != nil {
		log.Logger().Errorf("DownloadElementsByConfirmedLink: %v", err)
		return err
	}
	rsp, err := json.Marshal(elements)
	if err != nil {
		log.Logger().Errorf("DownloadElementsByConfirmedLink: json marshal error - %v", err)
		return err
	}

	file, err := ioutil.TempFile("", fmt.Sprintf("elements.json"))
	if err != nil {
		log.Logger().Errorf("DownloadElementsByConfirmedLink: failed to create temp json - %v", err)
		return defines.ComInnerError
	}
	defer os.Remove(file.Name())

	if _, err := file.Write(rsp); err != nil {
		log.Logger().Errorf("DownloadElementsByConfirmedLink: failed to write temp json - %v", err)
		return defines.ComInnerError
	}
	return e.Attachment(file.Name(), "elements.json")
}
