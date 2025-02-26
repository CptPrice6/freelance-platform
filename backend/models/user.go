package models

import (
	"github.com/beego/beego/v2/client/orm"
)

type User struct {
	Id       int    `orm:"pk;auto"`
	Email    string `orm:"unique"`
	Password string
	Role     string `orm:"default(user)"`
	Ban      bool   `orm:"default(false)"`
}

func init() {
	orm.RegisterModel(new(User))
}

func (u *User) TableName() string {
	return "users"
}

func CreateUser(email, password, role string, ban bool) error {
	o := orm.NewOrm()
	user := User{
		Email:    email,
		Password: password,
		Role:     role,
		Ban:      ban,
	}

	_, err := o.Insert(&user)
	return err
}

func GetUserByEmail(email string) (*User, error) {
	o := orm.NewOrm()
	user := User{Email: email}
	err := o.Read(&user, "Email")
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserById(id int) (*User, error) {
	o := orm.NewOrm()
	user := User{Id: id}
	err := o.Read(&user, "Id")
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func IsUserBanned(userID int) (bool, error) {
	o := orm.NewOrm()
	user := User{Id: userID}

	err := o.Read(&user, "Id")
	if err == orm.ErrNoRows {
		return false, nil // User not found, so not banned
	} else if err != nil {
		return false, err // Some other error occurred
	}

	return user.Ban, nil // Return the ban status
}

func BanUserByID(userID int) error {
	o := orm.NewOrm()
	user := User{Id: userID}

	if err := o.Read(&user); err != nil {
		return err
	}

	// Update the ban status
	user.Ban = true
	_, err := o.Update(&user, "Ban")
	return err
}
