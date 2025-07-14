package usecase

import (
	"context"
	"time"
	"todolist/internal/domain/person/entity"
	"todolist/internal/domain/person/repository"
	vo "todolist/internal/domain/person/valueobject"
	sharedvo "todolist/internal/domain/shared/valueobject"
	"todolist/internal/dto"
)

// CreatePersonUseCase handles person creation
type CreatePersonUseCase interface {
	Execute(ctx context.Context, input dto.CreatePersonRequest) (*dto.PersonResponse, error)
}

type createPersonUseCase struct {
	personRepository repository.PersonRepository
}

// NewCreatePersonUseCase creates a new instance of CreatePersonUseCase
func NewCreatePersonUseCase(personRepository repository.PersonRepository) CreatePersonUseCase {
	return &createPersonUseCase{
		personRepository: personRepository,
	}
}

// Execute creates a new person
func (uc *createPersonUseCase) Execute(ctx context.Context, input dto.CreatePersonRequest) (*dto.PersonResponse, error) {
	// Check if email already exists
	exists, err := uc.personRepository.ExistsByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, vo.ErrInvalidEmail
	}

	// Validate Tax ID
	taxID, err := vo.NewTaxID(input.TaxID)
	if err != nil {
		return nil, err
	}

	// Create email value object
	email, err := vo.NewEmail(input.Email)
	if err != nil {
		return nil, err
	}

	birthDate := (*sharedvo.Date)(nil)

	// If birth date is provided, create date value object
	if input.BirthDate != "" {
		birthDt, err := sharedvo.NewDate(input.BirthDate)
		if err != nil {
			return nil, err
		}

		birthDate = &birthDt
	}

	// Create person entity
	person, err := entity.NewPerson(
		time.Now().Unix(),
		input.Name,
		input.Phone,
		taxID,
		email,
		birthDate,
	)
	if err != nil {
		return nil, err
	}

	// Save person
	if err := uc.personRepository.Save(ctx, person); err != nil {
		return nil, err
	}

	// Convert to response
	return toPersonResponse(person), nil
}

// Helper function to convert entity to DTO
func toPersonResponse(person *entity.Person) *dto.PersonResponse {
	return &dto.PersonResponse{
		ID:        person.ID(),
		Name:      person.Name(),
		Email:     person.Email().Value(),
		Phone:     person.Phone(),
		CreatedAt: person.CreatedAt().Format("2006-01-02T15:04:05Z"),
		UpdatedAt: person.UpdatedAt().Format("2006-01-02T15:04:05Z"),
	}
}
