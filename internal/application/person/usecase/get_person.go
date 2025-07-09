package usecase

import (
	"context"
	"todolist/internal/domain/person/dto"
	"todolist/internal/domain/person/repository"
	"todolist/internal/domain/shared"
)

// GetPersonUseCase handles retrieving a person
type GetPersonUseCase interface {
	Execute(ctx context.Context, personID int64) (*dto.PersonResponse, error)
}

type getPersonUseCase struct {
	personRepo repository.PersonRepository
}

// NewGetPersonUseCase creates a new instance of GetPersonUseCase
func NewGetPersonUseCase(personRepo repository.PersonRepository) GetPersonUseCase {
	return &getPersonUseCase{
		personRepo: personRepo,
	}
}

// Execute retrieves a person by ID
func (uc *getPersonUseCase) Execute(ctx context.Context, personID int64) (*dto.PersonResponse, error) {
	person, err := uc.personRepo.FindByID(ctx, personID)
	if err != nil {
		return nil, shared.ErrNotFound
	}

	return toPersonResponse(person), nil
}
