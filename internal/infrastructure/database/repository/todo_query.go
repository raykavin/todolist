package repository

import (
	"context"
	"time"
	"todolist/internal/domain/shared"
	sharedvo "todolist/internal/domain/shared/valueobject"
	"todolist/internal/domain/todo/entity"
	"todolist/internal/domain/todo/repository"
	vo "todolist/internal/domain/todo/valueobject"
	"todolist/internal/infrastructure/database"
	"todolist/internal/infrastructure/database/mapper"
	"todolist/internal/infrastructure/database/model"

	"gorm.io/gorm"
)

// todoQueryRepository implements repository.TodoQueryRepository
type todoQueryRepository struct {
	db     *gorm.DB
	mapper *mapper.TodoMapper
}

// NewTodoQueryRepository creates a new todo query repository
func NewTodoQueryRepository(db *gorm.DB) repository.TodoQueryRepository {
	return &todoQueryRepository{
		db:     db,
		mapper: mapper.NewTodoMapper(),
	}
}

// FindAll finds all todos with pagination
func (r *todoQueryRepository) FindAll(ctx context.Context, options shared.QueryOptions) ([]*entity.Todo, error) {
	users := []*model.Todo{}

	query := r.db.WithContext(ctx).Model(&model.Todo{}).Preload("Tags")
	query = database.ApplyQueryOptions(query, options)

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	return r.mapper.ToDomainList(users)
}

// FindByFilters finds todos by multiple filters
func (r *todoQueryRepository) FindByFilters(
	ctx context.Context,
	filters vo.TodoFilterCriteria,
	options shared.QueryOptions,
) ([]*entity.Todo, error) {
	users := []*model.Todo{}

	query := r.db.WithContext(ctx).Model(&model.Todo{}).Preload("Tags")

	// Apply filters
	if filters.UserID != 0 {
		query = query.Where("user_id = ?", filters.UserID)
	}

	if len(filters.Status) > 0 {
		query = query.Where("status IN ?", filters.Status)
	}

	if len(filters.Priority) > 0 {
		priorities := make([]int, len(filters.Priority))
		for i, p := range filters.Priority {
			if priority, err := sharedvo.NewPriorityFromString(p); err == nil {
				priorities[i] = int(priority)
			}
		}
		query = query.Where("priority IN ?", priorities)
	}

	if filters.IsOverdue != nil && *filters.IsOverdue {
		query = query.Where(
			"due_date < ? AND status IN ?",
			time.Now(),
			[]string{"pending", "in_progress"},
		)
	}

	if filters.DueDateFrom != nil {
		query = query.Where("due_date >= ?", *filters.DueDateFrom)
	}

	if filters.DueDateTo != nil {
		query = query.Where("due_date <= ?", *filters.DueDateTo)
	}

	if filters.SearchTerm != "" {
		query = database.BuildSearchQuery(
			query,
			filters.SearchTerm,
			"title",
			"description",
		)
	}

	// Handle tag filters
	if len(filters.Tags) > 0 {
		query = query.Joins("JOIN todo_tags ON todo_tags.todo_id = todos.id").
			Joins("JOIN tags ON tags.id = todo_tags.tag_id").
			Where("tags.name IN ?", filters.Tags).
			Group("todos.id")
	}

	query = database.ApplyQueryOptions(query, options)

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	return r.mapper.ToDomainList(users)
}

// FindByUserAndStatus finds todos by user and status
func (r *todoQueryRepository) FindByUserAndStatus(
	ctx context.Context,
	userID int64,
	status vo.TodoStatus,
	options shared.QueryOptions,
) ([]*entity.Todo, error) {
	users := []*model.Todo{}

	query := r.db.WithContext(ctx).
		Model(&model.Todo{}).
		Preload("Tags").
		Where("user_id = ? AND status = ?", userID, string(status))
	query = database.ApplyQueryOptions(query, options)

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	return r.mapper.ToDomainList(users)
}

// FindByUserAndPriority finds todos by user and priority
func (r *todoQueryRepository) FindByUserAndPriority(
	ctx context.Context,
	userID int64,
	priority sharedvo.Priority,
	options shared.QueryOptions,
) ([]*entity.Todo, error) {
	users := []*model.Todo{}

	query := r.db.WithContext(ctx).
		Model(&model.Todo{}).
		Preload("Tags").
		Where("user_id = ? AND priority = ?", userID, int(priority))
	query = database.ApplyQueryOptions(query, options)

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	return r.mapper.ToDomainList(users)
}

