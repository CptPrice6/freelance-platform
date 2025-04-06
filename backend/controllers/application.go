package controllers

import (
	"backend/models"
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

	// only CLIENT can change and only STATUS + rejection reason

	// CHECK application.Job.Client.Id == user.Id
	// CHECK application.Job.Status == "open"

	// if applications status is accepted or rejected it cannot be changed

	// if rejected it can be only changed to accepted or rejected

	// when rejected -> change status + add rejection reason ( if present )
	// when accepeted ->
	// change status, reject ALL other applications and add some basic same rejection reason
	// also add that freelancer to the job and change status to "in progress"

}
