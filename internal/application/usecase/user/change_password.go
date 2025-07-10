package usecase

import (
	"context"
	"errors"
	"todolist/internal/application/dto"
	"todolist/internal/domain/shared"
	"todolist/internal/domain/user/repository"
	vo "todolist/internal/domain/user/valueobject"
)

var (
	ErrIncorrectOldPassword = errors.New("old password is incorrect")
)

// ChangePasswordUseCase handles password changes
type ChangePasswordUseCase interface {
	Execute(ctx context.Context, userID int64, req dto.ChangePasswordRequest) error
}

type changePasswordUseCase struct {
	userRepo repository.UserRepository
}

// NewChangePasswordUseCase creates a new instance of ChangePasswordUseCase
func NewChangePasswordUseCase(userRepo repository.UserRepository) ChangePasswordUseCase {
	return &changePasswordUseCase{
		userRepo: userRepo,
	}
}

// Execute changes user password
func (uc *changePasswordUseCase) Execute(ctx context.Context, userID int64, req dto.ChangePasswordRequest) error {
	// Get the user
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return shared.ErrNotFound
	}

	// Verify old password
	if !user.Password().Matches(req.OldPassword) {
		return ErrIncorrectOldPassword
	}

	// Create new password value object
	newPassword, err := vo.NewPassword(req.NewPassword)
	if err != nil {
		return err
	}

	// Change password
	user.ChangePassword(newPassword)

	// Save updated user
	return uc.userRepo.Save(ctx, user)
}
