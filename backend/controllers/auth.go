package controllers

import (
	"backend/models"
	"backend/utils"
	"backend/validators"
	"encoding/json"
	"net/http"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	web.Controller
}

// RegisterHandler - Handles user registration
func (c *AuthController) RegisterHandler() {

	registerRequest, err := validators.RegisterValidator(c.Ctx.Input.RequestBody)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": err.Error()}
		c.ServeJSON()
		return
	}

	alreadyExists := models.UserAlreadyExists(registerRequest.Email)
	if alreadyExists {
		c.Ctx.Output.SetStatus(http.StatusConflict)
		c.Data["json"] = map[string]string{"error": "Email already registered"}
		c.ServeJSON()
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Password hashing failed"}
		c.ServeJSON()
		return
	}

	err = models.CreateUser(registerRequest.Email, string(hashedPassword), "user")
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Registration failed"}
		c.ServeJSON()
		return
	}

	// Registration successful
	c.Ctx.Output.SetStatus(http.StatusCreated)
	c.Data["json"] = map[string]string{"message": "User registered successfully"}
	c.ServeJSON()
}

// LoginHandler - Handles user login and returns a JWT token
func (c *AuthController) LoginHandler() {

	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &loginData)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Invalid input"}
		c.ServeJSON()
		return
	}

	missingFields := []string{}
	if loginData.Email == "" {
		missingFields = append(missingFields, "email")
	}
	if loginData.Password == "" {
		missingFields = append(missingFields, "password")
	}

	if len(missingFields) > 0 {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{
			"error":          "Missing required fields",
			"missing_fields": missingFields,
		}
		c.ServeJSON()
		return
	}

	o := orm.NewOrm()
	user := models.User{Email: loginData.Email}
	err = o.Read(&user, "Email")
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusUnauthorized)
		c.Data["json"] = map[string]string{"error": "User not found"}
		c.ServeJSON()
		return
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password))
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusUnauthorized)
		c.Data["json"] = map[string]string{"error": "Incorrect password"}
		c.ServeJSON()
		return
	}

	// Generate JWT token
	accessToken, err := utils.GenerateAccessToken(user.Email, user.Role)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Access token generation failed"}
		c.ServeJSON()
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(user.Email, user.Role)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Refresh token generation failed"}
		c.ServeJSON()
		return
	}

	// Successful login
	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}
	c.ServeJSON()
}

// Refresh Token Handler
func (c *AuthController) RefreshTokenHandler() {
	var requestData struct {
		RefreshToken string `json:"refresh_token"`
	}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &requestData)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.ContentType("application/json")
		c.Data["json"] = map[string]string{"error": "Invalid input"}
		c.ServeJSON()
		return
	}
	// Validate Refresh Token
	claims, err := utils.ValidateRefreshToken(requestData.RefreshToken)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusUnauthorized)
		c.Data["json"] = map[string]string{"error": "Invalid refresh token"}
		c.ServeJSON()
		return
	}

	// Generate new Access Token
	newAccessToken, err := utils.GenerateAccessToken(claims.Email, claims.Role)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Could not generate a new access token"}
		c.ServeJSON()
		return
	}
	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = map[string]string{
		"access_token": newAccessToken,
	}
	c.ServeJSON()
}
