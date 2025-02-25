package controllers

import (
	"backend/models"
	"net/http"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
)

type UserController struct {
	web.Controller
}

// GET /random - Returns a user from JWT middleware context email
func (c *UserController) Get() {
	email := c.Ctx.Input.GetData("email").(string)
	o := orm.NewOrm()
	user := models.User{Email: email}
	err := o.Read(&user, "Email")
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusUnauthorized)
		c.Ctx.Output.JSON(map[string]string{"error": "User not found"}, false, false)
		return
	}

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Ctx.Output.JSON(map[string]string{"email": user.Email, "role": user.Role}, false, false)

}
