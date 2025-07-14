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
