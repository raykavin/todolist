package dto

import "time"

// AuthRequest represents the login request
type AuthRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token            string        `json:"token"`
	ExpiresAt        time.Time     `json:"expires_at"`
	RefreshToken     string        `json:"refresh_token,omitempty"`
	RefreshExpiresAt time.Time     `json:"refresh_expires_at,omitempty"`
	User             *UserResponse `json:"user"`
}
