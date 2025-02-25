package models

import (
	"github.com/beego/beego/v2/client/orm"
)

// Number represents a database entity
type User struct {
	Id       int    `orm:"pk;auto"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string
}

func init() {
	orm.RegisterModel(new(User))
}

func CreateUser(email, password, role string) error {
	o := orm.NewOrm()
	user := User{
		Email:    email,
		Password: password, // You should hash this before storing
		Role:     role,
	}

	_, err := o.Insert(&user)
	return err
}

func UserAlreadyExists(email string) bool {
	o := orm.NewOrm()
	user := User{Email: email}
	err := o.Read(&user, "Email")
	return err == nil
}
