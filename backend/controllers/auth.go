package controllers

import (
	"backend/models"
	"backend/utils"
	"backend/validators"
	"net/http"

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
		c.Ctx.Output.JSON(map[string]string{"error": err.Error()}, false, false)
		return
	}

	user, err := models.GetUserByEmail(registerRequest.Email)
	if user != nil || err == nil {
		c.Ctx.Output.SetStatus(http.StatusConflict)
		c.Ctx.Output.JSON(map[string]string{"error": "Email already registered"}, false, false)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Ctx.Output.JSON(map[string]string{"error": "Password hashing failed"}, false, false)
		return
	}

	err = models.CreateUser(registerRequest.Email, string(hashedPassword), "user", false)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Ctx.Output.JSON(map[string]string{"error": "Registration failed"}, false, false)
		return
	}

	// Registration successful
	c.Ctx.Output.SetStatus(http.StatusCreated)
	c.Data["json"] = map[string]string{"message": "User registered successfully"}
	c.ServeJSON()
}

// LoginHandler - Handles user login and returns a JWT token
func (c *AuthController) LoginHandler() {

	loginRequest, err := validators.LoginValidator(c.Ctx.Input.RequestBody)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": err.Error()}, false, false)
		return
	}

	user, err := models.GetUserByEmail(loginRequest.Email)
	if user == nil || err != nil {
		c.Ctx.Output.SetStatus(http.StatusUnauthorized)
		c.Ctx.Output.JSON(map[string]string{"error": "User not found"}, false, false)
		return
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusUnauthorized)
		c.Ctx.Output.JSON(map[string]string{"error": "Incorrect password"}, false, false)
		return
	}

	accessToken, refreshToken, err := utils.GenerateTokenPair(user.Id, user.Role)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Ctx.Output.JSON(map[string]string{"error": "Could not generate a new token pair"}, false, false)
		return
	}

	err = models.DeleteAllRefreshTokensForUser(user.Id)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Ctx.Output.JSON(map[string]string{"error": "Error deleting old refresh token"}, false, false)
		return
	}

	err = models.SaveRefreshToken(refreshToken, user.Id)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Ctx.Output.JSON(map[string]string{"error": "Could not store new refresh token"}, false, false)
		return
	}

	// Successful login
	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Ctx.Output.JSON(map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, false, false)
}

func (c *AuthController) LogoutHandler() {

	err := models.DeleteAllRefreshTokensForUser(c.Ctx.Input.GetData("id").(int))
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Ctx.Output.JSON(map[string]string{"error": "Error deleting old refresh token"}, false, false)
		return
	}

	c.Ctx.Output.SetStatus(http.StatusCreated)
	c.Data["json"] = map[string]string{"message": "User log out successfull"}
	c.ServeJSON()

}

// Refresh Token Handler
func (c *AuthController) RefreshTokenHandler() {

	refreshRequest, err := validators.RefreshValidator(c.Ctx.Input.RequestBody)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.JSON(map[string]string{"error": err.Error()}, false, false)
		return
	}

	// Validate Refresh Token
	claims, err := utils.ValidateRefreshToken(refreshRequest.RefreshToken)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusUnauthorized)
		c.Ctx.Output.JSON(map[string]string{"error": "Invalid refresh token"}, false, false)
		return
	}

	user, err := models.GetUserById(claims.Id)
	if user == nil || err != nil {
		c.Ctx.Output.SetStatus(http.StatusUnauthorized)
		c.Ctx.Output.JSON(map[string]string{"error": "User not found"}, false, false)
		return
	}

	// Validate the refresh token in the database
	isValid, err := models.ValidateRefreshTokenInDB(refreshRequest.RefreshToken, user.Id)
	if err != nil || !isValid {
		c.Ctx.Output.SetStatus(http.StatusUnauthorized)
		c.Ctx.Output.JSON(map[string]string{"error": "Refresh token is invalid or expired"}, false, false)
		return
	}

	err = models.DeleteAllRefreshTokensForUser(user.Id)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Ctx.Output.JSON(map[string]string{"error": "Error deleting old refresh token"}, false, false)
		return
	}

	// Generate new Access Token
	newAccessToken, newRefreshToken, err := utils.GenerateTokenPair(user.Id, user.Role)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Ctx.Output.JSON(map[string]string{"error": "Could not generate a new token pair"}, false, false)
		return
	}

	err = models.SaveRefreshToken(newRefreshToken, user.Id)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Ctx.Output.JSON(map[string]string{"error": "Could not store new refresh token"}, false, false)
		return
	}

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Ctx.Output.JSON(map[string]string{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	}, false, false)

}
