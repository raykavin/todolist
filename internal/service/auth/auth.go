package auth

import (
	"context"
	"ecommerce/internal/domain/user/entity"
	"ecommerce/internal/domain/user/repository"
	"ecommerce/internal/dto"
)

type authUsecase struct {
	userRepo repository.UserRepository
}

// NewAuthUsecase creates a new instance of authUsecase.
func NewAuthUsecase(userRepo repository.UserRepository) *authUsecase {
	return &authUsecase{
		userRepo: userRepo,
	}
}

// Register handles the user registration process.
func (uc *authUsecase) Register(ctx context.Context, req *dto.RegisterRequest) (*entity.User, error) {
	// Check if user already exists
	if _, err := uc.userRepo.GetByEmail(ctx, req.Email); err == nil {
		return nil, entity.ErrUserAlreadyExists
	}

	// Create a new user
	user, err := entity.NewUser(req.FirstName, req.LastName, req.Email, req.Password, "customer")
	if err != nil {
		return nil, err
	}

	// Save the user to the database
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUser retrieves a user by their ID.
func (uc *authUsecase) GetUser(ctx context.Context, id uint) (*entity.User, error) {
	return uc.userRepo.GetByID(ctx, id)
}
