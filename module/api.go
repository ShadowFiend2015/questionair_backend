package module

import (
	"crypto/md5"
	"fmt"
	"io"
	"sort"

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
		rsp.Data = append(rsp.Data, scope)
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
		rsp.Data = append(rsp.Data, scope)
	}
	return rsp, nil
}

func CreateLink(link *Link) (RspLinkId, error) {
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
	rsp.Id = link.Id
	return rsp, nil
}

func ReadLinksByScope(scopeName string) (RspData, error) {
	rsp := RspData{
		Data: make([]interface{}, 0),
	}
	scopes, err := readScopes()
	if err != nil {
		log.Logger().Errorf("ReadLinksByScope: read scopes error - %+v", err)
		return rsp, defines.SqlInsertError
	}
	scopeMap := make(map[int64]RspScope)
	for _, scope := range scopes {
		scopeMap[scope.Id] = scope
	}

	scope, err := readScopeByName(scopeName)
	if err != nil {
		log.Logger().Errorf("ReadLinksByScope: %+v", err)
		return rsp, defines.SqlInsertError
	}
	if scope.Id == 0 {
		log.Logger().Errorf("ReadLinksByScope: no scope[%s]", scopeName)
		return rsp, defines.ComBadParam

	}
	links1, err := readLinksByScopeFisrt(scope.Id)
	if err != nil {
		log.Logger().Errorf("ReadLinksByScope: read scope at first error - %v", err)
		return rsp, defines.SqlReadError
	}
	links2, err := readLinksByScopeSecond(scope.Id)
	if err != nil {
		log.Logger().Errorf("ReadLinksByScope: read scope at second error - %v", err)
		return rsp, defines.SqlReadError
	}
	var links []RspLinkSelf
	for _, l := range links1 {
		scopeName, ok := scopeMap[l.ScopeId1]
		if !ok {
			log.Logger().Errorf("ReadLinksByScope: no scope1[%d] in link[%v]", l.ScopeId1, l)
			continue
		}
		links = append(links, RspLinkSelf{
			Id:           l.Id,
			Name:         l.LinkElementName1,
			Code:         l.LinkElementCode1,
			LinkName:     l.LinkElementName2,
			LinkCode:     l.LinkElementCode2,
			LinkFullName: fmt.Sprintf("%s:%s", scopeName, l.LinkElementName2),
			Status:       l.Status & 1,
		})
	}
	for _, l := range links2 {
		scopeName, ok := scopeMap[l.ScopeId2]
		if !ok {
			log.Logger().Errorf("ReadLinksByScope: no scope2[%d] in link[%v]", l.ScopeId2, l)
			continue
		}
		links = append(links, RspLinkSelf{
			Id:           l.Id,
			Name:         l.LinkElementName2,
			Code:         l.LinkElementCode2,
			LinkName:     l.LinkElementName1,
			LinkCode:     l.LinkElementCode1,
			LinkFullName: fmt.Sprintf("%s:%s", scopeName, l.LinkElementName1),
			Status:       (l.Status & 2) >> 1,
		})
	}
	sort.Slice(links, func(i, j int) bool { return links[i].Code < links[j].Code })
	rsp.Total = len(links)
	for _, link := range links {
		rsp.Data = append(rsp.Data, link)
	}
	return rsp, nil

}

func UpdateLink(linkId int64, scopeName string, confirm int) error {
	link, err := readLinkById(linkId)
	if err != nil {
		log.Logger().Errorf("UpdateLink: read link by id[%d] error - %v", linkId, err)
		return rsp, defines.SqlReadError
	} else if link.Id == 0 {
		log.Logger().Errorf("UpdateLink: no link[%d] in mysql", linkId)
		return rsp, defines.SqlNoData
	}
	scope, err := readScopeByName(scopeName)
	if err != nil {
		log.Logger().Errorf("UpdateLink: read scope by name[%s] error - %+v", scopeName, err)
		return rsp, defines.SqlReadError
	} else if scope.Id == 0 {
		log.Logger().Errorf("UpdateLink: no scope[%s] in mysql", scopeName)
		return rsp, defines.SqlNoData
	}
	if link.ScopeId1 == scope.Id {

	}

}
