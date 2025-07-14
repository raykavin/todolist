package middleware

import (
	"strings"
	"todolist/internal/interfaces/http/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret string
}

// JWTClaims custom claims for JWT
type JWTClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// AuthMiddleware creates an authentication middleware
func AuthMiddleware(config JWTConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(
				dto.ErrorResponse("UNAUTHORIZED", "Missing authorization header", nil),
			)
		}

		// Extract token
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(
				dto.ErrorResponse("INVALID_TOKEN", "Invalid authorization format", nil),
			)
		}

		tokenString := tokenParts[1]

		// Parse and validate token
		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (any, error) {
			return []byte(config.Secret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(
				dto.ErrorResponse("INVALID_TOKEN", "Invalid or expired token", nil),
			)
		}

		// Get claims
		claims, ok := token.Claims.(*JWTClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(
				dto.ErrorResponse("INVALID_TOKEN", "Invalid token claims", nil),
			)
		}

		// Store user info in context
		c.Locals("userID", claims.UserID)
		c.Locals("username", claims.Username)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

// RequireRole creates a role-based authorization middleware
func RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole, ok := c.Locals("role").(string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(
				dto.ErrorResponse("FORBIDDEN", "Access denied", nil),
			)
		}

		// Check if user has required role
		for _, role := range roles {
			if userRole == role {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(
			dto.ErrorResponse("INSUFFICIENT_PERMISSIONS", "Insufficient permissions", nil),
		)
	}
}

// GetUserID gets user ID from context
func GetUserID(c *fiber.Ctx) string {
	userID, _ := c.Locals("userID").(string)
	return userID
}

// GetUsername gets username from context
func GetUsername(c *fiber.Ctx) string {
	username, _ := c.Locals("username").(string)
	return username
}

// GetUserRole gets user role from context
func GetUserRole(c *fiber.Ctx) string {
	role, _ := c.Locals("role").(string)
	return role
}
