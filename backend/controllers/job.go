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

// TODO: Implement personal filtering personal:true in type
// then filter out open jobs that have skills that user has (if skills are empty, return all) ,
// that hours per week is <= user hours per week,
// that rate is >= user rate (if hourly) include fixed?)
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

	userID := c.Ctx.Input.GetData("id").(int)
	user, err := models.GetUserById(userID)
	if user == nil || err != nil {
		c.Ctx.Output.SetStatus(http.StatusUnauthorized)
		c.Ctx.Output.JSON(map[string]string{"error": "User not found"}, false, false)
		return
	}

	jobID, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": "Invalid job ID"}, false, false)
		return
	}

	job, err := models.GetJobByID(jobID)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusNotFound)
		c.Ctx.Output.JSON(map[string]string{"error": "Job not found"}, false, false)
		return
	}

	if job.Status != "open" {
		if user.Role == "freelancer" {
			application, err := models.GetApplicationByUserAndJob(user.Id, jobID)
			if err != nil || application == nil {
				c.Ctx.Output.SetStatus(http.StatusNotFound)
				c.Ctx.Output.JSON(map[string]string{"error": "Job not found"}, false, false)
				return
			}
			// Allow freelancers to access their applied job
		} else if user.Role == "client" && job.Client.Id == user.Id {
			// Allow client to access their own job
		} else {
			c.Ctx.Output.SetStatus(http.StatusNotFound)
			c.Ctx.Output.JSON(map[string]string{"error": "Job not found"}, false, false)
			return
		}
	}

	applicationID := 0
	if user.Role == "freelancer" {
		application, err := models.GetApplicationByUserAndJob(userID, jobID)
		if err == nil && application != nil {
			applicationID = application.Id
		}
	}

	var skillList []types.Skill
	for _, skill := range job.Skills {
		skillList = append(skillList, types.Skill{
			Id:   skill.Id,
			Name: skill.Name,
		})
	}

	jobInfo := types.JobInfo{
		ID:            job.Id,
		Title:         job.Title,
		Description:   job.Description,
		Type:          job.Type,
		Rate:          job.Rate,
		Amount:        job.Amount,
		Length:        job.Length,
		HoursPerWeek:  job.HoursPerWeek,
		ClientID:      job.Client.Id,
		Skills:        skillList,
		ApplicationID: applicationID,
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

	jobID, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": "Invalid job ID"}, false, false)
		return
	}

	id := c.Ctx.Input.GetData("id").(int)
	user, err := models.GetUserById(id)
	if user == nil || err != nil {
		c.Ctx.Output.SetStatus(http.StatusUnauthorized)
		c.Ctx.Output.JSON(map[string]string{"error": "User not found"}, false, false)
		return
	}

	job, err := models.GetJobByID(jobID)

	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusNotFound)
		c.Ctx.Output.JSON(map[string]string{"error": "Job not found"}, false, false)
		return
	}

	if job.Status != "open" {
		c.Ctx.Output.SetStatus(http.StatusNotFound)
		c.Ctx.Output.JSON(map[string]string{"error": "It is only possible to update open jobs"}, false, false)
		return
	}

	if job.Client.Id != user.Id {
		c.Ctx.Output.SetStatus(http.StatusForbidden)
		c.Ctx.Output.JSON(map[string]string{"error": "Clients are only allowed to update their own jobs"}, false, false)
		return
	}

	updateJobRequest, err := validators.UpdateJobValidator(c.Ctx.Input.RequestBody)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": err.Error()}, false, false)
		return
	}
	if updateJobRequest.Title != "" {
		job.Title = updateJobRequest.Title
	}
	if updateJobRequest.Description != "" {
		job.Description = updateJobRequest.Description
	}
	if updateJobRequest.Type != "" {
		job.Type = updateJobRequest.Type
	}
	if updateJobRequest.Rate != "" {
		job.Rate = updateJobRequest.Rate
	}
	if updateJobRequest.Amount != 0 {
		job.Amount = updateJobRequest.Amount
	}
	if updateJobRequest.Length != "" {
		job.Length = updateJobRequest.Length
	}
	if updateJobRequest.HoursPerWeek != "" {
		job.HoursPerWeek = updateJobRequest.HoursPerWeek
	}

	err = models.UpdateJob(job, updateJobRequest.Skills)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Ctx.Output.JSON(map[string]string{"error": "Error updating job"}, false, false)
		return
	}

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = map[string]string{"message": "Job data updated successfully"}
	c.ServeJSON()

}

