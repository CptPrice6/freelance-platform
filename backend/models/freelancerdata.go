package models

import (
	"errors"
	"fmt"

	"github.com/beego/beego/v2/client/orm"
)

type FreelancerData struct {
	Id           int      `orm:"pk;auto"`
	User         *User    `orm:"rel(fk);on_delete(cascade);unique"`
	Title        string   `orm:"size(30);null"`
	Description  string   `orm:"type(text);null"`
	HourlyRate   float64  `orm:"null"`
	HoursPerWeek string   `orm:"size(30);null"` // <20, 20-40, 40-60, 60-80, 80+ ( hours )
	Skills       []*Skill `orm:"rel(m2m);rel_table(freelancer_skills);on_delete(cascade)"`
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

	affectedRows, err := o.QueryM2M(freelancerData, "Skills").Remove(skill)
	if err != nil {
		return err
	}

	if affectedRows == 0 {
		return fmt.Errorf("skill not found")
	}

	return nil
}
