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