// FindOverdue finds overdue todos
func (r *todoQueryRepository) FindOverdue(
	ctx context.Context,
	userID int64,
	options shared.QueryOptions,
) ([]*entity.Todo, error) {
	users := []*model.Todo{}

	query := r.db.WithContext(ctx).
		Model(&model.Todo{}).
		Preload("Tags").
		Where("user_id = ? AND due_date < ? AND status IN ?",
			userID, time.Now(), []string{"pending", "in_progress"})
	query = database.ApplyQueryOptions(query, options)

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	return r.mapper.ToDomainList(users)
}

// FindDueToday finds todos due today
func (r *todoQueryRepository) FindDueToday(ctx context.Context, userID int64) ([]*entity.Todo, error) {
	users := []*model.Todo{}

	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)

	if err := r.db.WithContext(ctx).
		Model(&model.Todo{}).
		Preload("Tags").
		Where("user_id = ? AND due_date >= ? AND due_date < ?", userID, today, tomorrow).
		Find(&users).Error; err != nil {
		return nil, err
	}

	return r.mapper.ToDomainList(users)
}

// FindDueBetween finds todos due between dates
func (r *todoQueryRepository) FindDueBetween(
	ctx context.Context,
	userID int64,
	start, end time.Time,
	options shared.QueryOptions,
) ([]*entity.Todo, error) {
	users := []*model.Todo{}

	query := r.db.WithContext(ctx).
		Model(&model.Todo{}).
		Preload("Tags").
		Where("user_id = ? AND due_date >= ? AND due_date <= ?", userID, start, end)
	query = database.ApplyQueryOptions(query, options)

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	return r.mapper.ToDomainList(users)
}

// FindByTag finds todos by tag
func (r *todoQueryRepository) FindByTag(
	ctx context.Context,
	userID int64,
	tag string,
	options shared.QueryOptions,
) ([]*entity.Todo, error) {
	users := []*model.Todo{}

	query := r.db.WithContext(ctx).
		Model(&model.Todo{}).
		Preload("Tags").
		Joins("JOIN todo_tags ON todo_tags.todo_id = todos.id").
		Joins("JOIN tags ON tags.id = todo_tags.tag_id").
		Where("todos.user_id = ? AND tags.name = ?", userID, tag)
	query = database.ApplyQueryOptions(query, options)

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	return r.mapper.ToDomainList(users)
}

// FindByTags finds todos by multiple tags
func (r *todoQueryRepository) FindByTags(
	ctx context.Context,
	userID int64,
	tags []string,
	options shared.QueryOptions,
) ([]*entity.Todo, error) {
	users := []*model.Todo{}

	query := r.db.WithContext(ctx).
		Model(&model.Todo{}).
		Preload("Tags").
		Joins("JOIN todo_tags ON todo_tags.todo_id = todos.id").
		Joins("JOIN tags ON tags.id = todo_tags.tag_id").
		Where("todos.user_id = ? AND tags.name IN ?", userID, tags).
		Group("todos.id").
		Having("COUNT(DISTINCT tags.name) = ?", len(tags))

	query = database.ApplyQueryOptions(query, options)

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	return r.mapper.ToDomainList(users)
}

// Search searches todos by title or description
func (r *todoQueryRepository) Search(
	ctx context.Context,
	userID int64,
	searchQuery string,
	options shared.QueryOptions,
) ([]*entity.Todo, error) {
	users := []*model.Todo{}

	query := r.db.WithContext(ctx).
		Model(&model.Todo{}).
		Preload("Tags").
		Where("user_id = ?", userID)
	query = database.BuildSearchQuery(query, searchQuery, "title", "description")
	query = database.ApplyQueryOptions(query, options)

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	return r.mapper.ToDomainList(users)
}

