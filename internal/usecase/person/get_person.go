package usecase

import (
	"context"
	"todolist/internal/domain/person/repository"
	"todolist/internal/domain/shared"
	"todolist/internal/dto"
)

// GetPersonUseCase handles retrieving a person
type GetPersonUseCase interface {
	Execute(ctx context.Context, personID int64) (*dto.PersonResponse, error)
}

type getPersonUseCase struct {
	personRepository repository.PersonRepository
}

// NewGetPersonUseCase creates a new instance of GetPersonUseCase
func NewGetPersonUseCase(personRepository repository.PersonRepository) GetPersonUseCase {
	return &getPersonUseCase{
		personRepository: personRepository,
	}
}

// Execute retrieves a person by ID
func (uc *getPersonUseCase) Execute(ctx context.Context, personID int64) (*dto.PersonResponse, error) {
	person, err := uc.personRepository.FindByID(ctx, personID)
	if err != nil {
		return nil, shared.ErrNotFound
	}

	return toPersonResponse(person), nil
}
