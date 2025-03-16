package models

import "github.com/beego/beego/v2/client/orm"

type Skill struct {
	Id   int    `orm:"pk;auto"`
	Name string `orm:"size(50);unique"`
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
	_, err := o.QueryTable(new(Skill)).OrderBy("id").All(&skills)
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

func DeleteSkillByID(skillID int) error {
	o := orm.NewOrm()

	skill := Skill{Id: skillID}

	// Delete the skill by ID
	_, err := o.Delete(&skill)
	if err != nil {
		return err
	}

	return nil
}

func UpdateSkill(skill *Skill) error {
	o := orm.NewOrm()

	// Update skill
	_, err := o.Update(skill)
	if err != nil {
		return err
	}

	return nil
}

func CreateSkill(name string) error {
	o := orm.NewOrm()
	skill := Skill{
		Name: name,
	}

	_, err := o.Insert(&skill)
	if err != nil {
		return err
	}

	return nil
}
