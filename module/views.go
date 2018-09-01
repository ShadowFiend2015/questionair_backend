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
	if err := SyncDatabase(new(User), new(Type), new(Scope), new(Link)); err != nil {
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
