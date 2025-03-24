package controllers

import (
	"backend/models"
	"backend/validators"
	"net/http"

	"github.com/beego/beego/v2/server/web"
)

type JobController struct {
	web.Controller
}

func (c *JobController) GetJobsHandler() {
}

func (c *JobController) GetJobHandler() {
}

func (c *JobController) CreateJobHandler() {

	id := c.Ctx.Input.GetData("id").(int)
	user, err := models.GetUserById(id)
	if user == nil || err != nil {
		c.Ctx.Output.SetStatus(http.StatusUnauthorized)
		c.Ctx.Output.JSON(map[string]string{"error": "User not found"}, false, false)
		return
	}

	if user.Role != "client" {
		c.Ctx.Output.SetStatus(http.StatusForbidden)
		c.Data["json"] = map[string]string{"error": "Only clients can create jobs"}
		c.ServeJSON()
		return
	}

	createJobRequest, err := validators.CreateJobValidator(c.Ctx.Input.RequestBody)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": err.Error()}, false, false)
		return
	}

	err = models.CreateJob(user, createJobRequest.Title, createJobRequest.Description, createJobRequest.Type, createJobRequest.Rate, createJobRequest.Length, createJobRequest.HoursPerWeek, createJobRequest.Amount, createJobRequest.Skills)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Ctx.Output.JSON(map[string]string{"error": "Error creating job"}, false, false)
		return
	}

	c.Ctx.Output.SetStatus(http.StatusCreated)
	c.Data["json"] = map[string]string{"message": "Job created successfully"}
	c.ServeJSON()

}

func (c *JobController) UpdateClientJobHandler() {
}

func (c *JobController) DeleteClientJobHandler() {
}

func (c *JobController) GetClientJobsHandler() {
}

func (c *JobController) GetClientJobHandler() {
}

func (c *JobController) GetFreelancerJobsHandler() {

}

func (c *JobController) GetFreelancerJobHandler() {

}
