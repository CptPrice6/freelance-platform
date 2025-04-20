package models

import (
	"errors"

	"github.com/beego/beego/v2/client/orm"
)

type ClientData struct {
	Id          int    `orm:"pk;auto"`
	User        *User  `orm:"rel(fk);on_delete(cascade);unique"`
	Description string `orm:"type(text);null"`
	CompanyName string `orm:"size(30);null"`
	Industry    string `orm:"size(30);null"`
	Location    string `orm:"size(30);null"`
}

func init() {
	orm.RegisterModel(new(ClientData))
}

func (c *ClientData) TableName() string {
	return "client_data"
}

func GetClientDataByUserID(userID int) (*ClientData, error) {
	o := orm.NewOrm()
	client := ClientData{User: &User{Id: userID}}
	err := o.Read(&client, "User")
	if err == orm.ErrNoRows {
		return nil, nil
	}
	return &client, err
}

func CreateClientData(userID int) error {
	o := orm.NewOrm()

	user := User{Id: userID}
	err := o.Read(&user)
	if err == orm.ErrNoRows {
		return errors.New("user not found")
	}

	clientData := &ClientData{
		User: &user,
	}

	_, err = o.Insert(clientData)
	if err != nil {
		return err
	}

	return nil
}

func UpdateClientData(clientData *ClientData) error {
	o := orm.NewOrm()

	_, err := o.Update(clientData)
	if err != nil {
		return err
	}

	return nil
}
