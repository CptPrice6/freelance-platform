package models

import (
	"backend/types"

	"github.com/beego/beego/v2/client/orm"
)

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

func CreateJob(client *User, title, description, projectType, rate, length, hoursPerWeek string, amount int, skills []*types.Skill) error {
	o := orm.NewOrm()

	// Create a new Job instance
	job := Job{
		Client:       client,
		Title:        title,
		Description:  description,
		Type:         projectType,
		Rate:         rate,
		Amount:       amount,
		Length:       length,
		HoursPerWeek: hoursPerWeek,
		Status:       "open", // Default status
	}

	_, err := o.Insert(&job)
	if err != nil {
		return err
	}

	// Associate skills if provided
	if len(skills) > 0 {
		m2m := o.QueryM2M(&job, "Skills")
		for _, skillReq := range skills {
			var skill Skill
			err := o.QueryTable("skills").Filter("Id", skillReq.Id).One(&skill)
			if err == nil { // Only add existing skills
				m2m.Add(&skill)
			}
		}
	}

	return nil
}
