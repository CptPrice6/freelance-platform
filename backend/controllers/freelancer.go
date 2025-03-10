package controllers

import (
	"backend/models"
	"backend/validators"
	"net/http"

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
