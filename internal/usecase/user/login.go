package usecase

import (
	"context"
	"errors"
	"fmt"
	rptPerson "todolist/internal/domain/person/repository"
	rptUser "todolist/internal/domain/user/repository"
	vo "todolist/internal/domain/user/valueobject"
	"todolist/internal/dto"
	"todolist/internal/service"
)

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUserNotActive      = errors.New("user is not active")
)

// LoginUseCase handles user authentication
type LoginUseCase interface {
	Execute(ctx context.Context, req dto.AuthRequest) (*dto.AuthResponse, error)
}

type loginUseCase struct {
	userRepository   rptUser.UserRepository
	personRepository rptPerson.PersonRepository
	tokenService     service.TokenService
	tokenIssuerName  string
}

// NewLoginUseCase creates a new instance of LoginUseCase
func NewLoginUseCase(
	userRepository rptUser.UserRepository,
	personRepository rptPerson.PersonRepository,
	tokenIssuerName string,
) LoginUseCase {
	return &loginUseCase{
		userRepository:   userRepository,
		personRepository: personRepository,
	}
}

// Execute authenticates a user
func (uc *loginUseCase) Execute(ctx context.Context, req dto.AuthRequest) (*dto.AuthResponse, error) {
	// Find user by username
	user, err := uc.userRepository.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Check if user is active
	if user.Status() != vo.StatusActive {
		return nil, ErrUserNotActive
	}

	// Verify password
	if !user.Password().Matches(req.Password) {
		return nil, ErrInvalidCredentials
	}

	// Get person info
	person, err := uc.personRepository.FindByID(ctx, user.PersonID())
	if err != nil {
		return nil, err
	}

	// Generate token
	token, err := uc.tokenService.Generate(uc.tokenIssuerName, user.ID())
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Return auth response
	return &dto.AuthResponse{
		Token: token,
		User:  toUserResponseWithPerson(user, person),
	}, nil
}
