package models

import "github.com/beego/beego/v2/client/orm"

type Skill struct {
	Id        int    `orm:"pk;auto"`
	SkillName string `orm:"size(50)"`
}

func init() {
	orm.RegisterModel(new(Skill))
}

func (u *Skill) TableName() string {
	return "skills"
}
