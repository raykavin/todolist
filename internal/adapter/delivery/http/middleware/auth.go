package middleware

import (
	"net/http"
	"strings"
	"todolist/internal/dto"
	"todolist/internal/service"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware creates an authentication middleware
func AuthMiddleware(tokenService service.TokenService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse("UNAUTHORIZED", "Missing authorization header", nil))
			ctx.Abort()
			return
		}

		// Extract Bearer token
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse("INVALID_TOKEN", "Invalid authorization format", nil))
			ctx.Abort()
			return
		}

		tokenString := tokenParts[1]

		// Validate token
		validationResult, err := tokenService.ValidateToken(ctx, tokenString)
		if err != nil || validationResult.UserID == 0 {
			ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse("INVALID_TOKEN", "Invalid or expired token", nil))
			ctx.Abort()
			return
		}

		// Extract custom claims safely
		username := ""
		role := ""

		if val, ok := validationResult.Claims["username"]; ok {
			if strVal, ok := val.(string); ok {
				username = strVal
			}
		}

		if val, ok := validationResult.Claims["role"]; ok {
			if strVal, ok := val.(string); ok {
				role = strVal
			}
		}

		// Store user info in Gin context
		ctx.Set("userID", validationResult.UserID)
		ctx.Set("username", username)
		ctx.Set("role", role)

		ctx.Next()
	}
}

// RequireRole creates a role-based authorization middleware
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRole, exists := ctx.Get("role")
		if !exists {
			ctx.JSON(http.StatusForbidden, dto.ErrorResponse("FORBIDDEN", "Access denied", nil))
			ctx.Abort()
			return
		}

		userRoleStr, ok := userRole.(string)
		if !ok {
			ctx.JSON(http.StatusForbidden, dto.ErrorResponse("FORBIDDEN", "Invalid role", nil))
			ctx.Abort()
			return
		}

		for _, role := range roles {
			if userRoleStr == role {
				ctx.Next()
				return
			}
		}

		ctx.JSON(http.StatusForbidden, dto.ErrorResponse("INSUFFICIENT_PERMISSIONS", "Insufficient permissions", nil))
		ctx.Abort()
	}
}
