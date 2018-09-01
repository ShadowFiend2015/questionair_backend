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

type Type struct {
	Id   int64  `json:"id" form:"id" query:"id" validate:"eq=0"`
	Name string `json:"name" form:"name" query:"name" validate:"required" xorm:"notnull"`
	Code string `json:"code" form:"code" query:"code" validate:"required" xorm:"notnull"`
}

type RspType struct {
	Type `xorm:"extends"`
}

type Scope struct {
	Id     int64  `json:"id" form:"id" query:"id" validate:"eq=0"`
	TypeId int64  `json:"type_id" form:"type_id" query:"type_id" validate:"required" xorm:"notnull"`
	Name   string `json:"name" form:"name" query:"name" validate:"required" xorm:"notnull"`
	Code   string `json:"code" form:"code" query:"code" validate:"required" xorm:"notnull"`
}

type RspScope struct {
	Scope `xorm:"extends"`
}

type Link struct {
	Id       int64 `json:"id" form:"id" query:"id" validate:"eq=0"`
	TypeId1  int64 `json:"type_id_1" form:"type_id_1" query:"type_id_1" validate:"required" xorm:"notnull"`
	TypeId2  int64 `json:"type_id_2" form:"type_id_2" query:"type_id_2" validate:"required" xorm:"notnull"`
	ScopeId1 int64 `json:"scope_id_1" form:"scope_id_1" query:"scope_id_1" validate:"required" xorm:"notnull"`
	ScopeId2 int64 `json:"scope_id_2" form:"scope_id_2" query:"scope_id_2" validate:"required" xorm:"notnull"`
	Status   int   `json:"status" form:"status" query:"status" validate:"eq=0" xorm:"notnull"`
}

type RspLink struct {
	Link `xorm:"extends"`
}
