package controllers

import (
	"backend/models"
	"backend/validators"
	"net/http"

	"github.com/beego/beego/v2/server/web"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	web.Controller
}

func (c *UserController) GetUserHandler() {
	id := c.Ctx.Input.GetData("id").(int)
	user, err := models.GetUserById(id)
	if user == nil || err != nil {
		c.Ctx.Output.SetStatus(http.StatusUnauthorized)
		c.Ctx.Output.JSON(map[string]string{"error": "User not found"}, false, false)
		return
	}

	response := map[string]interface{}{
		"id":          user.Id,
		"email":       user.Email,
		"role":        user.Role,
		"name":        user.Name,
		"surname":     user.Surname,
		"description": user.Description,
	}

	switch user.Role {
	case "freelancer":
		freelancerData, err := models.GetFreelancerDataByUserID(user.Id)
		if err == nil && freelancerData != nil {
			response["freelancer_data"] = map[string]interface{}{
				"skills":         freelancerData.Skills,
				"hourly_rate":    freelancerData.HourlyRate,
				"work_type":      freelancerData.WorkType,
				"hours_per_week": freelancerData.HoursPerWeek,
			}
		}
	case "client":
		clientData, err := models.GetClientDataByUserID(user.Id)
		if err == nil && clientData != nil {
			response["client_data"] = map[string]interface{}{
				"company_name": clientData.CompanyName,
				"industry":     clientData.Industry,
				"location":     clientData.Location,
			}
		}
	}

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = response
	c.ServeJSON()

}

func (c *UserController) UpdateUserHandler() {
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
	}
	if updateUserRequest.Password != "" {
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
	if updateUserRequest.Name != "" {
		user.Name = updateUserRequest.Name
	}
	if updateUserRequest.Surname != "" {
		user.Surname = updateUserRequest.Surname
	}
	if updateUserRequest.Description != "" {
		user.Description = updateUserRequest.Description
	}

	err = models.UpdateUser(user)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Ctx.Output.JSON(map[string]string{"error": err.Error()}, false, false)
		return
	}

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = map[string]string{"message": "User data updated successfully"}
	c.ServeJSON()

}

func (c *UserController) DeleteUserHandler() {
	id := c.Ctx.Input.GetData("id").(int)
	err := models.DeleteUserByID(id)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": "User not found"}, false, false)
		return
	}

	err = models.DeleteAllRefreshTokensForUser(id)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Ctx.Output.JSON(map[string]string{"error": "Error deleting old refresh token"}, false, false)
		return
	}

	// delete all cascading tables!

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = map[string]string{"message": "User deletion successful"}
	c.ServeJSON()
}
