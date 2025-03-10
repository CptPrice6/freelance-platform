package models

import "github.com/beego/beego/v2/client/orm"

type Skill struct {
	Id   int    `orm:"pk;auto"`
	Name string `orm:"unique"`
}

func init() {
	orm.RegisterModel(new(Skill))
}

func (s *Skill) TableName() string {
	return "skills"
}

func GetAllSkills() ([]Skill, error) {
	o := orm.NewOrm()
	var skills []Skill
	_, err := o.QueryTable(new(Skill)).All(&skills)
	return skills, err
}

func GetSkillById(skillID int) (*Skill, error) {
	o := orm.NewOrm()
	skill := Skill{Id: skillID}
	err := o.Read(&skill, "Id")
	if err != nil {
		return nil, err
	}
	return &skill, nil
}
