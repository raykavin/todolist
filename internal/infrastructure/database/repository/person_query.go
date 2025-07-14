package repository

import (
	"context"
	"todolist/internal/domain/person/entity"
	"todolist/internal/domain/person/repository"
	"todolist/internal/domain/shared"
	"todolist/internal/infrastructure/database"
	"todolist/internal/infrastructure/database/mapper"
	"todolist/internal/infrastructure/database/model"

	"gorm.io/gorm"
)

// personQueryRepository implements repository.PersonQueryRepository
type personQueryRepository struct {
	db     *gorm.DB
	mapper *mapper.PersonMapper
}

// NewPersonQueryRepository creates a new person query repository
func NewPersonQueryRepository(db *gorm.DB) repository.PersonQueryRepository {
	return &personQueryRepository{
		db:     db,
		mapper: mapper.NewPersonMapper(),
	}
}

// FindAll finds all persons with pagination
func (r *personQueryRepository) FindAll(ctx context.Context, options shared.QueryOptions) ([]*entity.Person, error) {
	people := []*model.Person{}

	query := r.db.WithContext(ctx).Model(&model.Person{})
	query = database.ApplyQueryOptions(query, options)

	if err := query.Find(&people).Error; err != nil {
		return nil, err
	}

	return r.mapper.ToDomainList(people)
}

// Search searches persons by name or email
func (r *personQueryRepository) Search(ctx context.Context, searchQuery string, options shared.QueryOptions) ([]*entity.Person, error) {
	people := []*model.Person{}

	query := r.db.WithContext(ctx).Model(&model.Person{})
	query = database.BuildSearchQuery(query, searchQuery, "name", "email", "phone")
	query = database.ApplyQueryOptions(query, options)

	if err := query.Find(&people).Error; err != nil {
		return nil, err
	}

	return r.mapper.ToDomainList(people)
}

// Count counts persons with filters
func (r *personQueryRepository) Count(ctx context.Context, filters []shared.Filter) (int64, error) {
	var count int64

	query := r.db.WithContext(ctx).Model(&model.Person{})
	query = database.ApplyFilters(query, filters)

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}
