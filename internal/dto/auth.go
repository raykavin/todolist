package dto

// AuthRequest represents the login request
type AuthRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token        string        `json:"token"`
	ExpiresAt    string        `json:"expires_at"`
	RefreshToken string        `json:"refresh_token"`
	User         *UserResponse `json:"user"`
}
