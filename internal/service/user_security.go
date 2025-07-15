package service

import (
	"context"
	"errors"
	"time"

	"todolist/internal/domain/shared"
	"todolist/internal/domain/user/entity"
	"todolist/internal/domain/user/repository"
	uservo "todolist/internal/domain/user/valueobject"
)

// UserSecurityService defines operations for managing user security and governance policies.
//
// This service provides methods to validate user permissions, handle inactive or suspicious users,
// and enforce password policies.
type UserSecurityService interface {
	// ValidateUserPermission verifies if a user has the specified permission.
	//
	// Returns an error if the user does not have the required permission or if any validation fails.
	ValidateUserPermission(ctx context.Context, userID int64, permission string) error

	// DeactivateInactiveUsers deactivates users who have been inactive for more than the specified number of days.
	//
	// The 'limit' parameter restricts the maximum number of users to process in a single call.
	// Returns the number of users deactivated and an error if the operation fails.
	DeactivateInactiveUsers(ctx context.Context, inactiveDays int, limit int) (int, error)

	// BlockSuspiciousUsers blocks users who meet the given suspicious activity criteria.
	//
	// The 'criteria' parameter defines the conditions that qualify users as suspicious.
	// The 'limit' parameter restricts the maximum number of users to block.
	// Returns a slice of blocked user IDs and an error if the operation fails.
	BlockSuspiciousUsers(ctx context.Context, criteria SuspiciousCriteria, limit int) ([]int64, error)

	// EnforcePasswordPolicy checks if the given password complies with the security policy.
	//
	// Returns an error if the password does not meet the policy requirements.
	EnforcePasswordPolicy(password string) error

	// ShouldForcePasswordChange determines whether the given user must change their password.
	//
	// Returns true if a password change should be enforced based on the current security policy.
	ShouldForcePasswordChange(user *entity.User) bool
}

// SuspiciousCriteria defines the criteria for blocking suspicious users.
type SuspiciousCriteria struct {
	FailedLoginAttempts int
	TimeWindow          time.Duration
}

type userSecurityService struct {
	userRepository      repository.UserRepository
	userQueryRepository repository.UserQueryRepository
}

// NewUserSecurityService creates a new instance.
func NewUserSecurityService(
	userRepository repository.UserRepository,
	userQueryRepository repository.UserQueryRepository,
) UserSecurityService {
	return &userSecurityService{
		userRepository:      userRepository,
		userQueryRepository: userQueryRepository,
	}
}

func (s *userSecurityService) ValidateUserPermission(
	ctx context.Context,
	userID int64,
	permission string,
) error {
	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	if err := user.CanPerformAction(); err != nil {
		return err
	}

	if !user.Role().HasPermission(permission) {
		return errors.New("insufficient permissions")
	}

	return nil
}

func (s *userSecurityService) DeactivateInactiveUsers(
	ctx context.Context,
	inactiveDays,
	limit int,
) (int, error) {
	users, err := s.userQueryRepository.FindInactiveUsers(
		ctx,
		inactiveDays,
		shared.QueryOptions{Limit: limit},
	)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, user := range users {
		user.Deactivate()
		if err := s.userRepository.Save(ctx, user); err != nil {
			continue
		}
		count++
	}

	return count, nil
}

func (s *userSecurityService) BlockSuspiciousUsers(
	ctx context.Context,
	criteria SuspiciousCriteria,
	limit int,
) ([]int64, error) {
	blockedUserIDs := []int64{}

	users, err := s.userQueryRepository.FindByStatus(
		ctx,
		uservo.StatusActive,
		shared.QueryOptions{Limit: limit},
	)
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if user.FailedLoginAttempts() >= criteria.FailedLoginAttempts {
			if time.Since(user.LastLoginAttemptAt()) <= criteria.TimeWindow {
				user.Block()
				if err := s.userRepository.Save(ctx, user); err != nil {
					continue
				}
				blockedUserIDs = append(blockedUserIDs, user.ID())
			}
		}
	}

	return blockedUserIDs, nil
}

func (s *userSecurityService) EnforcePasswordPolicy(password string) error {
	// Delegate the password validation for domain value object
	return uservo.ValidatePassword(password)
}

func (s *userSecurityService) ShouldForcePasswordChange(user *entity.User) bool {
	passwordAge := time.Since(user.UpdatedAt())
	return passwordAge > (90 * 24 * time.Hour)
}
