package models

import "github.com/beego/beego/v2/client/orm"

type FreelancerData struct {
	Id           int      `orm:"pk;auto"`
	User         *User    `orm:"rel(fk);on_delete(cascade);unique"`
	HourlyRate   float64  `orm:"null"`
	WorkType     string   `orm:"size(50);null"`
	HoursPerWeek int      `orm:"null"`
	Skill        []*Skill `orm:"rel(m2m);rel_table(freelancer_skills)"`
}

func init() {
	orm.RegisterModel(new(FreelancerData))
}

func (u *FreelancerData) TableName() string {
	return "freelancer_data"
}
