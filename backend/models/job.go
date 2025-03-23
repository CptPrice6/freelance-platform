package models

import "github.com/beego/beego/v2/client/orm"

type Job struct {
	Id           int      `orm:"pk;auto"`
	Client       *User    `orm:"rel(fk);on_delete(cascade)"`
	Freelancer   *User    `orm:"rel(fk);on_delete(cascade);null"`
	Title        string   `orm:"size(30)"`
	Description  string   `orm:"type(text)"`
	Type         string   `orm:"size(30)"` // ongoing, one-time
	Rate         string   `orm:"size(30)"` // hourly, fixed
	Amount       int      // if hourly, amount per hour, if fixed, total amount
	Length       string   `orm:"size(30)"`               // <1, 1-3, 3-6, 6-12, 12+ ( months )
	HoursPerWeek string   `orm:"size(30)"`               // <10, 10-20, 20-40, 40-60, 80+ ( hours )
	Status       string   `orm:"size(30);default(open)"` // open, in-progress, completed
	Skills       []*Skill `orm:"rel(m2m);rel_table(job_skills);on_delete(cascade)"`
}

func init() {
	orm.RegisterModel(new(Job))
}

func (s *Job) TableName() string {
	return "jobs"
}
