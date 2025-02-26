package controllers

import (
	"backend/models"
	"backend/validators"
	"net/http"

	"github.com/beego/beego/v2/server/web"
)

type AdminController struct {
	web.Controller
}

func (c *AdminController) LogoutUserHandler() {
	logoutUserRequest, err := validators.LogoutUserValidator(c.Ctx.Input.RequestBody)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": err.Error()}, false, false)
		return
	}

	user, err := models.GetUserById(logoutUserRequest.UserId)
	if user == nil || err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": "User not found"}, false, false)
		return
	}

	err = models.DeleteAllRefreshTokensForUser(logoutUserRequest.UserId)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Ctx.Output.JSON(map[string]string{"error": "Error deleting old refresh token"}, false, false)
		return
	}

	c.Ctx.Output.SetStatus(http.StatusCreated)
	c.Data["json"] = map[string]string{"message": "User log out successfull"}
	c.ServeJSON()

}
