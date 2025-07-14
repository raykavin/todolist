package usecase

import (
	"context"
	"time"
	entPerson "todolist/internal/domain/person/entity"
	repoPerson "todolist/internal/domain/person/repository"
	entUser "todolist/internal/domain/user/entity"
	repoUser "todolist/internal/domain/user/repository"
	vo "todolist/internal/domain/user/valueobject"
	"todolist/internal/dto"
)

// CreateUserUseCase handles user creation
type CreateUserUseCase interface {
	Execute(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error)
}

type createUserUseCase struct {
	userRepository   repoUser.UserRepository
	personRepository repoPerson.PersonRepository
}

// NewCreateUserUseCase creates a new instance of CreateUserUseCase
func NewCreateUserUseCase(
	userRepository repoUser.UserRepository,
	personRepository repoPerson.PersonRepository,
) CreateUserUseCase {
	return &createUserUseCase{
		userRepository:   userRepository,
		personRepository: personRepository,
	}
}

// Execute creates a new user
func (uc *createUserUseCase) Execute(ctx context.Context, input dto.CreateUserRequest) (*dto.UserResponse, error) {
	// Verify person exists
	person, err := uc.personRepository.FindByID(ctx, input.PersonID)
	if err != nil {
		return nil, entPerson.ErrPersonNotFound
	}

	// Check if username already exists
	exists, err := uc.userRepository.ExistsByUsername(ctx, input.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, entUser.ErrUsernameExists
	}

	// Check if person already has a user
	exists, err = uc.userRepository.ExistsByPersonID(ctx, input.PersonID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, entUser.ErrUserAlreadyExists
	}

	// Create password value object
	password, err := vo.NewPassword(input.Password)
	if err != nil {
		return nil, err
	}

	// Create user entity
	user, err := entUser.NewUser(
		time.Now().Unix(),
		input.PersonID,
		input.Username,
		password,
		vo.RoleUser,
	)
	if err != nil {
		return nil, err
	}

	// Save user
	if err := uc.userRepository.Save(ctx, user); err != nil {
		return nil, err
	}

	// Convert to response with person info
	return toUserResponseWithPerson(user, person), nil
}

// Helper function to convert entity to DTO with person info
func toUserResponseWithPerson(user *entUser.User, person *entPerson.Person) *dto.UserResponse {
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
