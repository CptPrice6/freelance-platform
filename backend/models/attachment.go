package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Attachment struct {
	Id          int          `orm:"pk;auto"`
	Application *Application `orm:"rel(fk);on_delete(cascade)"`
	FileName    string       `orm:"size(255)"` // Original filename
	FilePath    string       `orm:"size(255)"` // Path with generated UUID name
	CreatedAt   time.Time    `orm:"auto_now_add;type(timestamp)"`
}

func init() {
	orm.RegisterModel(new(Attachment))
}

func (s *Attachment) TableName() string {
	return "attachments"
}
