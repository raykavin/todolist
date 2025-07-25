package usecase

import (
	"context"
	"todolist/internal/domain/shared"
	"todolist/internal/domain/todo/repository"
	vo "todolist/internal/domain/todo/valueobject"
	"todolist/internal/dto"
)

// ListTodosUseCase handles listing todos with filters
type ListTodosUseCase interface {
	Execute(ctx context.Context, userID int64, filters vo.TodoFilterCriteria, options shared.QueryOptions) (*dto.TodoListResponse, error)
}

type listTodosUseCase struct {
	todoQueryRepository repository.TodoQueryRepository
}

// NewListTodosUseCase creates a new instance of ListTodosUseCase
func NewListTodosUseCase(todoQueryRepository repository.TodoQueryRepository) ListTodosUseCase {
	return &listTodosUseCase{
		todoQueryRepository: todoQueryRepository,
	}
}

// Execute lists todos based on filters
func (uc *listTodosUseCase) Execute(
	ctx context.Context,
	userID int64,
	filters vo.TodoFilterCriteria,
	options shared.QueryOptions,
) (*dto.TodoListResponse, error) {
	// Ensure user filter is set
	filters.UserID = userID

	// Get todos
	todos, err := uc.todoQueryRepository.FindByFilters(ctx, filters, options)
	if err != nil {
		return nil, err
	}

	// Get total count
	countFilters := []shared.Filter{
		{Field: "user_id", Operator: shared.FilterOperatorEqual, Value: userID},
	}
	totalCount, err := uc.todoQueryRepository.Count(ctx, countFilters)
	if err != nil {
		return nil, err
	}

	// Convert to response
	response := &dto.TodoListResponse{
		Todos:      make([]*dto.TodoResponse, len(todos)),
		TotalCount: totalCount,
		Page:       options.Offset/options.Limit + 1,
		PageSize:   options.Limit,
	}

	for i, todo := range todos {
		response.Todos[i] = toTodoResponse(todo)
	}

	return response, nil
}
