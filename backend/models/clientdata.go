package models

import "github.com/beego/beego/v2/client/orm"

type ClientData struct {
	Id          int    `orm:"pk;auto"`
	User        *User  `orm:"rel(fk);on_delete(cascade);unique"`
	CompanyName string `orm:"size(255);null"`
	Industry    string `orm:"size(100);null"`
	Location    string `orm:"size(255);null"`
}

func init() {
	orm.RegisterModel(new(ClientData))
}

func (u *ClientData) TableName() string {
	return "client_data"
}
