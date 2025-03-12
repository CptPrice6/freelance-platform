package controllers

import (
	"backend/models"
	"backend/types"
	"backend/validators"
	"net/http"
	"strconv"
	"strings"

	"github.com/beego/beego/v2/server/web"
)

type SkillController struct {
	web.Controller
}

func (c *SkillController) AddFreelancerSkillHandler() {
	userID := c.Ctx.Input.GetData("id").(int)
	user, err := models.GetUserById(userID)
	if user == nil || err != nil {
		c.Ctx.Output.SetStatus(http.StatusUnauthorized)
		c.Ctx.Output.JSON(map[string]string{"error": "User not found"}, false, false)
		return
	}

	addSkillRequest, err := validators.AddDeleteSkillValidator(c.Ctx.Input.RequestBody)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": err.Error()}, false, false)
		return
	}

	skill, err := models.GetSkillById(addSkillRequest.SkillID)
	if skill == nil || err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": "Skill not found"}, false, false)
		return
	}

	freelancerData, err := models.GetFreelancerDataByUserID(userID)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": "Freelancer data not found"}, false, false)
		return
	}

	for _, existingSkill := range freelancerData.Skills {
		if existingSkill.Id == skill.Id {
			c.Ctx.Output.SetStatus(http.StatusBadRequest)
			c.Ctx.Output.JSON(map[string]string{"error": "Freelancer already has this skill"}, false, false)
			return
		}
	}

	err = models.AddSkillToFreelancerData(freelancerData, skill)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Ctx.Output.JSON(map[string]string{"error": err.Error()}, false, false)
		return
	}

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Ctx.Output.JSON(map[string]string{"message": "Skill added successfully"}, false, false)
}

func (c *SkillController) DeleteFreelancerSkillHandler() {

	userID := c.Ctx.Input.GetData("id").(int)
	user, err := models.GetUserById(userID)
	if user == nil || err != nil {
		c.Ctx.Output.SetStatus(http.StatusUnauthorized)
		c.Ctx.Output.JSON(map[string]string{"error": "User not found"}, false, false)
		return
	}

	addSkillRequest, err := validators.AddDeleteSkillValidator(c.Ctx.Input.RequestBody)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": err.Error()}, false, false)
		return
	}

	skill, err := models.GetSkillById(addSkillRequest.SkillID)
	if skill == nil || err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": "Skill not found"}, false, false)
		return
	}

	freelancerData, err := models.GetFreelancerDataByUserID(userID)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": "Freelancer data not found"}, false, false)
		return
	}

	err = models.DeleteSkillFromFreelancerData(freelancerData, skill)
	if err != nil {
		if err.Error() == "skill not found" {
			c.Ctx.Output.SetStatus(http.StatusBadRequest)
			c.Ctx.Output.JSON(map[string]string{"error": "Freelancer does not have this skill"}, false, false)
		} else {
			c.Ctx.Output.SetStatus(http.StatusInternalServerError)
			c.Ctx.Output.JSON(map[string]string{"error": "Failed to delete skill from freelancer data"}, false, false)
		}
		return
	}
	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Ctx.Output.JSON(map[string]string{"message": "Skill deleted successfully"}, false, false)

}

func (c *SkillController) GetSkillsHandler() {

	skills, err := models.GetAllSkills()
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Ctx.Output.JSON(map[string]string{"error": "Failed to fetch skills"}, false, false)
		return
	}

	var skillsResponse []types.Skill

	for _, skill := range skills {
		skillsResponse = append(skillsResponse, types.Skill{
			Id:   skill.Id,
			Name: skill.Name,
		})
	}

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Ctx.Output.JSON(skillsResponse, false, false)
}

func (c *SkillController) AddSkillHandler() {

	updateSkillRequest, err := validators.AddUpdateSkillValidator(c.Ctx.Input.RequestBody)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": err.Error()}, false, false)
		return
	}

	err = models.CreateSkill(updateSkillRequest.SkillName)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			c.Ctx.Output.SetStatus(http.StatusBadRequest)
			c.Ctx.Output.JSON(map[string]string{"error": "Skill already exists!"}, false, false)
		} else {
			c.Ctx.Output.SetStatus(http.StatusInternalServerError)
			c.Ctx.Output.JSON(map[string]string{"error": err.Error()}, false, false)
		}
		return
	}

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = map[string]string{"message": "Skill added successfully"}
	c.ServeJSON()

}

func (c *SkillController) DeleteSkillHandler() {
	idStr := c.Ctx.Input.Param(":id")
	skillID, err := strconv.Atoi(idStr)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": "Invalid skill ID"}, false, false)
		return
	}

	skill, err := models.GetSkillById(skillID)
	if skill == nil || err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": "Skill not found"}, false, false)
		return
	}

	err = models.DeleteSkillByID(skillID)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": "Skill not found"}, false, false)
		return
	}

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = map[string]string{"message": "Skill deleted successfully"}
	c.ServeJSON()

}

func (c *SkillController) UpdateSkillHandler() {
	idStr := c.Ctx.Input.Param(":id")
	skillID, err := strconv.Atoi(idStr)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": "Invalid skill ID"}, false, false)
		return
	}

	updateSkillRequest, err := validators.AddUpdateSkillValidator(c.Ctx.Input.RequestBody)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": err.Error()}, false, false)
		return
	}

	skill, err := models.GetSkillById(skillID)
	if skill == nil || err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": "Skill not found"}, false, false)
		return
	}

	if updateSkillRequest.SkillName != "" {
		skill.Name = updateSkillRequest.SkillName
	}

	err = models.UpdateSkill(skill)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			c.Ctx.Output.SetStatus(http.StatusBadRequest)
			c.Ctx.Output.JSON(map[string]string{"error": "Skill already exists!"}, false, false)
		} else {
			c.Ctx.Output.SetStatus(http.StatusInternalServerError)
			c.Ctx.Output.JSON(map[string]string{"error": err.Error()}, false, false)
		}
		return
	}

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = map[string]string{"message": "Skill name updated successfully"}
	c.ServeJSON()

}
