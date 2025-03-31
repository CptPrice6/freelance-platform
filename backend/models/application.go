package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Application struct {
	Id              int       `orm:"pk;auto"`
	User            *User     `orm:"rel(fk);on_delete(cascade)"`
	Job             *Job      `orm:"rel(fk);on_delete(cascade)"`
	Description     string    `orm:"type(text);null"`
	RejectionReason string    `orm:"type(text);null"`           // Reason for rejection, if applicable
	Status          string    `orm:"size(30);default(pending)"` // "pending", "accepted", "rejected"
	CreatedAt       time.Time `orm:"auto_now_add;type(datetime)"`
}

func (s *Application) TableUnique() [][]string {
	return [][]string{
		{"User", "Job"}, // Ensures only one application per user per job
	}
}

func init() {
	orm.RegisterModel(new(Application))
}

func (s *Application) TableName() string {
	return "applications"
}

func GetApplicationByUserAndJob(userID, jobID int) (*Application, error) {
	o := orm.NewOrm()
	var application Application
	err := o.QueryTable("applications").Filter("User__Id", userID).Filter("Job__Id", jobID).One(&application)
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &application, nil
}
