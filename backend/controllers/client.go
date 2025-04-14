package controllers

import (
	"backend/models"
	"backend/types"
	"backend/validators"
	"net/http"
	"strconv"

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
	if updateClientRequest.Description != "" {
		clientData.Description = updateClientRequest.Description
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

func (c *ClientController) GetClientsHandler() {
	users, err := models.GetUsersByRole("client")
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Ctx.Output.JSON(map[string]string{"error": "Error fetching clients"}, false, false)
		return
	}

	var clients []types.ClientInfo

	for _, user := range users {
		clientData, err := models.GetClientDataByUserID(user.Id)
		if err != nil || clientData == nil {
			continue
		}

		clientInfo := types.ClientInfo{
			ID:          user.Id,
			Name:        user.Name,
			Surname:     user.Surname,
			CompanyName: clientData.CompanyName,
		}

		clients = append(clients, clientInfo)
	}

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = clients
	c.ServeJSON()
}

func (c *ClientController) GetClientHandler() {
	idStr := c.Ctx.Input.Param(":id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": "Invalid client ID"}, false, false)
		return
	}

	user, err := models.GetUserById(userID)
	if user == nil || err != nil {
		c.Ctx.Output.SetStatus(http.StatusNotFound)
		c.Ctx.Output.JSON(map[string]string{"error": "User not found"}, false, false)
		return
	}
	if user.Role != "client" {
		c.Ctx.Output.SetStatus(http.StatusNotFound)
		c.Ctx.Output.JSON(map[string]string{"error": "Client not found"}, false, false)
		return
	}

	response := types.UserResponse{
		ID:      user.Id,
		Email:   user.Email,
		Role:    user.Role,
		Name:    user.Name,
		Surname: user.Surname,
	}

	clientData, err := models.GetClientDataByUserID(user.Id)
	if err == nil && clientData != nil {

		response.ClientData = &types.ClientData{
			CompanyName: clientData.CompanyName,
			Industry:    clientData.Industry,
			Location:    clientData.Location,
			Description: clientData.Description,
		}
	}

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = response
	c.ServeJSON()

}
