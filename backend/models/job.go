package models

import (
	"backend/types"
	"fmt"

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
	HoursPerWeek string   `orm:"size(30)"`               // <20, 20-40, 40-60, 60-80, 80+ ( hours )
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

func UpdateJob(job *Job, skills []*types.Skill) error {
	o := orm.NewOrm()

	// Update the job basic information
	_, err := o.Update(job)
	if err != nil {
		return err
	}

	// Get M2M relationship handler
	m2m := o.QueryM2M(job, "Skills")
	if m2m == nil {
		return fmt.Errorf("failed to get M2M relationship")
	}

	_, err = m2m.Clear()
	if err != nil {
		return err
	}

	// Add new skills if provided
	if len(skills) > 0 {
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

func GetOpenJobs() ([]Job, error) {
	o := orm.NewOrm()
	var jobs []Job

	_, err := o.QueryTable(new(Job)).Filter("Status", "open").All(&jobs)
	if err != nil {
		return nil, err
	}

	for i := range jobs {
		_, err := o.LoadRelated(&jobs[i], "Skills")
		if err != nil {
			return nil, err
		}
	}

	return jobs, nil
}

func GetJobByID(jobID int) (*Job, error) {
	o := orm.NewOrm()
	job := Job{Id: jobID}

	err := o.Read(&job)
	if err != nil {
		return nil, err
	}

	_, err = o.LoadRelated(&job, "Skills")
	if err != nil {
		return nil, err
	}

	return &job, nil
}

func DeleteJobByID(jobID int) error {

	o := orm.NewOrm()
	job := Job{Id: jobID}

	_, err := o.Delete(&job)
	if err != nil {
		return err
	}

	return nil
}

func GetJobsByClientID(clientID int) ([]Job, error) {
	o := orm.NewOrm()
	var jobs []Job

	_, err := o.QueryTable(new(Job)).Filter("Client__Id", clientID).All(&jobs)
	if err != nil {
		return nil, err
	}

	for i := range jobs {
		_, err := o.LoadRelated(&jobs[i], "Skills")
		if err != nil {
			return nil, err
		}
	}

	return jobs, nil
}

func GetJobsByFreelancerID(freelancerID int) ([]Job, error) {
	o := orm.NewOrm()
	var jobs []Job

	_, err := o.QueryTable(new(Job)).Filter("Freelancer__Id", freelancerID).All(&jobs)
	if err != nil {
		return nil, err
	}

	for i := range jobs {
		_, err := o.LoadRelated(&jobs[i], "Skills")
		if err != nil {
			return nil, err
		}
	}

	return jobs, nil
}
