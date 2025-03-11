package controllers

import (
	"backend/models"
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
	if updateFreelancerRequest.WorkType != "" {
		freelancerData.WorkType = updateFreelancerRequest.WorkType
	}
	if updateFreelancerRequest.HoursPerWeek > 0 {
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

	var freelancersData []map[string]interface{}

	for _, user := range users {
		freelancerData, err := models.GetFreelancerDataByUserID(user.Id)
		if err != nil || freelancerData == nil {
			continue
		}

		freelancerInfo := map[string]interface{}{
			"id":          user.Id,
			"name":        user.Name,
			"surname":     user.Surname,
			"title":       freelancerData.Title,
			"hourly_rate": freelancerData.HourlyRate,
		}

		freelancersData = append(freelancersData, freelancerInfo)
	}

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = freelancersData
	c.ServeJSON()
}

func (c *FreelancerController) GetFreelancerHandler() {
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
	if user.Role != "freelancer" {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": "Freelancer not found"}, false, false)
		return
	}

	response := map[string]interface{}{
		"id":      user.Id,
		"email":   user.Email,
		"name":    user.Name,
		"surname": user.Surname,
	}

	freelancerData, err := models.GetFreelancerDataByUserID(user.Id)
	if err == nil && freelancerData != nil {
		response["freelancer_data"] = map[string]interface{}{
			"title":          freelancerData.Title,
			"description":    freelancerData.Description,
			"skills":         freelancerData.Skills,
			"hourly_rate":    freelancerData.HourlyRate,
			"work_type":      freelancerData.WorkType,
			"hours_per_week": freelancerData.HoursPerWeek,
		}
	}

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = response
	c.ServeJSON()

}
