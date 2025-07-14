package service

import (
	"context"
	"time"
	"todolist/internal/domain/shared"
	"todolist/internal/domain/todo/entity"
	"todolist/internal/domain/todo/repository"
	vo "todolist/internal/domain/todo/valueobject"
)

// TodoService provides domain services for Todo operations
type TodoService interface {
	// Validation services
	ValidateUserOwnership(ctx context.Context, todoID, userID int64) error

	// Business logic services
	GetUserProductivity(ctx context.Context, userID int64, period time.Duration) (*ProductivityMetrics, error)
	SuggestDueDate(ctx context.Context, priority string, userWorkload map[time.Time]int) (*time.Time, error)

	// Bulk operations
	MarkOverdueAsInProgress(ctx context.Context, userID int64) (int, error)
	AutoCancelOldPendingTodos(ctx context.Context, userID int64, olderThan time.Duration) (int, error)
}

// ProductivityMetrics represents user productivity metrics
type ProductivityMetrics struct {
	TotalCreated          int64
	TotalCompleted        int64
	CompletionRate        float64
	AverageTimeToComplete time.Duration
	MostProductiveDay     string
	MostUsedTags          []string
}

// todoService implements TodoService
type todoService struct {
	todoRepository      repository.TodoRepository
	todoQueryRepository repository.TodoQueryRepository
}

// NewTodoService creates a new todo service
func NewTodoService(
	todoRepository repository.TodoRepository,
	todoQueryRepository repository.TodoQueryRepository,
) TodoService {
	return &todoService{
		todoRepository:      todoRepository,
		todoQueryRepository: todoQueryRepository,
	}
}

// ValidateUserOwnership validates if a user owns a todo
func (s *todoService) ValidateUserOwnership(ctx context.Context, todoID, userID int64) error {
	todo, err := s.todoRepository.FindByID(ctx, todoID)
	if err != nil {
		return entity.ErrTodoNotFound
	}

	if todo.UserID() != userID {
		return entity.ErrUnauthorizedTodoAccess
	}

	return nil
}

// GetUserProductivity calculates user productivity metrics
func (s *todoService) GetUserProductivity(ctx context.Context, userID int64, period time.Duration) (*ProductivityMetrics, error) {
	// This would involve complex queries and calculations
	// For now, returning a simplified version
	stats, err := s.todoQueryRepository.GetStatistics(ctx, userID)
	if err != nil {
		return nil, err
	}

	metrics := &ProductivityMetrics{
		TotalCreated:   stats.Total,
		TotalCompleted: stats.ByStatus["completed"],
		CompletionRate: stats.CompletionRate,
	}

	return metrics, nil
}

// SuggestDueDate suggests an optimal due date based on priority and user workload
func (s *todoService) SuggestDueDate(ctx context.Context, priority string, userWorkload map[time.Time]int) (*time.Time, error) {
	// Simple implementation: suggest dates based on priority
	now := time.Now()
	var suggestedDate time.Time

	switch priority {
	case "critical":
		suggestedDate = now.AddDate(0, 0, 1) // Tomorrow
	case "high":
		suggestedDate = now.AddDate(0, 0, 3) // 3 days
	case "medium":
		suggestedDate = now.AddDate(0, 0, 7) // 1 week
	case "low":
		suggestedDate = now.AddDate(0, 0, 14) // 2 weeks
	default:
		suggestedDate = now.AddDate(0, 0, 7) // Default to 1 week
	}

	// Find a day with lower workload
	for i := 0; i < 7; i++ {
		checkDate := suggestedDate.AddDate(0, 0, i)
		if workload, exists := userWorkload[checkDate]; !exists || workload < 3 {
			return &checkDate, nil
		}
	}

	return &suggestedDate, nil
}

// MarkOverdueAsInProgress marks all overdue todos as in progress
func (s *todoService) MarkOverdueAsInProgress(ctx context.Context, userID int64) (int, error) {
	// Create filter to find overdue pending todos
	isOverdue := true
	filter := vo.TodoFilterCriteria{
		UserID:    userID,
		Status:    []string{string(vo.StatusPending)},
		IsOverdue: &isOverdue,
	}

	// Find all overdue pending todos
	overdueTodos, err := s.todoQueryRepository.FindByFilters(ctx, filter, shared.QueryOptions{Limit: 1000})
	if err != nil {
		return 0, err
	}

	updatedCount := 0

	// Update each overdue todo to in_progress status
	for _, todo := range overdueTodos {
		// Verify todo is actually overdue and pending
		if !todo.IsOverdue() || todo.Status() != vo.StatusPending {
			continue
		}

		// Change status to in progress
		if err := todo.StartProgress(); err != nil {
			// Skip this todo and continue with others
			continue
		}

		// Save the updated todo
		if err := s.todoRepository.Save(ctx, todo); err != nil {
			// Skip this todo and continue with others
			continue
		}

		updatedCount++
	}

	return updatedCount, nil
}

// AutoCancelOldPendingTodos cancels old pending todos
func (s *todoService) AutoCancelOldPendingTodos(ctx context.Context, userID int64, olderThan time.Duration) (int, error) {
	// Calculate cutoff date
	cutoffDate := time.Now().Add(-olderThan)

	// Create filter to find pending todos with due date before cutoff
	filter := vo.TodoFilterCriteria{
		UserID:    userID,
		Status:    []string{string(vo.StatusPending)},
		DueDateTo: &cutoffDate,
	}

	// Find all old pending todos
	oldTodos, err := s.todoQueryRepository.FindByFilters(ctx, filter, shared.QueryOptions{Limit: 1000})
	if err != nil {
		return 0, err
	}

	cancelledCount := 0

	// Process each old todo
	for _, todo := range oldTodos {
		// Verify todo is pending
		if todo.Status() != vo.StatusPending {
			continue
		}

		// Check if todo is actually old enough
		dueDate := todo.DueDate()
		if dueDate == nil || !dueDate.Before(cutoffDate) {
			continue
		}

		// Cancel the todo
		if err := todo.Cancel(); err != nil {
			// Skip this todo and continue with others
			continue
		}

		// Save the updated todo
		if err := s.todoRepository.Save(ctx, todo); err != nil {
			// Skip this todo and continue with others
			continue
		}

		cancelledCount++
	}

	return cancelledCount, nil
}
