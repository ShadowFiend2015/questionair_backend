package module

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"

	"questionair_backend/conf"
)

var engine *xorm.Engine

func InitSql() error {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", conf.Conf.Sql.User, conf.Conf.Sql.Password, conf.Conf.Sql.Addr, conf.Conf.Sql.DB)
	if e, err := xorm.NewEngine("mysql", dataSourceName); err != nil {
		return err
	} else {
		engine = e
	}
	if err := SyncDatabase(new(User), new(Scope), new(Element), new(Link)); err != nil {
		return err
	}
	return nil
}

func SyncDatabase(data ...interface{}) error {
	return engine.Sync2(data...)
}

func readUserByAccount(account string) (User, error) {
	var user User
	_, err := engine.Table("user").Where("user.account = ?", account).Get(&user)
	return user, err
}

func readUserById(id int64) (User, error) {
	var user User
	_, err := engine.Table("user").Id(id).Get(&user)
	return user, err
}

func readScopeByName(name string) (RspScope, error) {
	var scope RspScope
	_, err := engine.Table("scope").Where("name = ?", name).Get(&scope)
	return scope, err
}

func readScopes() ([]RspScope, error) {
	var scopes []RspScope
	err := engine.Table("scope").Asc("scope.id").Find(&scopes)
	return scopes, err
}

func readScopesExceptOne(name string) ([]RspScope, error) {
	var scopes []RspScope
	err := engine.Table("scope").Where("scope.name <> ?", name).Asc("scope.id").Find(&scopes)
	return scopes, err
}

func readElementByName(name string, scopeId int64) (Element, error) {
	var element Element
	_, err := engine.Table("element").Where("element.scope_id = ? and element.name = ?", scopeId, name).Get(&element)
	return element, err
}

func createLink(link *Link) error {
	_, err := engine.Table("link").Insert(link)
	return err
}

func readLinkById(id int64) (Link, error) {
	var link Link
	_, err := engine.Table("link").Id(id).Get(&link)
	return link, err
}

func readLinksByScopeFisrt(scopeId int64) ([]RspLink, error) {
	var links []RspLink
	err := engine.Table("link").Where("link.scope_id1 = ?", scopeId).Join("INNER", []string{"element", "element1"}, "link.element_code1 = element1.code").Join("INNER", []string{"element", "element2"}, "link.element_code2 = element2.code").Find(&links)
	return links, err
}

func readLinksByScopeSecond(scopeId int64) ([]RspLink, error) {
	var links []RspLink
	err := engine.Table("link").Where("link.scope_id2 = ?", scopeId).Join("INNER", []string{"element", "element1"}, "link.element_code1 = element1.code").Join("INNER", []string{"element", "element2"}, "link.element_code2 = element2.code").Find(&links)
	return links, err
}

func readLinksElementConfirmed() ([]RspLink, error) {
	var links []RspLink
	err := engine.Table("link").Where("link.status <> 0").Join("INNER", []string{"element", "element1"}, "link.element_code1 = element1.code").Join("INNER", []string{"element", "element2"}, "link.element_code2 = element2.code").Find(&links)
	return links, err
}

func readLinksConfirmed() ([]Link, error) {
	var links []Link
	err := engine.Table("link").Where("link.status <> 0").Find(&links)
	return links, err
}

func countLinkByCode(code1, code2 string) (int, error) {
	var link Link
	total, err := engine.Table("link").Where("link.element_code1 = ? and link.element_code2 = ?", code1, code2).Count(&link)
	return int(total), err
}

func updateLink(link *Link) error {
	_, err := engine.Table("link").Id(link.Id).Update(link)
	return err
}

func updateLinkStatus(link *Link) error {
	_, err := engine.Table("link").Cols("status").Id(link.Id).Update(link)
	return err
}