func (c *JobController) DeleteClientJobHandler() {

	jobID, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": "Invalid job ID"}, false, false)
		return
	}

	id := c.Ctx.Input.GetData("id").(int)
	user, err := models.GetUserById(id)
	if user == nil || err != nil {
		c.Ctx.Output.SetStatus(http.StatusUnauthorized)
		c.Ctx.Output.JSON(map[string]string{"error": "User not found"}, false, false)
		return
	}

	job, err := models.GetJobByID(jobID)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusNotFound)
		c.Ctx.Output.JSON(map[string]string{"error": "Job not found"}, false, false)
		return
	}

	if job.Status != "open" {
		c.Ctx.Output.SetStatus(http.StatusNotFound)
		c.Ctx.Output.JSON(map[string]string{"error": "It is only possible to delete open jobs"}, false, false)
		return
	}

	if job.Client.Id != user.Id {
		c.Ctx.Output.SetStatus(http.StatusForbidden)
		c.Ctx.Output.JSON(map[string]string{"error": "Clients are only allowed to delete their own jobs"}, false, false)
		return
	}

	err = models.DeleteJobByID(jobID)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Ctx.Output.JSON(map[string]string{"error": "Error deleting job"}, false, false)
		return
	}

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = map[string]string{"message": "Job deleted successfully"}
	c.ServeJSON()

}

func (c *JobController) GetClientJobsHandler() {

	userID := c.Ctx.Input.GetData("id").(int)
	user, err := models.GetUserById(userID)
	if user == nil || err != nil {
		c.Ctx.Output.SetStatus(http.StatusUnauthorized)
		c.Ctx.Output.JSON(map[string]string{"error": "User not found"}, false, false)
		return
	}

	if user.Role != "client" {
		c.Ctx.Output.SetStatus(http.StatusForbidden)
		c.Ctx.Output.JSON(map[string]string{"error": "Only clients can access their jobs"}, false, false)
		return
	}

	jobs, err := models.GetJobsByClientID(userID)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Ctx.Output.JSON(map[string]string{"error": "Error fetching jobs"}, false, false)
		return
	}

	var jobList []types.ClientJobInfo
	for _, job := range jobs {
		var skillList []types.Skill
		for _, skill := range job.Skills {
			skillList = append(skillList, types.Skill{
				Id:   skill.Id,
				Name: skill.Name,
			})
		}

		// Get application count for this job
		applicationCount, err := models.GetApplicationCountForJob(job.Id)
		if err != nil {
			c.Ctx.Output.SetStatus(http.StatusInternalServerError)
			c.Ctx.Output.JSON(map[string]string{"error": "Error fetching application count"}, false, false)
			return
		}

		jobInfo := types.ClientJobInfo{
			ID:               job.Id,
			Title:            job.Title,
			Description:      job.Description,
			Type:             job.Type,
			Rate:             job.Rate,
			Amount:           job.Amount,
			Length:           job.Length,
			HoursPerWeek:     job.HoursPerWeek,
			Status:           job.Status,
			ClientID:         job.Client.Id,
			Skills:           skillList,
			ApplicationCount: applicationCount,
		}

		jobList = append(jobList, jobInfo)
	}

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = jobList
	c.ServeJSON()
}

func (c *JobController) GetClientJobHandler() {

	jobID, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": "Invalid job ID"}, false, false)
		return
	}

	userID := c.Ctx.Input.GetData("id").(int)
	user, err := models.GetUserById(userID)
	if user == nil || err != nil {
		c.Ctx.Output.SetStatus(http.StatusUnauthorized)
		c.Ctx.Output.JSON(map[string]string{"error": "User not found"}, false, false)
		return
	}
	if user.Role != "client" {
		c.Ctx.Output.SetStatus(http.StatusForbidden)
		c.Ctx.Output.JSON(map[string]string{"error": "Only clients can access their job details"}, false, false)
		return
	}

	job, err := models.GetJobByID(jobID)
	if err != nil || job.Client.Id != user.Id {
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

	applications, err := models.GetApplicationsByJobID(jobID)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Ctx.Output.JSON(map[string]string{"error": "Error fetching applications"}, false, false)
		return
	}

	var applicationList []types.Application
	for _, application := range applications {

		attachment, err := models.GetAttachmentByApplicationID(application.Id)
		if err != nil {
			c.Ctx.Output.SetStatus(http.StatusInternalServerError)
			c.Ctx.Output.JSON(map[string]string{"error": "Error fetching attachment"}, false, false)
			return
		}

		var attachmentInfo *types.Attachment
		if attachment != nil {
			attachmentInfo = &types.Attachment{
				ID:            attachment.Id,
				ApplicationID: attachment.Application.Id,
				FileName:      attachment.FileName,
				FilePath:      attachment.FilePath,
				CreatedAt:     attachment.CreatedAt,
			}
		}

		applicationList = append(applicationList, types.Application{
			ID:              application.Id,
			UserID:          application.User.Id,
			JobID:           application.Job.Id,
			Description:     application.Description,
			RejectionReason: application.RejectionReason,
			Status:          application.Status,
			CreatedAt:       application.CreatedAt,
			Attachment:      attachmentInfo,
		})
	}

	jobInfo := types.ClientJobDetailedInfo{
		ID:           job.Id,
		Title:        job.Title,
		Description:  job.Description,
		Type:         job.Type,
		Rate:         job.Rate,
		Amount:       job.Amount,
		Length:       job.Length,
		HoursPerWeek: job.HoursPerWeek,
		Status:       job.Status,
		ClientID:     job.Client.Id,
		Skills:       skillList,
		Applications: applicationList,
	}

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = jobInfo
	c.ServeJSON()
}

