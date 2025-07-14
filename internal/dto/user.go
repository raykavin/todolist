package dto

// CreateUserRequest represents the request to create a user
type CreateUserRequest struct {
	PersonID int64  `json:"person_id" validate:"required"`
	Username string `json:"username"  validate:"required,min=3,max=50"`
	Password string `json:"password"  validate:"required,min=8"`
}

// ChangePasswordRequest represents the password change request
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

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
