package usecase

import (
	"context"
	"errors"
	"todolist/internal/domain/shared"
	"todolist/internal/domain/user/repository"
	vo "todolist/internal/domain/user/valueobject"
	"todolist/internal/dto"
)

var (
	ErrIncorrectOldPassword = errors.New("old password is incorrect")
)

// ChangePasswordUseCase handles password changes
type ChangePasswordUseCase interface {
	Execute(ctx context.Context, userID int64, input dto.ChangePasswordRequest) error
}

type changePasswordUseCase struct {
	userRepository repository.UserRepository
}

// NewChangePasswordUseCase creates a new instance of ChangePasswordUseCase
func NewChangePasswordUseCase(userRepository repository.UserRepository) ChangePasswordUseCase {
	return &changePasswordUseCase{
		userRepository: userRepository,
	}
}

// Execute changes user password
func (uc *changePasswordUseCase) Execute(ctx context.Context, userID int64, input dto.ChangePasswordRequest) error {
	// Get the user
	user, err := uc.userRepository.FindByID(ctx, userID)
	if err != nil {
		return shared.ErrNotFound
	}

	// Verify old password
	if !user.Password().Matches(input.OldPassword) {
		return ErrIncorrectOldPassword
	}

	// Create new password value object
	newPassword, err := vo.NewPassword(input.NewPassword)
	if err != nil {
		return err
	}

	// Change password
	user.ChangePassword(newPassword)

	// Save updated user
	return uc.userRepository.Save(ctx, user)
}
