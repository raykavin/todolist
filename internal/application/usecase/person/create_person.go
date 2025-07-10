package usecase

import (
	"context"
	"time"
	"todolist/internal/application/dto"
	"todolist/internal/domain/person/entity"
	"todolist/internal/domain/person/repository"
	vo "todolist/internal/domain/person/valueobject"
	sharedvo "todolist/internal/domain/shared/valueobject"
)

// CreatePersonUseCase handles person creation
type CreatePersonUseCase interface {
	Execute(ctx context.Context, req dto.CreatePersonRequest) (*dto.PersonResponse, error)
}

type createPersonUseCase struct {
	personRepo repository.PersonRepository
}

// NewCreatePersonUseCase creates a new instance of CreatePersonUseCase
func NewCreatePersonUseCase(personRepo repository.PersonRepository) CreatePersonUseCase {
	return &createPersonUseCase{
		personRepo: personRepo,
	}
}

// Execute creates a new person
func (uc *createPersonUseCase) Execute(ctx context.Context, req dto.CreatePersonRequest) (*dto.PersonResponse, error) {
	// Check if email already exists
	exists, err := uc.personRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, vo.ErrInvalidEmail
	}

	// Validate Tax ID
	taxID, err := vo.NewTaxID(req.TaxID)
	if err != nil {
		return nil, err
	}

	// Create email value object
	email, err := vo.NewEmail(req.Email)
	if err != nil {
		return nil, err
	}

	birthDate := (*sharedvo.Date)(nil)

	// If birth date is provided, create date value object
	if req.BirthDate != "" {
		birthDt, err := sharedvo.NewDate(req.BirthDate)
		if err != nil {
			return nil, err
		}

		birthDate = &birthDt
	}

	// Create person entity
	person, err := entity.NewPerson(
		time.Now().Unix(),
		req.Name,
		req.Phone,
		taxID,
		email,
		birthDate,
	)
	if err != nil {
		return nil, err
	}

	// Save person
	if err := uc.personRepo.Save(ctx, person); err != nil {
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
