package controllers

import (
	"backend/models"
	"backend/types"
	"backend/validators"
	"encoding/base64"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/beego/beego/v2/server/web"
)

type ApplicationController struct {
	web.Controller
}

func (c *ApplicationController) GetFreelancerApplications() {

	userID := c.Ctx.Input.GetData("id").(int)
	user, err := models.GetUserById(userID)
	if user == nil || err != nil {
		c.Ctx.Output.SetStatus(http.StatusUnauthorized)
		c.Ctx.Output.JSON(map[string]string{"error": "User not found"}, false, false)
		return
	}

	if user.Role != "freelancer" {
		c.Ctx.Output.SetStatus(http.StatusForbidden)
		c.Ctx.Output.JSON(map[string]string{"error": "Only freelancer can access their applications"}, false, false)
		return
	}

	applications, err := models.GetApplicationsByUserID(userID)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Ctx.Output.JSON(map[string]string{"error": "Error fetching jobs"}, false, false)
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

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Ctx.Output.JSON(applicationList, false, false)
	c.ServeJSON()

}

func (c *ApplicationController) GetFreelancerApplication() {

	applicationID, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": "Invalid application ID"}, false, false)
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
		c.Ctx.Output.JSON(map[string]string{"error": "Only freelancer can access their application details"}, false, false)
		return
	}

	application, err := models.GetApplicationByID(applicationID)
	if application == nil || err != nil || application.User.Id != userID {
		c.Ctx.Output.SetStatus(http.StatusNotFound)
		c.Ctx.Output.JSON(map[string]string{"error": "Application not found"}, false, false)
		return
	}

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

	applicationInfo := types.Application{
		ID:              application.Id,
		UserID:          application.User.Id,
		JobID:           application.Job.Id,
		Description:     application.Description,
		RejectionReason: application.RejectionReason,
		Status:          application.Status,
		CreatedAt:       application.CreatedAt,
		Attachment:      attachmentInfo,
	}

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Ctx.Output.JSON(applicationInfo, false, false)
	c.ServeJSON()

}

func (c *ApplicationController) SubmitApplication() {

	submitApplicationRequest, err := validators.SubmitApplicationValidator(c.Ctx.Input.RequestBody)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": err.Error()}, false, false)
		return
	}
	job, err := models.GetJobByID(submitApplicationRequest.JobID)
	if job == nil || err != nil {
		c.Ctx.Output.SetStatus(http.StatusNotFound)
		c.Ctx.Output.JSON(map[string]string{"error": "Job not found"}, false, false)
		return
	}

	if job.Status != "open" {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": "Job is not open for applications"}, false, false)
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
		c.Ctx.Output.JSON(map[string]string{"error": "Only freelancers can submit applications"}, false, false)
		return
	}

	// Check if freelancer already applied
	existingApplication, err := models.GetApplicationByUserAndJob(userID, submitApplicationRequest.JobID)

	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Ctx.Output.JSON(map[string]string{"error": "Error checking existing application"}, false, false)
		return
	}

	if existingApplication != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": "You have already applied"}, false, false)
		return
	}

	applicationId, err := models.CreateApplication(user, job, submitApplicationRequest.Description)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Ctx.Output.JSON(map[string]string{"error": "Failed to save application"}, false, false)
		return
	}

	// Handle file (if provided)
	if submitApplicationRequest.FileName != "" && submitApplicationRequest.FileBase64 != "" {
		decodedFile, err := base64.StdEncoding.DecodeString(submitApplicationRequest.FileBase64)
		if err != nil {
			c.Ctx.Output.SetStatus(http.StatusBadRequest)
			c.Ctx.Output.JSON(map[string]string{"error": "Invalid base64 encoding"}, false, false)
			return
		}

		// Generate unique filename
		uniqueFileName := time.Now().Format("20060102_150405") + "_" + submitApplicationRequest.FileName
		uploadDir := "uploads"
		os.MkdirAll(uploadDir, os.ModePerm)

		// Save file
		filePath := filepath.Join(uploadDir, uniqueFileName)
		err = os.WriteFile(filePath, decodedFile, 0644)
		if err != nil {
			c.Ctx.Output.SetStatus(http.StatusInternalServerError)
			c.Ctx.Output.JSON(map[string]string{"error": "Failed to save attachment"}, false, false)
			return
		}

		err = models.CreateAttachment(applicationId, submitApplicationRequest.FileName, filePath)

		if err != nil {
			c.Ctx.Output.SetStatus(http.StatusInternalServerError)
			c.Ctx.Output.JSON(map[string]string{"error": "Failed to save attachment link in the database"}, false, false)
			return
		}
	}

	// Success response
	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Ctx.Output.JSON(map[string]string{"message": "Application submitted successfully"}, false, false)
	c.ServeJSON()
}

