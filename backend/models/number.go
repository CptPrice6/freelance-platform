package models

import (
	"math/rand"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

// Number represents a database entity
type Number struct {
	Id    int `orm:"pk;auto"`
	Value int
}

func init() {
	orm.RegisterModel(new(Number))
}

// GetRandomNumber fetches a random number from the database
func GetRandomNumber() (int, error) {
	o := orm.NewOrm()
	var numbers []Number

	_, err := o.QueryTable(new(Number)).All(&numbers)
	if err != nil {
		return 0, err
	}

	if len(numbers) == 0 {
		return 0, nil
	}

	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(numbers))

	return numbers[randomIndex].Value, nil
}
