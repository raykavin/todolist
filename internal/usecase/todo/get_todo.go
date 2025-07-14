package usecase

import (
	"context"
	"todolist/internal/domain/shared"
	"todolist/internal/domain/todo/repository"
	"todolist/internal/domain/todo/service"
	"todolist/internal/dto"
)

// GetTodoUseCase handles retrieving a single todo
type GetTodoUseCase interface {
	Execute(ctx context.Context, userID, todoID int64) (*dto.TodoResponse, error)
}

type getTodoUseCase struct {
	todoRepository repository.TodoRepository
	todoService    service.TodoService
}

// NewGetTodoUseCase creates a new instance of GetTodoUseCase
func NewGetTodoUseCase(
	todoRepository repository.TodoRepository,
	todoService service.TodoService,
) GetTodoUseCase {
	return &getTodoUseCase{
		todoRepository: todoRepository,
		todoService:    todoService,
	}
}

// Execute retrieves a todo by ID
func (uc *getTodoUseCase) Execute(ctx context.Context, userID, todoID int64) (*dto.TodoResponse, error) {
	// Validate user ownership
	if err := uc.todoService.ValidateUserOwnership(ctx, todoID, userID); err != nil {
		return nil, err
	}

	// Get the todo
	todo, err := uc.todoRepository.FindByID(ctx, todoID)
	if err != nil {
		return nil, shared.ErrNotFound
	}

	return toTodoResponse(todo), nil
}
