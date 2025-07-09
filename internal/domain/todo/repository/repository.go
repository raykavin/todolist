package repository

import (
	"context"
	"time"
	"todolist/internal/domain/shared"
	sharedvo "todolist/internal/domain/shared/valueobject"
	"todolist/internal/domain/todo/dto"
	"todolist/internal/domain/todo/entity"
	vo "todolist/internal/domain/todo/valueobject"
)

// TodoRepository defines persistence operations for Todo
type TodoRepository interface {
	// Commands
	Save(ctx context.Context, todo *entity.Todo) error
	Delete(ctx context.Context, id int64) error

	// Queries
	FindByID(ctx context.Context, id int64) (*entity.Todo, error)
	FindByUserID(ctx context.Context, userID int64) ([]*entity.Todo, error)

	// Batch operations
	DeleteByUserID(ctx context.Context, userID int64) error
}

// TodoQueryRepository defines complex query operations for Todo
type TodoQueryRepository interface {
	// List operations
	FindAll(ctx context.Context, options shared.QueryOptions) ([]*entity.Todo, error)
	FindByFilters(ctx context.Context, filters dto.TodoFilter, options shared.QueryOptions) ([]*entity.Todo, error)

	// User-specific queries
	FindByUserAndStatus(ctx context.Context, userID int64, status vo.TodoStatus, options shared.QueryOptions) ([]*entity.Todo, error)
	FindByUserAndPriority(ctx context.Context, userID int64, priority sharedvo.Priority, options shared.QueryOptions) ([]*entity.Todo, error)

	// Date-based queries
	FindOverdue(ctx context.Context, userID int64, options shared.QueryOptions) ([]*entity.Todo, error)
	FindDueToday(ctx context.Context, userID int64) ([]*entity.Todo, error)
	FindDueBetween(ctx context.Context, userID int64, start, end time.Time, options shared.QueryOptions) ([]*entity.Todo, error)

	// Tag queries
	FindByTag(ctx context.Context, userID int64, tag string, options shared.QueryOptions) ([]*entity.Todo, error)
	FindByTags(ctx context.Context, userID int64, tags []string, options shared.QueryOptions) ([]*entity.Todo, error)

	// Search
	Search(ctx context.Context, userID int64, query string, options shared.QueryOptions) ([]*entity.Todo, error)

	// Aggregations
	Count(ctx context.Context, filters []shared.Filter) (int64, error)
	CountByStatus(ctx context.Context, userID int64) (map[vo.TodoStatus]int64, error)
	CountByPriority(ctx context.Context, userID int64) (map[sharedvo.Priority]int64, error)
	GetStatistics(ctx context.Context, userID int64) (*dto.TodoStatistics, error)

	// Tag aggregations
	GetPopularTags(ctx context.Context, userID int64, limit int) ([]dto.TagCount, error)
}
