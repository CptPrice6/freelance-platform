package models

import (
	"errors"

	"github.com/beego/beego/v2/client/orm"
)

type FreelancerData struct {
	Id           int      `orm:"pk;auto"`
	User         *User    `orm:"rel(fk);on_delete(cascade);unique"`
	HourlyRate   float64  `orm:"null"`
	WorkType     string   `orm:"size(50);null"`
	HoursPerWeek int      `orm:"null"`
	Skills       []*Skill `orm:"rel(m2m);rel_table(freelancer_skills)"`
}

func init() {
	orm.RegisterModel(new(FreelancerData))
}

func (f *FreelancerData) TableName() string {
	return "freelancer_data"
}

func GetFreelancerDataByUserID(userID int) (*FreelancerData, error) {
	o := orm.NewOrm()
	freelancerData := FreelancerData{User: &User{Id: userID}}
	err := o.Read(&freelancerData, "User")
	if err == orm.ErrNoRows {
		return nil, nil
	}

	_, err = o.LoadRelated(&freelancerData, "Skills")
	if err != nil {
		return nil, err
	}

	return &freelancerData, err
}

func CreateFreelancerData(userID int) error {
	o := orm.NewOrm()

	user := User{Id: userID}
	err := o.Read(&user)
	if err == orm.ErrNoRows {
		return errors.New("user not found")
	}

	if user.Role != "freelancer" {
		return errors.New("user is not a freelancer")
	}

	freelancerData := &FreelancerData{
		User: &user,
	}

	_, err = o.Insert(freelancerData)
	if err != nil {
		return err
	}

	return nil
}

func UpdateFreelancerData(freelancerData *FreelancerData) error {
	o := orm.NewOrm()

	_, err := o.Update(freelancerData)
	if err != nil {
		return err
	}

	return nil
}

func AddSkillToFreelancerData(freelancerData *FreelancerData, skill *Skill) error {
	o := orm.NewOrm()

	_, err := o.QueryM2M(freelancerData, "Skills").Add(skill)
	if err != nil {
		return err
	}

	return nil
}
func DeleteSkillFromFreelancerData(freelancerData *FreelancerData, skill *Skill) error {
	o := orm.NewOrm()

	_, err := o.QueryM2M(freelancerData, "Skills").Remove(skill)
	if err != nil {
		return err
	}

	return nil
}
