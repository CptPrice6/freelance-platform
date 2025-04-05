package models

import (
	"errors"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Attachment struct {
	Id          int          `orm:"pk;auto"`
	Application *Application `orm:"rel(fk);on_delete(cascade)"`
	FileName    string       `orm:"size(255)"` // Original filename
	FilePath    string       `orm:"size(255)"` // Path with unique name
	CreatedAt   time.Time    `orm:"auto_now_add;type(timestamp)"`
}

func init() {
	orm.RegisterModel(new(Attachment))
}

func (s *Attachment) TableName() string {
	return "attachments"
}

func GetAttachmentByApplicationID(applicationID int) (*Attachment, error) {
	o := orm.NewOrm()
	var attachment Attachment

	err := o.QueryTable(new(Attachment)).
		Filter("Application__Id", applicationID).
		One(&attachment)

	if err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &attachment, nil
}

func CreateAttachment(applicationId int, fileName, filePath string) error {
	o := orm.NewOrm()

	exists := o.QueryTable(new(Application)).Filter("Id", applicationId).Exist()
	if !exists {
		return errors.New("application not found")
	}

	attachment := Attachment{
		Application: &Application{Id: applicationId},
		FileName:    fileName,
		FilePath:    filePath,
		CreatedAt:   time.Now(),
	}
	_, err := o.Insert(&attachment)
	if err != nil {
		return err
	}

	return nil
}

func GetAttachmentByID(attachmentID int) (*Attachment, error) {
	o := orm.NewOrm()
	attachment := Attachment{Id: attachmentID}

	err := o.Read(&attachment, "Id")

	if err != nil {
		return nil, err
	}

	if _, err := o.LoadRelated(&attachment, "Application"); err != nil {
		return nil, err
	}

	return &attachment, nil
}

func DeleteAttachmentByID(attachmentID int) error {

	o := orm.NewOrm()
	attachment := Attachment{Id: attachmentID}

	_, err := o.Delete(&attachment)
	if err != nil {
		return err
	}

	return nil
}
