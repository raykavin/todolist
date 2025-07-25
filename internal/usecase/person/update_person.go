package usecase

import (
	"context"
	"todolist/internal/domain/person/repository"
	vo "todolist/internal/domain/person/valueobject"
	"todolist/internal/domain/shared"
	"todolist/internal/dto"
)

// UpdatePersonUseCase handles person updates
type UpdatePersonUseCase interface {
	Execute(ctx context.Context, personID int64, input dto.UpdatePersonRequest) (*dto.PersonResponse, error)
}

type updatePersonUseCase struct {
	personRepository repository.PersonRepository
}

// NewUpdatePersonUseCase creates a new instance of UpdatePersonUseCase
func NewUpdatePersonUseCase(personRepository repository.PersonRepository) UpdatePersonUseCase {
	return &updatePersonUseCase{
		personRepository: personRepository,
	}
}

// Execute updates a person
func (uc *updatePersonUseCase) Execute(ctx context.Context, personID int64, input dto.UpdatePersonRequest) (*dto.PersonResponse, error) {
	// Get the person
	person, err := uc.personRepository.FindByID(ctx, personID)
	if err != nil {
		return nil, shared.ErrNotFound
	}

	// Update email if provided
	if input.Email != nil {
		// Check if new email already exists
		existingPerson, _ := uc.personRepository.FindByEmail(ctx, *input.Email)
		if existingPerson != nil && existingPerson.ID() != personID {
			return nil, vo.ErrInvalidEmail
		}

		email, err := vo.NewEmail(*input.Email)
		if err != nil {
			return nil, err
		}
		person.SetEmail(email)
	}

	// Update phone if provided
	if input.Phone != nil {
		if err := person.SetPhone(*input.Phone); err != nil {
			return nil, err
		}
	}

	// Save updated person
	if err := uc.personRepository.Save(ctx, person); err != nil {
		return nil, err
	}

	return toPersonResponse(person), nil
}
