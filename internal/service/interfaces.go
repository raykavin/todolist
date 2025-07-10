package service

import (
	"context"
	"todolist/internal/domain/user/entity"
)

// TokenService defines an interface for token generation and validation.
type TokenService interface {
	Generate(issuerName string, userID int64) (string, error)
	Validate(token string) (int64, error)
}

// UserSecurityService provides user security and governance services.
type UserSecurityService interface {
	ValidateUserPermission(ctx context.Context, userID int64, permission string) error
	DeactivateInactiveUsers(ctx context.Context, inactiveDays int, limit int) (int, error)
	BlockSuspiciousUsers(ctx context.Context, criteria SuspiciousCriteria, limit int) ([]int64, error)

	EnforcePasswordPolicy(password string) error
	ShouldForcePasswordChange(user *entity.User) bool
}
