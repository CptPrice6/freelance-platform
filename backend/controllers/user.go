package controllers

import (
	"backend/models"
	"backend/validators"
	"net/http"
	"strconv"

	"github.com/beego/beego/v2/server/web"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	web.Controller
}

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

func (c *UserController) Put() {
	id := c.Ctx.Input.GetData("id").(int)
	user, err := models.GetUserById(id)
	if user == nil || err != nil {
		c.Ctx.Output.SetStatus(http.StatusUnauthorized)
		c.Ctx.Output.JSON(map[string]string{"error": "User not found"}, false, false)
		return
	}

	updateUserRequest, err := validators.UpdateUserValidator(c.Ctx.Input.RequestBody)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": err.Error()}, false, false)
		return
	}

	if updateUserRequest.Email != "" {
		user.Email = updateUserRequest.Email
	} else if updateUserRequest.Password != "" {
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(updateUserRequest.Password))
		if err != nil {
			c.Ctx.Output.SetStatus(http.StatusBadRequest)
			c.Ctx.Output.JSON(map[string]string{"error": "Incorrect old password"}, false, false)
			return
		}
		if updateUserRequest.NewPassword != "" {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateUserRequest.NewPassword), bcrypt.DefaultCost)
			if err != nil {
				c.Ctx.Output.SetStatus(http.StatusInternalServerError)
				c.Ctx.Output.JSON(map[string]string{"error": "Password hashing failed"}, false, false)
				return
			}
			user.Password = string(hashedPassword)
		} else {
			c.Ctx.Output.SetStatus(http.StatusBadRequest)
			c.Ctx.Output.JSON(map[string]string{"error": "Missing new password"}, false, false)
			return
		}
	}

	err = models.UpdateUser(user)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Ctx.Output.JSON(map[string]string{"error": "User update failed"}, false, false)
		return
	}

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = map[string]string{"message": "User data updated successfully"}
	c.ServeJSON()

}
