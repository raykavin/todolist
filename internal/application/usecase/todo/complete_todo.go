package usecase

import (
	"context"
	"todolist/internal/application/dto"
	"todolist/internal/domain/shared"
	"todolist/internal/domain/todo/repository"
	"todolist/internal/domain/todo/service"
)

// CompleteTodoUseCase handles completing todos
type CompleteTodoUseCase interface {
	Execute(ctx context.Context, userID, todoID int64) (*dto.TodoResponse, error)
}

type completeTodoUseCase struct {
	todoRepo    repository.TodoRepository
	todoService service.TodoService
}

// NewCompleteTodoUseCase creates a new instance of CompleteTodoUseCase
func NewCompleteTodoUseCase(
	todoRepo repository.TodoRepository,
	todoService service.TodoService,
) CompleteTodoUseCase {
	return &completeTodoUseCase{
		todoRepo:    todoRepo,
		todoService: todoService,
	}
}

// Execute completes a todo
func (uc *completeTodoUseCase) Execute(ctx context.Context, userID, todoID int64) (*dto.TodoResponse, error) {
	// Validate user ownership
	if err := uc.todoService.ValidateUserOwnership(ctx, todoID, userID); err != nil {
		return nil, err
	}

	// Get the todo
	todo, err := uc.todoRepo.FindByID(ctx, todoID)
	if err != nil {
		return nil, shared.ErrNotFound
	}

	// Complete the todo
	if err := todo.Complete(); err != nil {
		return nil, err
	}

	// Save updated todo
	if err := uc.todoRepo.Save(ctx, todo); err != nil {
		return nil, err
	}

	return toTodoResponse(todo), nil
}
