package controllers

import (
	"backend/models"
	"backend/types"
	"backend/validators"
	"net/http"
	"strconv"

	"github.com/beego/beego/v2/server/web"
)

type AdminController struct {
	web.Controller
}

func (c *AdminController) DeleteUserHandler() {

	idStr := c.Ctx.Input.Param(":id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": "Invalid user ID"}, false, false)
		return
	}

	err = models.DeleteUserByID(userID)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": "User not found"}, false, false)
		return
	}

	err = models.DeleteAllRefreshTokensForUser(userID)
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

func (c *AdminController) UpdateUserHandler() {
	idStr := c.Ctx.Input.Param(":id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": "Invalid user ID"}, false, false)
		return
	}

	user, err := models.GetUserById(userID)
	if user == nil || err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": "User not found"}, false, false)
		return
	}

	updateUserRequest, err := validators.UpdateUserValidatorAdmin(c.Ctx.Input.RequestBody)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": err.Error()}, false, false)
		return
	}

	if updateUserRequest.Role != "" {
		user.Role = updateUserRequest.Role
	}
	user.Ban = updateUserRequest.Ban

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

func (c *AdminController) GetUsersHandler() {

	users, err := models.GetUsers()
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Ctx.Output.JSON(map[string]string{"error": err.Error()}, false, false)
		return
	}

	var usersResponse []types.UserResponseForAdmins

	for _, user := range users {

		response := types.UserResponseForAdmins{
			ID:        user.Id,
			Email:     user.Email,
			Role:      user.Role,
			Name:      user.Name,
			Surname:   user.Surname,
			Ban:       user.Ban,
			CreatedAt: user.CreatedAt,
		}

		switch user.Role {
		case "freelancer":
			freelancerData, err := models.GetFreelancerDataByUserID(user.Id)
			if err == nil && freelancerData != nil {
				var skillList []types.Skill

				for _, skill := range freelancerData.Skills {
					skillList = append(skillList, types.Skill{
						Id:   skill.Id,
						Name: skill.Name,
					})
				}

				response.FreelancerData = &types.FreelancerData{
					Title:        freelancerData.Title,
					Description:  freelancerData.Description,
					Skills:       skillList,
					HourlyRate:   freelancerData.HourlyRate,
					WorkType:     freelancerData.WorkType,
					HoursPerWeek: freelancerData.HoursPerWeek,
				}
			}
		case "client":
			clientData, err := models.GetClientDataByUserID(user.Id)
			if err == nil && clientData != nil {
				response.ClientData = &types.ClientData{
					Description: clientData.Description,
					CompanyName: clientData.CompanyName,
					Industry:    clientData.Industry,
					Location:    clientData.Location,
				}
			}
		}

		usersResponse = append(usersResponse, response)
	}

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = usersResponse
	c.ServeJSON()

}
