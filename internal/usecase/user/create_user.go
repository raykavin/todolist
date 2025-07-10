package usecase

import (
	"context"
	"errors"
	"time"
	"todolist/internal/application/dto"
	personEnt "todolist/internal/domain/person/entity"
	"todolist/internal/domain/person/repository"
	"todolist/internal/domain/user/entity"
	userRepo "todolist/internal/domain/user/repository"
	vo "todolist/internal/domain/user/valueobject"
)

var (
	ErrPersonNotFound    = errors.New("person not found")
	ErrUsernameExists    = errors.New("username already exists")
	ErrUserAlreadyExists = errors.New("user already exists for this person")
)

// CreateUserUseCase handles user creation
type CreateUserUseCase interface {
	Execute(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error)
}

type createUserUseCase struct {
	userRepo   userRepo.UserRepository
	personRepo repository.PersonRepository
}

// NewCreateUserUseCase creates a new instance of CreateUserUseCase
func NewCreateUserUseCase(
	userRepo userRepo.UserRepository,
	personRepo repository.PersonRepository,
) CreateUserUseCase {
	return &createUserUseCase{
		userRepo:   userRepo,
		personRepo: personRepo,
	}
}

// Execute creates a new user
func (uc *createUserUseCase) Execute(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	// Verify person exists
	person, err := uc.personRepo.FindByID(ctx, req.PersonID)
	if err != nil {
		return nil, ErrPersonNotFound
	}

	// Check if username already exists
	exists, err := uc.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUsernameExists
	}

	// Check if person already has a user
	exists, err = uc.userRepo.ExistsByPersonID(ctx, req.PersonID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUserAlreadyExists
	}

	// Create password value object
	password, err := vo.NewPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create role value object
	// role := vo.UserRole(req.Role)
	// if !role.IsValid() {
	// 	return nil, errors.New("invalid role")
	// }

	// Create user entity
	user, err := entity.NewUser(
		time.Now().Unix(),
		req.PersonID,
		req.Username,
		password,
		// role, Automatically set by entity
	)
	if err != nil {
		return nil, err
	}

	// Save user
	if err := uc.userRepo.Save(ctx, user); err != nil {
		return nil, err
	}

	// Convert to response with person info
	return toUserResponseWithPerson(user, person), nil
}

// Helper function to convert entity to DTO with person info
func toUserResponseWithPerson(user *entity.User, person *personEnt.Person) *dto.UserResponse {
	return &dto.UserResponse{
		ID:       user.ID(),
		PersonID: user.PersonID(),
		Username: user.Username(),
		Status:   user.Status().String(),
		Role:     user.Role().String(),
		Person: &dto.PersonInfo{
			ID:    person.ID(),
			Name:  person.Name(),
			Email: person.Email().Value(),
		},
		CreatedAt: user.CreatedAt().Format("2006-01-02T15:04:05Z"),
		UpdatedAt: user.UpdatedAt().Format("2006-01-02T15:04:05Z"),
	}
}
