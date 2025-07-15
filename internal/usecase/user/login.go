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
	Execute(ctx context.Context, input dto.AuthRequest) (*dto.AuthResponse, error)
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
	tokenService service.TokenService,
	tokenIssuerName string,
) LoginUseCase {
	return &loginUseCase{
		userRepository:   userRepository,
		personRepository: personRepository,
		tokenService:     tokenService,
		tokenIssuerName:  tokenIssuerName,
	}
}

// Execute authenticates a user
func (uc *loginUseCase) Execute(ctx context.Context, input dto.AuthRequest) (*dto.AuthResponse, error) {
	// Find user by username
	user, err := uc.userRepository.FindByUsername(ctx, input.Username)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Check if user is active
	if user.Status() != vo.StatusActive {
		return nil, ErrUserNotActive
	}

	// Verify password
	if !user.Password().Matches(input.Password) {
		return nil, ErrInvalidCredentials
	}

	// Get person info
	person, err := uc.personRepository.FindByID(ctx, user.PersonID())
	if err != nil {
		return nil, err
	}

	// Generate token
	authTokens, err := uc.tokenService.GenerateTokens(ctx, uc.tokenIssuerName, user.ID())
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Return auth response
	return &dto.AuthResponse{
		Token:            authTokens.AccessToken,
		RefreshToken:     authTokens.RefreshToken,
		ExpiresAt:        authTokens.AccessMeta.ExpiresAt,
		RefreshExpiresAt: authTokens.RefreshMeta.ExpiresAt,
		User:             toUserResponseWithPerson(user, person),
	}, nil
}