func (c *ApplicationController) UpdateApplication() {
	updateApplicationRequest, err := validators.UpdateApplicationValidator(c.Ctx.Input.RequestBody)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": err.Error()}, false, false)
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
		c.Ctx.Output.JSON(map[string]string{"error": "Only freelancers can update applications"}, false, false)
		return
	}

	applicationID, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": "Invalid application ID"}, false, false)
		return
	}

	application, err := models.GetApplicationByID(applicationID)
	if application == nil || err != nil {
		c.Ctx.Output.SetStatus(http.StatusNotFound)
		c.Ctx.Output.JSON(map[string]string{"error": "Application not found"}, false, false)
		return
	}

	if application.User.Id != userID {
		c.Ctx.Output.SetStatus(http.StatusForbidden)
		c.Ctx.Output.JSON(map[string]string{"error": "Freelancers are only allowed to update their own applications"}, false, false)
		return
	}

	if application.Status != "pending" {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": "Only pending applications can be updated"}, false, false)
		return
	}

	if updateApplicationRequest.Description != "" {
		application.Description = updateApplicationRequest.Description
		err = models.UpdateApplication(application)
		if err != nil {
			c.Ctx.Output.SetStatus(http.StatusInternalServerError)
			c.Ctx.Output.JSON(map[string]string{"error": "Failed to update application description"}, false, false)
			return
		}
	}

	if updateApplicationRequest.FileName != "" && updateApplicationRequest.FileBase64 != "" {
		attachment, err := models.GetAttachmentByApplicationID(application.Id)
		if err == nil && attachment != nil {

			err := os.Remove(attachment.FilePath)
			if err != nil {
				c.Ctx.Output.SetStatus(http.StatusInternalServerError)
				c.Ctx.Output.JSON(map[string]string{"error": "Failed to delete old attachment"}, false, false)
				return
			}

			err = models.DeleteAttachmentByID(attachment.Id)
			if err != nil {
				c.Ctx.Output.SetStatus(http.StatusInternalServerError)
				c.Ctx.Output.JSON(map[string]string{"error": "Failed to delete old attachment link"}, false, false)
				return
			}
		}

		decodedFile, err := base64.StdEncoding.DecodeString(updateApplicationRequest.FileBase64)
		if err != nil {
			c.Ctx.Output.SetStatus(http.StatusBadRequest)
			c.Ctx.Output.JSON(map[string]string{"error": "Invalid base64 encoding"}, false, false)
			return
		}

		uniqueFileName := time.Now().Format("20060102_150405") + "_" + updateApplicationRequest.FileName
		uploadDir := "uploads"
		os.MkdirAll(uploadDir, os.ModePerm)

		filePath := filepath.Join(uploadDir, uniqueFileName)
		err = os.WriteFile(filePath, decodedFile, 0644)
		if err != nil {
			c.Ctx.Output.SetStatus(http.StatusInternalServerError)
			c.Ctx.Output.JSON(map[string]string{"error": "Failed to save new attachment"}, false, false)
			return
		}

		err = models.CreateAttachment(application.Id, updateApplicationRequest.FileName, filePath)
		if err != nil {
			c.Ctx.Output.SetStatus(http.StatusInternalServerError)
			c.Ctx.Output.JSON(map[string]string{"error": "Failed to save new attachment link in database"}, false, false)
			return
		}

	}

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Ctx.Output.JSON(map[string]string{"message": "Application updated successfully"}, false, false)
	c.ServeJSON()

}

