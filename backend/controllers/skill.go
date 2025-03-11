package controllers

import (
	"backend/models"
	"backend/validators"
	"net/http"

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
			c.Ctx.Output.JSON(map[string]string{"error": "This skill is already present"}, false, false)
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
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Ctx.Output.JSON(map[string]string{"error": "Failed to delete skill to freelancer data"}, false, false)
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
	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Ctx.Output.JSON(skills, false, false)
}

func (c *SkillController) AddSkillHandler() {

}

func (c *SkillController) DeleteSkillHandler() {

}

func (c *SkillController) UpdateSkillHandler() {

}
