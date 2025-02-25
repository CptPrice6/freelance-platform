package middleware

import (
	"backend/utils"
	"net/http"
	"strings"

	"github.com/beego/beego/v2/server/web/context"
)

// JWTMiddleware is a middleware to protect routes with JWT
func UserMiddleware(ctx *context.Context) {
	// Get the token from the Authorization header
	tokenString := ctx.Input.Header("Authorization")
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	if tokenString == "" {
		ctx.Output.SetStatus(http.StatusUnauthorized)
		ctx.Output.Body([]byte("Authorization token missing"))
		return
	}

	// Parse and validate JWT
	claims, err := utils.ParseJWT(tokenString)
	if err != nil {
		ctx.Output.SetStatus(http.StatusUnauthorized)
		ctx.Output.Body([]byte("Invalid token"))
		return
	}
	// check if user is banned

	// Attach user email to the context for further use
	ctx.Input.SetData("email", claims.Email)
}

func AdminMiddleware(ctx *context.Context) {
	tokenString := strings.TrimPrefix(ctx.Input.Header("Authorization"), "Bearer ")

	if tokenString == "" {
		ctx.Output.SetStatus(http.StatusUnauthorized)
		ctx.Output.Body([]byte("Authorization token missing"))
		return
	}

	claims, err := utils.ParseJWT(tokenString)
	if err != nil {
		ctx.Output.SetStatus(http.StatusUnauthorized)
		ctx.Output.Body([]byte("Invalid token"))
		return
	}
	// check if user is banned

	if claims.Role != "admin" {
		ctx.Output.SetStatus(http.StatusForbidden)
		ctx.Output.Body([]byte("Access denied: Admins only"))
		return
	}

	ctx.Input.SetData("email", claims.Email)
}
