package repository

import (
	"context"
	"errors"
	"todolist/internal/domain/person/entity"
	"todolist/internal/domain/person/repository"
	"todolist/internal/domain/shared"
	"todolist/internal/infrastructure/database/mapper"
	"todolist/internal/infrastructure/database/model"

	"gorm.io/gorm"
)

// personRepository implements repository.PersonRepository
type personRepository struct {
	db     *gorm.DB
	mapper *mapper.PersonMapper
}

// NewPersonRepository creates a new person repository
func NewPersonRepository(db *gorm.DB) repository.PersonRepository {
	return &personRepository{
		db:     db,
		mapper: mapper.NewPersonMapper(),
	}
}

// Save saves or updates a person
func (r *personRepository) Save(ctx context.Context, person *entity.Person) error {
	personModel := r.mapper.ToModel(person)

	// Use Save to handle both create and update
	if err := r.db.WithContext(ctx).Save(personModel).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return shared.ErrDuplicateEntry
		}
		return err
	}

	return nil
}

// Delete deletes a person (soft delete)
func (r *personRepository) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Delete(&model.Person{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return shared.ErrNotFound
	}

	return nil
}

// FindByID finds a person by ID
func (r *personRepository) FindByID(ctx context.Context, id int64) (*entity.Person, error) {
	person := &model.Person{}

	if err := r.db.WithContext(ctx).First(person, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.ErrNotFound
		}
		return nil, err
	}

	return r.mapper.ToDomain(person)
}

// FindByEmail finds a person by email
func (r *personRepository) FindByEmail(ctx context.Context, email string) (*entity.Person, error) {
	person := &model.Person{}

	if err := r.db.WithContext(ctx).Where("email = ?", email).First(person).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.ErrNotFound
		}
		return nil, err
	}

	return r.mapper.ToDomain(person)
}

// ExistsByTaxID checks if a person tax id exists
func (r *personRepository) ExistsByTaxID(ctx context.Context, taxID string) (bool, error) {
	var count int64

	if err := r.db.WithContext(ctx).Model(&model.Person{}).
		Where("tax_id = ?", taxID).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
