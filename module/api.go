package module

import (
	"crypto/md5"
	"fmt"
	"io"

	"questionair_backend/defines"
	log "questionair_backend/util/logger"
)

func CheckUser(account, passMD5 string) (RspUserCheck, error) {
	rsp := RspUserCheck{
		Pass: false,
	}
	user, err := readUserByAccount(account)
	if err != nil {
		log.Logger().Errorf("CheckUser: %+v", err)
		return rsp, defines.SqlReadError
	}
	if user.Id == 0 {
		log.Logger().Errorf("CheckUser: account [%s] not found in the database", account)
		return rsp, nil
	}
	h := md5.New()
	io.WriteString(h, user.Password)
	if passMD5 != fmt.Sprintf("%x", h.Sum(nil)) {
		log.Logger().Errorf("CheckUser: account [%s] can't match password(md5) [%s]", account, passMD5)
		return rsp, nil
	}
	rsp.Pass = true
	return rsp, nil
}
