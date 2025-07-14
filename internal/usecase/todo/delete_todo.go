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
	todoRepository repository.TodoRepository
	todoService    service.TodoService
}

// NewDeleteTodoUseCase creates a new instance of DeleteTodoUseCase
func NewDeleteTodoUseCase(
	todoRepository repository.TodoRepository,
	todoService service.TodoService,
) DeleteTodoUseCase {
	return &deleteTodoUseCase{
		todoRepository: todoRepository,
		todoService:    todoService,
	}
}

// Execute deletes a todo
func (uc *deleteTodoUseCase) Execute(ctx context.Context, userID, todoID int64) error {
	// Validate user ownership
	if err := uc.todoService.ValidateUserOwnership(ctx, todoID, userID); err != nil {
		return err
	}

	// Delete the todo
	return uc.todoRepository.Delete(ctx, todoID)
}
