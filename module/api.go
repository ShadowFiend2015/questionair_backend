package module

import (
	"crypto/md5"
	"fmt"
	"io"
	"sort"

	"questionair_backend/defines"
	log "questionair_backend/util/logger"
)

var (
	scopeIdMap   = make(map[int64]RspScope)
	scopeNameMap = make(map[string]RspScope)
)

func InitScopeMap() error {
	scopes, err := readScopes()
	if err != nil {
		log.Logger().Errorf("InitScopeMap: read scopes error - %+v", err)
		return err
	}
	for _, scope := range scopes {
		scopeIdMap[scope.Id] = scope
		scopeNameMap[scope.Name] = scope
	}
	return nil
}

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
		return rsp, defines.ComLoginFailed
	}
	h := md5.New()
	io.WriteString(h, user.Password)
	if passMD5 != fmt.Sprintf("%x", h.Sum(nil)) {
		log.Logger().Errorf("CheckUser: account [%s] can't match password(md5) [%s]", account, passMD5)
		return rsp, defines.ComLoginFailed
	}
	rsp.Pass = true
	rsp.UserId = user.Id
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

func ReadElementsByConfirmedLink() ([]RspElementLinked, error) {
	var rsp []RspElementLinked
	elementMap := make(map[string]RspElementLinked)
	links, err := readLinksElementConfirmed()
	if err != nil {
		log.Logger().Errorf("ReadElementsByConfirmedLink: read links error - %v", err)
		return rsp, defines.SqlReadError
	}
	for _, link := range links {
		if _, ok := elementMap[link.ElementCode1]; !ok {
			elementMap[link.ElementCode1] = RspElementLinked{
				ScopeId:     link.ScopeId1,
				ScopeCode:   scopeIdMap[link.ScopeId1].Code,
				ElementCode: link.ElementCode1,
				ElementName: link.LinkElementName1,
			}
		}
		if _, ok := elementMap[link.ElementCode2]; !ok {
			elementMap[link.ElementCode2] = RspElementLinked{
				ScopeId:     link.ScopeId2,
				ScopeCode:   scopeIdMap[link.ScopeId2].Code,
				ElementCode: link.ElementCode2,
				ElementName: link.LinkElementName2,
			}
		}
	}
	for _, element := range elementMap {
		rsp = append(rsp, element)
	}
	sort.Slice(rsp, func(i, j int) bool {
		if rsp[i].ScopeId == rsp[j].ScopeId {
			return rsp[i].ElementCode < rsp[j].ElementCode
		}
		return rsp[i].ScopeId < rsp[j].ScopeId
	})
	return rsp, nil

}

