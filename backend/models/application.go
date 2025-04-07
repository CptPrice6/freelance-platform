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

func GetApplicationCountForJob(jobID int) (int, error) {
	o := orm.NewOrm()
	count, err := o.QueryTable(new(Application)).Filter("Job__Id", jobID).Count()
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func GetApplicationsByJobID(jobID int) ([]Application, error) {
	o := orm.NewOrm()
	var applications []Application

	_, err := o.QueryTable(new(Application)).
		Filter("Job__Id", jobID).All(&applications)

	if err != nil {
		return nil, err
	}

	return applications, nil
}

func CreateApplication(user *User, job *Job, description string) (int, error) {
	o := orm.NewOrm()

	application := Application{
		User:        user,
		Job:         job,
		Description: description,
		Status:      "pending",
		CreatedAt:   time.Now(),
	}

	id, err := o.Insert(&application)
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func GetApplicationByID(applicationID int) (*Application, error) {
	o := orm.NewOrm()
	application := Application{Id: applicationID}

	err := o.Read(&application, "Id")
	if err != nil {
		return nil, err
	}

	_, err = o.LoadRelated(&application, "User")
	if err != nil {
		return nil, err
	}

	_, err = o.LoadRelated(&application, "Job")
	if err != nil {
		return nil, err
	}

	return &application, nil
}

func GetApplicationByUserID(userID int) ([]Application, error) {
	o := orm.NewOrm()
	var applications []Application

	_, err := o.QueryTable(new(Application)).
		Filter("User__Id", userID).All(&applications)

	if err != nil {
		return nil, err
	}

	for i := range applications {
		_, err := o.LoadRelated(&applications[i], "Job")
		if err != nil {
			return nil, err
		}
	}

	return applications, nil
}

func UpdateApplication(application *Application) error {
	o := orm.NewOrm()

	_, err := o.Update(application)
	if err != nil {
		return err
	}

	return nil
}

func DeleteApplicationByID(applicationID int) error {

	o := orm.NewOrm()
	application := Application{Id: applicationID}

	_, err := o.Delete(&application)
	if err != nil {
		return err
	}

	return nil
}
