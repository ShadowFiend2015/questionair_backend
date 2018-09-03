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

func ReadScopes() (RspData, error) {
	rsp := RspData{
		Data: make([]interface{}, 0),
	}
	scopes, err := readScopes()
	if err != nil {
		log.Logger().Errorf("ReadScopes: %+v", err)
		return rsp, defines.SqlReadError
	}
	rsp.Total = len(scopes)
	for _, scope := range scopes {
		rsp.Data = append(rsp.Data, link)
	}
	return rsp, nil
}

func ReadScopesExceptOne(name string) (RspData, error) {
	rsp := RspData{
		Data: make([]interface{}, 0),
	}
	scopes, err := readScopesExceptOne(name)
	if err != nil {
		log.Logger().Errorf("ReadScopesExceptOne: %+v", err)
		return rsp, defines.SqlReadError
	}
	rsp.Total = len(scopes)
	for _, scope := range scopes {
		rsp.Data = append(rsp.Data, link)
	}
	return rsp, nil
}

func CreateLink(link *Link) (int64, error) {
	var rsp RspLinkId
	if link.ScopeId1 >= link.ScopeId2 {
		log.Logger().Errorf("CreateLink: link scope_id_1[%d] >= scope_id_2[%d]", link.ScopeId1, link.ScopeId2)
		return rsp, defines.ComInnerError
	}
	if link.ElementCode1 >= link.ElementCode2 {
		log.Logger().Errorf("CreateLink: link element_code_1[%d] >= element_code_2[%d]", link.ElementCode1, link.ElementCode2)
		return rsp, defines.ComInnerError
	}
	if has, err := countLinkByCode(link.ElementCode1, link.ElementCode2); err != nil {
		log.Logger().Errorf("CreateLink: count repeat link error - ", err)
		return rsp, defines.SqlReadError
	} else if has != 0 {
		log.Logger().Errorf("CreateLink: count repeat link error - ", err)
		return rsp, defines.ComDuplicate
	}
	if err := createLink(link); err != nil {
		log.Logger().Errorf("CreateLink: %+v", err)
		return rsp, defines.SqlInsertError
	}
	rsp.Id = u.Id
	return rsp, nil
}

func ReadLinksByScope(scopeName string) (RspData, error) {

}
