package repository

import (
	"context"
	"todolist/internal/domain/shared"
	"todolist/internal/domain/user/entity"
	"todolist/internal/domain/user/valueobject"
)

// UserRepository defines persistence operations for User
type UserRepository interface {
	// Commands
	Save(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id int64) error

	// Queries
	FindByID(ctx context.Context, id int64) (*entity.User, error)
	FindByUsername(ctx context.Context, username string) (*entity.User, error)
	FindByPersonID(ctx context.Context, personID int64) (*entity.User, error)

	// Validations
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	ExistsByPersonID(ctx context.Context, personID int64) (bool, error)
}

// UserQueryRepository defines complex query operations for User
type UserQueryRepository interface {
	// List operations
	FindAll(ctx context.Context, options shared.QueryOptions) ([]*entity.User, error)
	FindByStatus(ctx context.Context, status valueobject.UserStatus, options shared.QueryOptions) ([]*entity.User, error)
	FindByRole(ctx context.Context, role valueobject.UserRole, options shared.QueryOptions) ([]*entity.User, error)
	FindInactiveUsers(ctx context.Context, inactiveDays int, options shared.QueryOptions) ([]*entity.User, error)

	// Aggregations
	Count(ctx context.Context, filters []shared.Filter) (int64, error)
	CountByStatus(ctx context.Context) (map[valueobject.UserStatus]int64, error)
	CountByRole(ctx context.Context) (map[valueobject.UserRole]int64, error)
}
