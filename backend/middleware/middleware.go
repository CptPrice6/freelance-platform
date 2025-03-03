package middleware

import (
	"backend/models"
	"backend/utils"
	"net/http"
	"strings"

	"github.com/beego/beego/v2/server/web/context"
)

// A middleware to protect routes with JWT authentication
func UserAuthMiddleware(ctx *context.Context) {

	if ctx.Request.Method == "OPTIONS" {
		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.Output.SetStatus(http.StatusOK)
		return
	}

	ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	ctx.ResponseWriter.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	ctx.ResponseWriter.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
	ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")

	// Get the token from the Authorization header
	tokenString := ctx.Input.Header("Authorization")
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	if tokenString == "" {
		ctx.Output.SetStatus(http.StatusUnauthorized)
		ctx.Output.JSON(map[string]string{"error": "Missing access token"}, false, false)
		return
	}

	// Parse and validate JWT
	claims, err := utils.ValidateAccessToken(tokenString)
	if err != nil {
		ctx.Output.SetStatus(http.StatusUnauthorized)
		ctx.Output.JSON(map[string]string{"error": "Invalid access token"}, false, false)
		return
	}

	banStatus, err := models.IsUserBanned(claims.Id)
	if err != nil {
		ctx.Output.SetStatus(http.StatusInternalServerError)
		ctx.Output.JSON(map[string]string{"error": "Error fetching user data"}, false, false)
		return
	}

	if banStatus {
		ctx.Output.SetStatus(http.StatusForbidden)
		ctx.Output.JSON(map[string]string{"error": "User is banned"}, false, false)
		return
	}

	// Attach user id to the context for further use
	ctx.Input.SetData("id", claims.Id)
}

// A middleware to protect routes with JWT authentication
// This middleware is used to protect routes that require admin access
func AdminAuthMiddleware(ctx *context.Context) {

	if ctx.Request.Method == "OPTIONS" {
		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.Output.SetStatus(http.StatusOK)
		return
	}

	ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	ctx.ResponseWriter.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	ctx.ResponseWriter.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
	ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")

	tokenString := strings.TrimPrefix(ctx.Input.Header("Authorization"), "Bearer ")

	if tokenString == "" {
		ctx.Output.SetStatus(http.StatusUnauthorized)
		ctx.Output.JSON(map[string]string{"error": "Missing access token"}, false, false)
		return
	}

	claims, err := utils.ValidateAccessToken(tokenString)
	if err != nil {
		ctx.Output.SetStatus(http.StatusUnauthorized)
		ctx.Output.JSON(map[string]string{"error": "Invalid access token"}, false, false)
		return
	}

	if claims.Role != "admin" {
		ctx.Output.SetStatus(http.StatusForbidden)
		ctx.Output.JSON(map[string]string{"error": "Access denied: Admins only"}, false, false)
		return
	}

	// Attach user id to the context for further use
	ctx.Input.SetData("id", claims.Id)
}
