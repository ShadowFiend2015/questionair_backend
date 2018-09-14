package module

import (
	"time"
)

type omit *struct{}

type RspData struct {
	Total int           `json:"total"`
	Data  []interface{} `json:"data"`
}

type User struct {
	Id        int64     `json:"id" form:"id" query:"id" validate:"eq=0"`
	Account   string    `json:"account" form:"account" query:"account" validate:"required" xorm:"unique notnull"`
	Password  string    `json:"password" form:"password" query:"password" validate:"required" xorm:"notnull"`
	Name      string    `json:"name" form:"name" query:"name" validate:"required" xorm:"notnull"`
	Status    int       `json:"status" form:"status" query:"status" validate:"eq=0" xorm:"notnull"`
	CreatedAt time.Time `json:"-" xorm:"created"`
	UpdatedAt time.Time `json:"-" xorm:"updated"`
}

type RspUser struct {
	User     `xorm:"extends"`
	Password omit `json:"password,omitempty" xorm:"-"`
}

type RspUserCheck struct {
	Pass   bool   `json:"pass"`
	UserId int64  `json:"-"`
	Token  string `json:"token"`
}

type Scope struct {
	Id   int64  `json:"id" form:"id" query:"id" validate:"eq=0"`
	Name string `json:"name" form:"name" query:"name" validate:"required" xorm:"notnull unique"`
	Code string `json:"code" form:"code" query:"code" validate:"required" xorm:"notnull unique"`
}

type RspScope struct {
	Scope `xorm:"extends"`
}

type Element struct {
	Id      int64  `json:"id" form:"id" query:"id" validate:"eq=0"`
	ScopeId int64  `json:"scope_id" form:"scope_id" query:"scope_id" validate:"required" xorm:"notnull unique('1')"`
	Name    string `json:"name" form:"name" query:"name" validate:"required" xorm:"notnull unique('1')"`
	Code    string `json:"code" form:"code" query:"code" validate:"required" xorm:"notnull unique"`
}

type RspElement struct {
	Element `xorm:"extends"`
	Name    string `json:"scope_name" xorm:"name"`
	Code    string `json:"scope_code" xorm:"code"`
}

type RspElementLinked struct {
	ScopeId     int64  `json:"-"`
	ScopeCode   string `json:"scope"`
	ElementCode string `json:"code"`
	ElementName string `json:"name"`
}

// ScopeId1 < ScopeId2
// status: 0 - no confirm, 1 - element_code_1 confirmed, 2 - element_code_2 confirmed, 3 - both confirmed
type Link struct {
	Id           int64  `json:"id" form:"id" query:"id" validate:"eq=0"`
	ScopeId1     int64  `json:"scope_id_1" form:"scope_id_1" query:"scope_id_1" validate:"required" xorm:"notnull"`
	ScopeId2     int64  `json:"scope_id_2" form:"scope_id_2" query:"scope_id_2" validate:"required" xorm:"notnull"`
	ElementCode1 string `json:"element_code_1" form:"element_code_1" query:"element_code_1" validate:"required" xorm:"notnull unique('1')"`
	ElementCode2 string `json:"element_code_2" form:"element_code_2" query:"element_code_2" validate:"required" xorm:"notnull unique('1')"`
	Status       int    `json:"status" form:"status" query:"status" validate:"eq=0" xorm:"notnull"`
}

type RspLink struct {
	Link             `xorm:"extends"`
	LinkElementName1 string `xorm:"name"`
	LinkElementCode1 string `xorm:"code"`
	LinkElementName2 string `xorm:"name"`
	LinkElementCode2 string `xorm:"code"`
}

type RspLinkSelf struct {
	Id           int64  `json:"id"`
	Name         string `json:"name"`
	Code         string `json:"code"`
	LinkName     string `json:"link_name"`
	LinkCode     string `json:"link_code"`
	LinkFullName string `json:"link_full_name"`
	Status       int    `json:"status"`
}

type RspLinkDownload struct {
	Element1 string `json:"a"`
	Element2 string `json:"b"`
}

type RspLinkId struct {
	Id int64 `json:"id"`
}
