package controllers

import (
	"backend/models"
	"net/http"
	"os"
	"strconv"

	"github.com/beego/beego/v2/server/web"
)

type AttachmentController struct {
	web.Controller
}

func (c *AttachmentController) DownloadAttachment() {
	attachmentID, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
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

	attachment, err := models.GetAttachmentByID(attachmentID)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusNotFound)
		c.Ctx.Output.JSON(map[string]string{"error": "Attachment not found"}, false, false)
		return
	}

	if user.Role == "freelancer" {

		if user.Id != attachment.Application.User.Id {
			c.Ctx.Output.SetStatus(http.StatusForbidden)
			c.Ctx.Output.JSON(map[string]string{"error": "You do not have permission to download this attachment"}, false, false)
			return
		}

	} else if user.Role == "client" {

		job, err := models.GetJobByID(attachment.Application.Job.Id)
		if err != nil {
			c.Ctx.Output.SetStatus(http.StatusNotFound)
			c.Ctx.Output.JSON(map[string]string{"error": "Job not found"}, false, false)
			return
		}

		if user.Id != job.Client.Id {
			c.Ctx.Output.SetStatus(http.StatusForbidden)
			c.Ctx.Output.JSON(map[string]string{"error": "You do not have permission to download this attachment"}, false, false)
			return
		}

	}

	if _, err := os.Stat(attachment.FilePath); os.IsNotExist(err) {
		c.Ctx.Output.SetStatus(http.StatusNotFound)
		c.Ctx.Output.JSON(map[string]string{"error": "File not found"}, false, false)
		return
	}

	c.Ctx.Output.Download(attachment.FilePath, attachment.FileName)
}
