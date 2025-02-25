package models

import (
	"github.com/beego/beego/v2/client/orm"
)

// Number represents a database entity
type User struct {
	Id       int    `orm:"pk;auto"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string
}

func init() {
	orm.RegisterModel(new(User))
}
