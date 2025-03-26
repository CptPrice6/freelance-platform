package controllers

import (
	"backend/models"
	"backend/types"
	"backend/validators"
	"net/http"
	"strconv"

	"github.com/beego/beego/v2/server/web"
)

type JobController struct {
	web.Controller
}

func (c *JobController) GetJobsHandler() {
	jobs, err := models.GetOpenJobs()
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Ctx.Output.JSON(map[string]string{"error": "Error fetching jobs"}, false, false)
		return
	}

	var jobList []types.JobInfo

	for _, job := range jobs {
		var skillList []types.Skill

		for _, skill := range job.Skills {
			skillList = append(skillList, types.Skill{
				Id:   skill.Id,
				Name: skill.Name,
			})
		}

		jobInfo := types.JobInfo{
			ID:           job.Id,
			Title:        job.Title,
			Description:  job.Description,
			Type:         job.Type,
			Rate:         job.Rate,
			Amount:       job.Amount,
			Length:       job.Length,
			HoursPerWeek: job.HoursPerWeek,
			ClientID:     job.Client.Id,
			Skills:       skillList,
		}

		jobList = append(jobList, jobInfo)
	}

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = jobList
	c.ServeJSON()
}

func (c *JobController) GetJobHandler() {
	jobID, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": "Invalid job ID"}, false, false)
		return
	}

	job, err := models.GetJobByID(jobID)
	if err != nil || job.Status != "open" {
		c.Ctx.Output.SetStatus(http.StatusNotFound)
		c.Ctx.Output.JSON(map[string]string{"error": "Job not found"}, false, false)
		return
	}

	var skillList []types.Skill
	for _, skill := range job.Skills {
		skillList = append(skillList, types.Skill{
			Id:   skill.Id,
			Name: skill.Name,
		})
	}

	jobInfo := types.JobInfo{
		ID:           job.Id,
		Title:        job.Title,
		Description:  job.Description,
		Type:         job.Type,
		Rate:         job.Rate,
		Amount:       job.Amount,
		Length:       job.Length,
		HoursPerWeek: job.HoursPerWeek,
		ClientID:     job.Client.Id,
		Skills:       skillList,
	}

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = jobInfo
	c.ServeJSON()
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
