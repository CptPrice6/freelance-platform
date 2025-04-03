package controllers

import (
	"backend/models"
	"backend/validators"
	"encoding/base64"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/beego/beego/v2/server/web"
)

type ApplicationController struct {
	web.Controller
}

func (c *ApplicationController) GetFreelancerApplications() {}

func (c *ApplicationController) GetFreelancerApplication() {}

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

	// update only pending applications!
}

func (c *ApplicationController) DeleteApplication() {

	// delete only pending applications!
}

func (c *ApplicationController) ChangeApplicationStatus() {

	// only CLIENT can change and only STATUS

	// if applications status is accepted or rejected it cannot be changed

	// if rejected it can be only changed to accepted or rejected

	// when rejected -> just change status
	// when accepeted -> change status and reject ALL other applications
}
