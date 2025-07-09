package handler

import (
	"ecommerce/internal/domain/user/usecase"
	"ecommerce/internal/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication-related HTTP requests.
type AuthHandler struct {
	authUsecase usecase.UserUsecase
}

// NewAuthHandler creates a new instance of AuthHandler.
func NewAuthHandler(authUsecase usecase.UserUsecase) *AuthHandler {
	return &AuthHandler{authUsecase}
}

// Register is the handler for user registration.
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authUsecase.Register(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}
