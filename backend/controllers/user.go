package controllers

import (
	"backend/models"
	"net/http"
	"strconv"

	"github.com/beego/beego/v2/server/web"
)

type UserController struct {
	web.Controller
}

// GET /random - Returns a user from JWT middleware context email
func (c *UserController) Get() {
	id := c.Ctx.Input.GetData("id").(int)
	user, err := models.GetUserById(id)
	if user == nil || err != nil {
		c.Ctx.Output.SetStatus(http.StatusUnauthorized)
		c.Ctx.Output.JSON(map[string]string{"error": "User not found"}, false, false)
		return
	}

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Ctx.Output.JSON(map[string]string{"id": strconv.Itoa(user.Id), "email": user.Email, "role": user.Role}, false, false)

}
