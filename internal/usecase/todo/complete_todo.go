package usecase

import (
	"context"
	"todolist/internal/domain/shared"
	"todolist/internal/domain/todo/repository"
	"todolist/internal/domain/todo/service"
	"todolist/internal/dto"
)

// CompleteTodoUseCase handles completing todos
type CompleteTodoUseCase interface {
	Execute(ctx context.Context, userID, todoID int64) (*dto.TodoResponse, error)
}

type completeTodoUseCase struct {
	todoRepository repository.TodoRepository
	todoService    service.TodoService
}

// NewCompleteTodoUseCase creates a new instance of CompleteTodoUseCase
func NewCompleteTodoUseCase(
	todoRepository repository.TodoRepository,
	todoService service.TodoService,
) CompleteTodoUseCase {
	return &completeTodoUseCase{
		todoRepository: todoRepository,
		todoService:    todoService,
	}
}

// Execute completes a todo
func (uc *completeTodoUseCase) Execute(ctx context.Context, userID, todoID int64) (*dto.TodoResponse, error) {
	// Validate user ownership
	if err := uc.todoService.ValidateUserOwnership(ctx, todoID, userID); err != nil {
		return nil, err
	}

	// Get the todo
	todo, err := uc.todoRepository.FindByID(ctx, todoID)
	if err != nil {
		return nil, shared.ErrNotFound
	}

	// Complete the todo
	if err := todo.Complete(); err != nil {
		return nil, err
	}

	// Save updated todo
	if err := uc.todoRepository.Save(ctx, todo); err != nil {
		return nil, err
	}

	return toTodoResponse(todo), nil
}