func (c *ApplicationController) DeleteApplication() {

	applicationID, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": "Invalid application ID"}, false, false)
		return
	}

	userID := c.Ctx.Input.GetData("id").(int)
	user, err := models.GetUserById(userID)
	if user == nil || err != nil {
		c.Ctx.Output.SetStatus(http.StatusUnauthorized)
		c.Ctx.Output.JSON(map[string]string{"error": "User not found"}, false, false)
		return
	}

	application, err := models.GetApplicationByID(applicationID)
	if err != nil || application == nil {
		c.Ctx.Output.SetStatus(http.StatusNotFound)
		c.Ctx.Output.JSON(map[string]string{"error": "Application not found"}, false, false)
		return
	}

	if application.User.Id != user.Id {
		c.Ctx.Output.SetStatus(http.StatusForbidden)
		c.Ctx.Output.JSON(map[string]string{"error": "Freelancers are only allowed to delete their own applications"}, false, false)
		return
	}

	if application.Status != "pending" {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": "Only pending applications can be deleted"}, false, false)
		return
	}

	attachment, err := models.GetAttachmentByApplicationID(application.Id)
	if err == nil && attachment != nil {
		err := os.Remove(attachment.FilePath)
		if err != nil {
			c.Ctx.Output.SetStatus(http.StatusInternalServerError)
			c.Ctx.Output.JSON(map[string]string{"error": "Failed to delete attachment"}, false, false)
			return
		}
	}

	err = models.DeleteApplicationByID(applicationID)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Ctx.Output.JSON(map[string]string{"error": "Failed to delete application"}, false, false)
		return
	}

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = map[string]string{"message": "Application deleted successfully"}
	c.ServeJSON()

}

func (c *ApplicationController) ChangeApplicationStatus() {

	changeApplicationStatusRequest, err := validators.ChangeApplicationStatusValidator(c.Ctx.Input.RequestBody)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": err.Error()}, false, false)
		return
	}

	applicationID, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": "Invalid application ID"}, false, false)
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
		c.Ctx.Output.JSON(map[string]string{"error": "Only clients can change application status"}, false, false)
		return
	}

	application, err := models.GetApplicationByID(applicationID)
	if err != nil || application == nil {
		c.Ctx.Output.SetStatus(http.StatusNotFound)
		c.Ctx.Output.JSON(map[string]string{"error": "Application not found"}, false, false)
		return
	}

	if application.Job.Client.Id != user.Id {
		c.Ctx.Output.SetStatus(http.StatusForbidden)
		c.Ctx.Output.JSON(map[string]string{"error": "You can only change application status of your jobs"}, false, false)
		return
	}

	if application.Job.Status != "open" {
		c.Ctx.Output.SetStatus(http.StatusForbidden)
		c.Ctx.Output.JSON(map[string]string{"error": "You can only change application status of open jobs"}, false, false)
		return
	}

	if application.Status != "pending" {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": "You can only change application status of pending applications"}, false, false)
		return
	}

	if changeApplicationStatusRequest.Status == "rejected" {

		if changeApplicationStatusRequest.RejectionReason == "" { // if no reason is provided, set a default one
			application.RejectionReason = "Your application was rejected."
		} else {
			application.RejectionReason = changeApplicationStatusRequest.RejectionReason
		}

		application.Status = "rejected"

		err = models.UpdateApplication(application)
		if err != nil {
			c.Ctx.Output.SetStatus(http.StatusInternalServerError)
			c.Ctx.Output.JSON(map[string]string{"error": "Failed to update application status"}, false, false)
			return
		}

	} else if changeApplicationStatusRequest.Status == "accepted" {

		// reject all other applications for this job
		applications, err := models.GetApplicationsByJobID(application.Job.Id)
		if err != nil {
			c.Ctx.Output.SetStatus(http.StatusInternalServerError)
			c.Ctx.Output.JSON(map[string]string{"error": "Failed to get applications for this job"}, false, false)
			return
		}
		for _, app := range applications {
			if app.Id != applicationID && app.Status == "pending" {

				app.Status = "rejected"
				app.RejectionReason = "Your application was automatically rejected because another application was accepted."

				err = models.UpdateApplication(&app)
				if err != nil {
					c.Ctx.Output.SetStatus(http.StatusInternalServerError)
					c.Ctx.Output.JSON(map[string]string{"error": "Failed to reject other applications"}, false, false)
					return
				}
			}
		}

		// update the status of the accepted application
		application.Status = "accepted"
		application.Job.Status = "in-progress"
		application.Job.Freelancer = application.User

		err = models.UpdateApplication(application)
		if err != nil {
			c.Ctx.Output.SetStatus(http.StatusInternalServerError)
			c.Ctx.Output.JSON(map[string]string{"error": "Failed to update application status"}, false, false)
			return
		}
		err = models.UpdateJob(application.Job)
		if err != nil {
			c.Ctx.Output.SetStatus(http.StatusInternalServerError)
			c.Ctx.Output.JSON(map[string]string{"error": "Failed to update job status"}, false, false)
			return
		}

	}

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Ctx.Output.JSON(map[string]string{"message": "Application status updated successfully"}, false, false)
	c.ServeJSON()

}
