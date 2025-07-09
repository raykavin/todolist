package usecase

import (
	"context"
	"todolist/internal/domain/todo/repository"
	"todolist/internal/domain/todo/service"
)

// DeleteTodoUseCase handles deleting todos
type DeleteTodoUseCase interface {
	Execute(ctx context.Context, userID, todoID int64) error
}

type deleteTodoUseCase struct {
	todoRepo    repository.TodoRepository
	todoService service.TodoService
}

// NewDeleteTodoUseCase creates a new instance of DeleteTodoUseCase
func NewDeleteTodoUseCase(
	todoRepo repository.TodoRepository,
	todoService service.TodoService,
) DeleteTodoUseCase {
	return &deleteTodoUseCase{
		todoRepo:    todoRepo,
		todoService: todoService,
	}
}

// Execute deletes a todo
func (uc *deleteTodoUseCase) Execute(ctx context.Context, userID, todoID int64) error {
	// Validate user ownership
	if err := uc.todoService.ValidateUserOwnership(ctx, todoID, userID); err != nil {
		return err
	}

	// Delete the todo
	return uc.todoRepo.Delete(ctx, todoID)
}