// Count counts todos with filters
func (r *todoQueryRepository) Count(ctx context.Context, filters []shared.Filter) (int64, error) {
	var count int64

	query := r.db.WithContext(ctx).Model(&model.Todo{})
	query = database.ApplyFilters(query, filters)

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// CountByStatus counts todos by status for a user
func (r *todoQueryRepository) CountByStatus(ctx context.Context, userID int64) (map[vo.TodoStatus]int64, error) {
	var results []struct {
		Status string
		Count  int64
	}

	if err := r.db.WithContext(ctx).
		Model(&model.Todo{}).
		Select("status, COUNT(*) as count").
		Where("user_id = ?", userID).
		Group("status").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	counts := make(map[vo.TodoStatus]int64)
	for _, result := range results {
		status := vo.TodoStatus(result.Status)
		if status.IsValid() {
			counts[status] = result.Count
		}
	}

	return counts, nil
}

// CountByPriority counts todos by priority for a user
func (r *todoQueryRepository) CountByPriority(ctx context.Context, userID int64) (map[sharedvo.Priority]int64, error) {
	var results []struct {
		Priority int
		Count    int64
	}

	if err := r.db.WithContext(ctx).
		Model(&model.Todo{}).
		Select("priority, COUNT(*) as count").
		Where("user_id = ?", userID).
		Group("priority").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	counts := make(map[sharedvo.Priority]int64)
	for _, result := range results {
		priority := sharedvo.Priority(result.Priority)
		if priority.IsValid() {
			counts[priority] = result.Count
		}
	}

	return counts, nil
}

// GetStatistics gets todo statistics for a user
func (r *todoQueryRepository) GetStatistics(ctx context.Context, userID int64) (*vo.TodoStatistics, error) {
	stats := &vo.TodoStatistics{
		ByStatus:   make(map[string]int64),
		ByPriority: make(map[string]int64),
	}

	// Get total count
	if err := r.db.WithContext(ctx).
		Model(&model.Todo{}).
		Where("user_id = ?", userID).
		Count(&stats.Total).Error; err != nil {
		return nil, err
	}

	// Get counts by status
	statusCounts, err := r.CountByStatus(ctx, userID)
	if err != nil {
		return nil, err
	}
	for status, count := range statusCounts {
		stats.ByStatus[status.String()] = count
	}

	// Get counts by priority
	priorityCounts, err := r.CountByPriority(ctx, userID)
	if err != nil {
		return nil, err
	}
	for priority, count := range priorityCounts {
		stats.ByPriority[priority.String()] = count
	}

	// Get overdue count
	if err := r.db.WithContext(ctx).
		Model(&model.Todo{}).
		Where("user_id = ? AND due_date < ? AND status IN ?",
			userID, time.Now(), []string{"pending", "in_progress"}).
		Count(&stats.Overdue).Error; err != nil {
		return nil, err
	}

	// Get due today count
	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)
	if err := r.db.WithContext(ctx).
		Model(&model.Todo{}).
		Where("user_id = ? AND due_date >= ? AND due_date < ?", userID, today, tomorrow).
		Count(&stats.DueToday).Error; err != nil {
		return nil, err
	}

	// Get due this week count
	weekEnd := today.Add(7 * 24 * time.Hour)
	if err := r.db.WithContext(ctx).
		Model(&model.Todo{}).
		Where("user_id = ? AND due_date >= ? AND due_date < ?", userID, today, weekEnd).
		Count(&stats.DueThisWeek).Error; err != nil {
		return nil, err
	}

	// Get completed today count
	if err := r.db.WithContext(ctx).
		Model(&model.Todo{}).
		Where("user_id = ? AND status = ? AND completed_at >= ? AND completed_at < ?",
			userID, "completed", today, tomorrow).
		Count(&stats.CompletedToday).Error; err != nil {
		return nil, err
	}

	// Calculate completion rate
	if stats.Total > 0 {
		completed := stats.ByStatus["completed"]
		stats.CompletionRate = float64(completed) / float64(stats.Total) * 100
	}

	return stats, nil
}

// GetPopularTags gets the most used tags for a user
func (r *todoQueryRepository) GetPopularTags(ctx context.Context, userID int64, limit int) ([]vo.TagCount, error) {
	tagCounts := []vo.TagCount{}

	query := `
		SELECT t.name as tag, COUNT(DISTINCT tt.todo_id) as count
		FROM tags t
		JOIN todo_tags tt ON tt.tag_id = t.id
		JOIN todos td ON td.id = tt.todo_id
		WHERE td.user_id = ? AND td.deleted_at IS NULL
		GROUP BY t.name
		ORDER BY count DESC
		LIMIT ?
	`

	if err := r.db.WithContext(ctx).
		Raw(query, userID, limit).
		Scan(&tagCounts).Error; err != nil {
		return nil, err
	}

	return tagCounts, nil
}
