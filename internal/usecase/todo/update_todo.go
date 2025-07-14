package usecase

import (
	"context"
	"todolist/internal/domain/shared"
	sharedvo "todolist/internal/domain/shared/valueobject"
	"todolist/internal/domain/todo/repository"
	"todolist/internal/domain/todo/service"
	vo "todolist/internal/domain/todo/valueobject"
	"todolist/internal/dto"
)

// UpdateTodoUseCase handles updating todos
type UpdateTodoUseCase interface {
	Execute(ctx context.Context, userID, todoID int64, input dto.UpdateTodoRequest) (*dto.TodoResponse, error)
}

type updateTodoUseCase struct {
	todoRepository repository.TodoRepository
	todoService    service.TodoService
}

// NewUpdateTodoUseCase creates a new instance of UpdateTodoUseCase
func NewUpdateTodoUseCase(
	todoRepository repository.TodoRepository,
	todoService service.TodoService,
) UpdateTodoUseCase {
	return &updateTodoUseCase{
		todoRepository: todoRepository,
		todoService:    todoService,
	}
}

// Execute updates a todo
func (uc *updateTodoUseCase) Execute(
	ctx context.Context,
	userID, todoID int64,
	input dto.UpdateTodoRequest,
) (*dto.TodoResponse, error) {
	// Validate user ownership
	if err := uc.todoService.ValidateUserOwnership(ctx, todoID, userID); err != nil {
		return nil, err
	}

	// Get the todo
	todo, err := uc.todoRepository.FindByID(ctx, todoID)
	if err != nil {
		return nil, shared.ErrNotFound
	}

	// Update title if provided
	if input.Title != nil {
		title, err := vo.NewTodoTitle(*input.Title)
		if err != nil {
			return nil, err
		}
		todo.UpdateTitle(title)
	}

	// Update description if provided
	if input.Description != nil {
		description, err := vo.NewTodoDescription(*input.Description)
		if err != nil {
			return nil, err
		}
		todo.UpdateDescription(description)
	}

	// Update priority if provided
	if input.Priority != nil {
		priority, err := sharedvo.NewPriorityFromString(*input.Priority)
		if err != nil {
			return nil, err
		}
		todo.UpdatePriority(priority)
	}

	// Update due date if provided
	if input.DueDate != nil {
		if err := todo.UpdateDueDate(input.DueDate); err != nil {
			return nil, err
		}
	}

	// Update status if provided
	if input.Status != nil {
		status, err := vo.NewTodoStatusFromString(*input.Status)
		if err != nil {
			return nil, err
		}
		if err := todo.ChangeStatus(status); err != nil {
			return nil, err
		}
	}

	// Save updated todo
	if err := uc.todoRepository.Save(ctx, todo); err != nil {
		return nil, err
	}

	return toTodoResponse(todo), nil
}
