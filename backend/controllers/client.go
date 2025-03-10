package controllers

import (
	"backend/models"
	"backend/validators"
	"net/http"

	"github.com/beego/beego/v2/server/web"
)

type ClientController struct {
	web.Controller
}

func (c *ClientController) UpdateClientDataHandler() {
	userID := c.Ctx.Input.GetData("id").(int)

	clientData, err := models.GetClientDataByUserID(userID)
	if err != nil || clientData == nil {
		c.Ctx.Output.SetStatus(http.StatusUnauthorized)
		c.Ctx.Output.JSON(map[string]string{"error": "Client profile not found"}, false, false)
		return
	}

	updateClientRequest, err := validators.UpdateClientDataValidator(c.Ctx.Input.RequestBody)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": err.Error()}, false, false)
		return
	}

	if updateClientRequest.CompanyName != "" {
		clientData.CompanyName = updateClientRequest.CompanyName
	}
	if updateClientRequest.Industry != "" {
		clientData.Industry = updateClientRequest.Industry
	}
	if updateClientRequest.Location != "" {
		clientData.Location = updateClientRequest.Location
	}

	err = models.UpdateClientData(clientData)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Ctx.Output.JSON(map[string]string{"error": "Error updating client data"}, false, false)
		return
	}

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = map[string]string{"message": "Client data updated successfully"}
	c.ServeJSON()
}
