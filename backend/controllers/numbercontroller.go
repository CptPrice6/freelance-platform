package controllers

import (
	"backend/models"

	"github.com/beego/beego/v2/server/web"
)

type NumberController struct {
	web.Controller
}

// GET /random - Returns a random number
func (c *NumberController) Get() {
	number, err := models.GetRandomNumber()
	if err != nil {
		c.Data["json"] = map[string]string{"error": err.Error()}
	} else {
		c.Data["json"] = map[string]int{"random_number": number}
	}
	c.ServeJSON()
}
