package dto

// UserResponse represents a user in API responses
type UserResponse struct {
	ID        int64       `json:"id"`
	PersonID  int64       `json:"person_id"`
	Username  string      `json:"username"`
	Status    string      `json:"status"`
	Role      string      `json:"role"`
	Person    *PersonInfo `json:"person,omitempty"`
	CreatedAt string      `json:"created_at"`
	UpdatedAt string      `json:"updated_at"`
}

// PersonInfo represents basic person information
type PersonInfo struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	Token string        `json:"token"`
	User  *UserResponse `json:"user"`
}
