package models

import (
	"github.com/beego/beego/v2/client/orm"
)

// Number represents a database entity
type User struct {
	Id       int `orm:"pk;auto"`
	Name     string
	Email    string
	Password string
}

func init() {
	orm.RegisterModel(new(User))
}
