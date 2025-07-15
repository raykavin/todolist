package repository

import (
	"context"
	"todolist/internal/domain/person/entity"
	"todolist/internal/domain/shared"
)

// PersonRepository defines persistence operations for Person
type PersonRepository interface {
	// Commands
	Save(ctx context.Context, person *entity.Person) error
	Delete(ctx context.Context, id int64) error

	// Queries
	FindByID(ctx context.Context, id int64) (*entity.Person, error)
	FindByEmail(ctx context.Context, email string) (*entity.Person, error)

	// Validations
	ExistsByTaxID(ctx context.Context, taxID string) (bool, error)
}

// PersonQueryRepository defines complex query operations for Person
type PersonQueryRepository interface {
	// List operations
	FindAll(ctx context.Context, options shared.QueryOptions) ([]*entity.Person, error)

	// Search
	Search(ctx context.Context, query string, options shared.QueryOptions) ([]*entity.Person, error)

	// Count
	Count(ctx context.Context, filters []shared.Filter) (int64, error)
}
