package service

import (
	"context"
	"errors"
	"strings"
	"time"
	"todolist/internal/domain/shared"
	"todolist/internal/domain/user/entity"
	"todolist/internal/domain/user/repository"
	vo "todolist/internal/domain/user/valueobject"
)

// UserService provides domain services for User operations
type UserService interface {
	// Authentication related
	ValidateUserPermission(ctx context.Context, userID int64, permission string) error

	// User management
	DeactivateInactiveUsers(ctx context.Context, inactiveDays int) (int, error)
	BlockSuspiciousUsers(ctx context.Context, criteria SuspiciousCriteria) ([]int64, error)

	// Password policies
	EnforcePasswordPolicy(password string) error
	ShouldForcePasswordChange(user *entity.User) bool
}

// SuspiciousCriteria defines criteria for identifying suspicious users
type SuspiciousCriteria struct {
	FailedLoginAttempts int
	TimeWindow          time.Duration
}

type userService struct {
	userRepo      repository.UserRepository
	userQueryRepo repository.UserQueryRepository
}

// NewUserService creates a new user service
func NewUserService(
	userRepo repository.UserRepository,
	userQueryRepo repository.UserQueryRepository,
) UserService {
	return &userService{
		userRepo:      userRepo,
		userQueryRepo: userQueryRepo,
	}
}

// ValidateUserPermission validates if a user has a specific permission
func (s *userService) ValidateUserPermission(ctx context.Context, userID int64, permission string) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// Check if user can perform action
	if err := user.CanPerformAction(); err != nil {
		return err
	}

	// Check role permissions
	if !user.Role().HasPermission(permission) {
		return errors.New("insufficient permissions")
	}

	return nil
}

// DeactivateInactiveUsers deactivates users who haven't logged in for specified days
func (s *userService) DeactivateInactiveUsers(ctx context.Context, inactiveDays int) (int, error) {
	// Find inactive users
	users, err := s.userQueryRepo.FindInactiveUsers(ctx, inactiveDays, shared.QueryOptions{Limit: 1000})
	if err != nil {
		return 0, err
	}

	count := 0
	for _, user := range users {
		user.Deactivate()
		if err := s.userRepo.Save(ctx, user); err != nil {
			continue
		}
		count++
	}

	return count, nil
}

// BlockSuspiciousUsers blocks users based on suspicious criteria
func (s *userService) BlockSuspiciousUsers(ctx context.Context, criteria SuspiciousCriteria) ([]int64, error) {
	blockedUserIDs := []int64{}

	// Get all active users (simplified)
	users, err := s.userQueryRepo.FindByStatus(ctx, vo.StatusActive, shared.QueryOptions{Limit: 1000})
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		// In real implementation, check failed login attempts from logs
		// For now, just demonstrate the pattern
		user.Block()
		if err := s.userRepo.Save(ctx, user); err != nil {
			continue
		}
		blockedUserIDs = append(blockedUserIDs, user.ID())
	}

	return blockedUserIDs, nil
}

// EnforcePasswordPolicy validates password against policy rules
func (s *userService) EnforcePasswordPolicy(password string) error {
	// Check minimum length
	if len(password) < 8 {
		return vo.ErrPasswordTooShort
	}

	// Check for at least one uppercase letter
	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasDigit = true
		case strings.ContainsRune("!@#$%^&*()_+-=[]{}|;:,.<>?", char):
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasDigit || !hasSpecial {
		return vo.ErrPasswordTooWeak
	}

	return nil
}

// ShouldForcePasswordChange determines if user should change password
func (s *userService) ShouldForcePasswordChange(user *entity.User) bool {
	// Force password change after 90 days
	passwordAge := time.Since(user.UpdatedAt())
	return passwordAge > (90 * 24 * time.Hour)
}