func CreateLink(scopeName1, scopeName2, elementName1, elementName2, hostScope string) (RspLinkSelf, error) {
	var rsp RspLinkSelf
	scope1, ok := scopeNameMap[scopeName1]
	if !ok {
		log.Logger().Errorf("CreateLink: no scope[%s] in mysql", scopeName1)
		return rsp, defines.SqlNoData
	}
	scope2, ok := scopeNameMap[scopeName2]
	if !ok {
		log.Logger().Errorf("CreateLink: no scope[%s] in mysql", scopeName2)
		return rsp, defines.SqlNoData
	}
	if scope1.Id > scope2.Id {
		scope1, scope2 = scope2, scope1
		scopeName1, scopeName2 = scopeName2, scopeName1
		elementName1, elementName2 = elementName2, elementName1
	}
	element1, err := readElementByName(elementName1, scope1.Id)
	if err != nil {
		log.Logger().Errorf("CreateLink: read element1 by name[%s] scope[%d] error - ", elementName1, scope1.Id, err)
		return rsp, defines.SqlReadError
	}
	element2, err := readElementByName(elementName2, scope2.Id)
	if err != nil {
		log.Logger().Errorf("CreateLink: read element2 by name[%s] scope[%d] error - ", elementName2, scope2.Id, err)
		return rsp, defines.SqlReadError
	}
	if element1.Id == 0 || element2.Id == 0 {
		log.Logger().Errorf("CreateLink:  element1 scope[%d] name[%s] id[%d] element2 scope[%d] name[%s] id[%d]", scope1.Id, elementName1, element1.Id, scope2.Id, elementName2, element2.Id)
		return rsp, defines.ComBadParam
	}
	status := 0
	if hostScope == scopeName1 {
		status = 1
		rsp.Name = element1.Name
		rsp.Code = element1.Code
		rsp.LinkName = element2.Name
		rsp.LinkCode = element2.Code
		rsp.LinkFullName = fmt.Sprintf("%s:%s", scope2.Name, rsp.LinkName)
		rsp.Status = 1
	} else if hostScope == scopeName2 {
		status = 2
		rsp.Name = element2.Name
		rsp.Code = element2.Code
		rsp.LinkName = element1.Name
		rsp.LinkCode = element1.Code
		rsp.LinkFullName = fmt.Sprintf("%s:%s", scope1.Name, rsp.LinkName)
		rsp.Status = 1
	} else {
		log.Logger().Errorf("CreateLink: hostScope[%s] cant't match scope1[%s] or scope2[%s]", hostScope, scopeName1, scopeName2)
		return rsp, defines.ComBadParam
	}
	link := &Link{
		ScopeId1:     scope1.Id,
		ScopeId2:     scope2.Id,
		ElementCode1: element1.Code,
		ElementCode2: element2.Code,
		Status:       status,
	}

	if has, err := countLinkByCode(element1.Code, element2.Code); err != nil {
		log.Logger().Errorf("CreateLink: count repeated link error - ", err)
		return rsp, defines.SqlReadError
	} else if has != 0 {
		log.Logger().Errorf("CreateLink: repeated link element_code[%s][%s]", element1.Code, element2.Code)
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

	scope, ok := scopeNameMap[scopeName]
	if !ok {
		log.Logger().Errorf("ReadLinksByScope: no scope[%s] in mysql", scopeName)
		return rsp, defines.SqlNoData
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
		scope2, ok := scopeIdMap[l.ScopeId2]
		if !ok {
			log.Logger().Errorf("ReadLinksByScope: no scope2[%d] of link[%v] in mysql", l.ScopeId2, l)
			continue
		}
		links = append(links, RspLinkSelf{
			Id:           l.Id,
			Name:         l.LinkElementName1,
			Code:         l.ElementCode1,
			LinkName:     l.LinkElementName2,
			LinkCode:     l.ElementCode2,
			LinkFullName: fmt.Sprintf("%s:%s", scope2.Name, l.LinkElementName2),
			Status:       l.Status & 1,
		})
	}
	for _, l := range links2 {
		scope1, ok := scopeIdMap[l.ScopeId1]
		if !ok {
			log.Logger().Errorf("ReadLinksByScope: no scope1[%d] of link[%v] in mysql", l.ScopeId1, l)
			continue
		}
		links = append(links, RspLinkSelf{
			Id:           l.Id,
			Name:         l.LinkElementName2,
			Code:         l.ElementCode2,
			LinkName:     l.LinkElementName1,
			LinkCode:     l.ElementCode1,
			LinkFullName: fmt.Sprintf("%s:%s", scope1.Name, l.LinkElementName1),
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

func ReadLinksConfirmed() ([]RspLinkDownload, error) {
	var rsp []RspLinkDownload
	links, err := readLinksConfirmed()
	if err != nil {
		log.Logger().Errorf("ReadLinksConfirmed: read links error - %v", err)
		return rsp, defines.SqlReadError
	}
	for _, link := range links {
		rsp = append(rsp, RspLinkDownload{
			Element1: link.ElementCode1,
			Element2: link.ElementCode2,
		})
	}
	return rsp, nil
}

func UpdateLink(linkId int64, hostScope string, confirm int) error {
	link, err := readLinkById(linkId)
	if err != nil {
		log.Logger().Errorf("UpdateLink: read link by id[%d] error - %v", linkId, err)
		return defines.SqlReadError
	} else if link.Id == 0 {
		log.Logger().Errorf("UpdateLink: no link[%d] in mysql", linkId)
		return defines.SqlNoData
	}
	scope, ok := scopeNameMap[hostScope]
	if !ok {
		log.Logger().Errorf("UpdateLink: no scope[%s] in mysql", hostScope)
		return defines.SqlNoData
	}
	if link.ScopeId1 == scope.Id {
		link.Status = (link.Status & 2) + confirm
	} else if link.ScopeId2 == scope.Id {
		link.Status = (link.Status & 1) + confirm<<1
	}
	if err := updateLinkStatus(&link); err != nil {
		log.Logger().Errorf("UpdateLink: update link[%d] status error - %v", linkId, err)
		return defines.SqlUpdateError
	}
	return nil
}
