package controllers

import (
	"backend/models"
	"backend/types"
	"backend/validators"
	"net/http"
	"strconv"

	"github.com/beego/beego/v2/server/web"
)

type FreelancerController struct {
	web.Controller
}

func (c *FreelancerController) UpdateFreelancerDataHandler() {
	userID := c.Ctx.Input.GetData("id").(int)

	freelancerData, err := models.GetFreelancerDataByUserID(userID)
	if err != nil || freelancerData == nil {
		c.Ctx.Output.SetStatus(http.StatusUnauthorized)
		c.Ctx.Output.JSON(map[string]string{"error": "Freelancer profile not found"}, false, false)
		return
	}

	updateFreelancerRequest, err := validators.UpdateFreelancerDataValidator(c.Ctx.Input.RequestBody)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": err.Error()}, false, false)
		return
	}

	if updateFreelancerRequest.HourlyRate > 0 {
		freelancerData.HourlyRate = updateFreelancerRequest.HourlyRate
	}
	if updateFreelancerRequest.HoursPerWeek != "" {
		freelancerData.HoursPerWeek = updateFreelancerRequest.HoursPerWeek
	}
	if updateFreelancerRequest.Description != "" {
		freelancerData.Description = updateFreelancerRequest.Description
	}
	if updateFreelancerRequest.Title != "" {
		freelancerData.Title = updateFreelancerRequest.Title
	}

	err = models.UpdateFreelancerData(freelancerData)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Ctx.Output.JSON(map[string]string{"error": "Error updating freelancer data"}, false, false)
		return
	}

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = map[string]string{"message": "Freelancer data updated successfully"}
	c.ServeJSON()
}

func (c *FreelancerController) GetFreelancersHandler() {
	users, err := models.GetUsersByRole("freelancer")
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Ctx.Output.JSON(map[string]string{"error": "Error fetching freelancers"}, false, false)
		return
	}

	var freelancers []types.FreelancerInfo

	for _, user := range users {
		freelancerData, err := models.GetFreelancerDataByUserID(user.Id)
		if err != nil || freelancerData == nil {
			continue
		}

		freelancerInfo := types.FreelancerInfo{
			ID:         user.Id,
			Name:       user.Name,
			Surname:    user.Surname,
			Title:      freelancerData.Title,
			HourlyRate: freelancerData.HourlyRate,
		}

		freelancers = append(freelancers, freelancerInfo)
	}

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = freelancers
	c.ServeJSON()
}

func (c *FreelancerController) GetFreelancerHandler() {
	idStr := c.Ctx.Input.Param(":id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": "Invalid freelancer ID"}, false, false)
		return
	}

	user, err := models.GetUserById(userID)
	if user == nil || err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": "User not found"}, false, false)
		return
	}
	if user.Role != "freelancer" {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": "Freelancer not found"}, false, false)
		return
	}

	response := types.UserResponse{
		ID:      user.Id,
		Email:   user.Email,
		Role:    user.Role,
		Name:    user.Name,
		Surname: user.Surname,
	}

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
			HoursPerWeek: freelancerData.HoursPerWeek,
		}
	}

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = response
	c.ServeJSON()

}
