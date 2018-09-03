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

func readScopes() ([]RspScope, error) {
	var scopes []RspScope
	err := engine.Table("scope").Asc("scope.id").Find(&scopes)
	return scopes, err
}

func readScopesExceptOne(name string) ([]RspScope, error) {
	var scopes []RspScope
	err := engine.Table("scope").Where("scope.Name <> ?", name).Asc("scope.id").Find(&scopes)
	return scopes, err
}

func createLink(link *Link) error {
	_, err := engine.Table("link").Insert(link)
	return err
}

func readLinksByScopeFisrt(scopeName string) ([]RspLink, error) {
	var links []RspLink
	err := engine.Table("link").Join("INNER", "scope", "link.scope_id_1 = scope.id").Where("scope.name = ?", scopeName).Find(&links)
	return links, err
}

func readLinksByScopeSecond(scopeName string) ([]RspLink, error) {
	var links []RspLink
	err := engine.Table("link").Join("INNER", "scope", "link.scope_id_2 = scope.id").Where("scope.name = ?", scopeName).Find(&links)
	return links, err
}

func countLinkByCode(code1, code2 string) (int, error) {
	var link Link
	total, err := engine.Table("link").Where("link.element_code_1 = ? and link.element_code_2 = ?", code1, code2).Count(&link)
	return int(total), err
}

func updateLink(link *Link) error {
	_, err := engine.Table("link").Where("link.id = ?", link.Id).Update(link)
	return err
}
