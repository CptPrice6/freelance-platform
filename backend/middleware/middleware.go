package middleware

import (
	"backend/utils"
	"net/http"
	"strings"

	"github.com/beego/beego/v2/server/web/context"
)

// A middleware to protect routes with JWT authentication
func UserAuthMiddleware(ctx *context.Context) {
	// Get the token from the Authorization header
	tokenString := ctx.Input.Header("Authorization")
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	if tokenString == "" {
		ctx.Output.SetStatus(http.StatusUnauthorized)
		ctx.Output.JSON(map[string]string{"error": "Authorization token missing"}, false, false)
		return
	}

	// Parse and validate JWT
	claims, err := utils.ValidateAccessToken(tokenString)
	if err != nil {
		ctx.Output.SetStatus(http.StatusUnauthorized)
		ctx.Output.JSON(map[string]string{"error": "Invalid token"}, false, false)
		return
	}

	// check if user is banned

	// Attach user id to the context for further use
	ctx.Input.SetData("id", claims.Id)
}

// A middleware to protect routes with JWT authentication
// This middleware is used to protect routes that require admin access
func AdminAuthMiddleware(ctx *context.Context) {
	tokenString := strings.TrimPrefix(ctx.Input.Header("Authorization"), "Bearer ")

	if tokenString == "" {
		ctx.Output.SetStatus(http.StatusUnauthorized)
		ctx.Output.JSON(map[string]string{"error": "Authorization token missing"}, false, false)
		return
	}

	claims, err := utils.ValidateAccessToken(tokenString)
	if err != nil {
		ctx.Output.SetStatus(http.StatusUnauthorized)
		ctx.Output.JSON(map[string]string{"error": "Invalid token"}, false, false)
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