func (c *JobController) GetFreelancerJobsHandler() {

	userID := c.Ctx.Input.GetData("id").(int)
	user, err := models.GetUserById(userID)
	if user == nil || err != nil {
		c.Ctx.Output.SetStatus(http.StatusUnauthorized)
		c.Ctx.Output.JSON(map[string]string{"error": "User not found"}, false, false)
		return
	}

	if user.Role != "freelancer" {
		c.Ctx.Output.SetStatus(http.StatusForbidden)
		c.Ctx.Output.JSON(map[string]string{"error": "Only freelancers can access their jobs"}, false, false)
		return
	}

	jobs, err := models.GetJobsByFreelancerID(userID)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Ctx.Output.JSON(map[string]string{"error": "Error fetching jobs"}, false, false)
		return
	}

	var jobList []types.FreelancerJobInfo
	for _, job := range jobs {
		var skillList []types.Skill
		for _, skill := range job.Skills {
			skillList = append(skillList, types.Skill{
				Id:   skill.Id,
				Name: skill.Name,
			})
		}

		application, err := models.GetApplicationByUserAndJob(userID, job.Id)
		if err != nil || application == nil {
			c.Ctx.Output.SetStatus(http.StatusInternalServerError)
			c.Ctx.Output.JSON(map[string]string{"error": "Error fetching application"}, false, false)
			return
		}

		jobInfo := types.FreelancerJobInfo{
			ID:            job.Id,
			Title:         job.Title,
			Description:   job.Description,
			Type:          job.Type,
			Rate:          job.Rate,
			Amount:        job.Amount,
			Length:        job.Length,
			HoursPerWeek:  job.HoursPerWeek,
			Status:        job.Status,
			ClientID:      job.Client.Id,
			Skills:        skillList,
			ApplicationID: application.Id,
		}

		jobList = append(jobList, jobInfo)
	}

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = jobList
	c.ServeJSON()

}

func (c *JobController) GetFreelancerJobHandler() {

	jobID, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": "Invalid job ID"}, false, false)
		return
	}

	userID := c.Ctx.Input.GetData("id").(int)
	user, err := models.GetUserById(userID)
	if user == nil || err != nil {
		c.Ctx.Output.SetStatus(http.StatusUnauthorized)
		c.Ctx.Output.JSON(map[string]string{"error": "User not found"}, false, false)
		return
	}

	if user.Role != "freelancer" {
		c.Ctx.Output.SetStatus(http.StatusForbidden)
		c.Ctx.Output.JSON(map[string]string{"error": "Only freelancers can access their jobs"}, false, false)
		return
	}

	job, err := models.GetJobByID(jobID)

	if err != nil || job.Freelancer == nil || job.Freelancer.Id != user.Id {
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

	application, err := models.GetApplicationByUserAndJob(userID, job.Id)
	if err != nil || application == nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Ctx.Output.JSON(map[string]string{"error": "Error fetching application"}, false, false)
		return
	}

	jobInfo := types.FreelancerJobInfo{
		ID:            job.Id,
		Title:         job.Title,
		Description:   job.Description,
		Type:          job.Type,
		Rate:          job.Rate,
		Amount:        job.Amount,
		Length:        job.Length,
		HoursPerWeek:  job.HoursPerWeek,
		Status:        job.Status,
		ClientID:      job.Client.Id,
		Skills:        skillList,
		ApplicationID: application.Id,
	}

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = jobInfo
	c.ServeJSON()

}

// Admin function
func (c *JobController) DeleteJobHandler() {

	jobID, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": "Invalid job ID"}, false, false)
		return
	}

	id := c.Ctx.Input.GetData("id").(int)
	user, err := models.GetUserById(id)
	if user == nil || err != nil {
		c.Ctx.Output.SetStatus(http.StatusUnauthorized)
		c.Ctx.Output.JSON(map[string]string{"error": "User not found"}, false, false)
		return
	}

	err = models.DeleteJobByID(jobID)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Ctx.Output.JSON(map[string]string{"error": "Error deleting job"}, false, false)
		return
	}

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = map[string]string{"message": "Job deleted successfully"}
	c.ServeJSON()
}
