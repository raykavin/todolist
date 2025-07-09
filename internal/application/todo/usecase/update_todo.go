package usecase

import (
	"context"
	"todolist/internal/domain/shared"
	sharedvo "todolist/internal/domain/shared/valueobject"
	"todolist/internal/domain/todo/dto"
	"todolist/internal/domain/todo/repository"
	"todolist/internal/domain/todo/service"
	vo "todolist/internal/domain/todo/valueobject"
)

// UpdateTodoUseCase handles updating todos
type UpdateTodoUseCase interface {
	Execute(ctx context.Context, userID, todoID int64, req dto.UpdateTodoRequest) (*dto.TodoResponse, error)
}

type updateTodoUseCase struct {
	todoRepo    repository.TodoRepository
	todoService service.TodoService
}

// NewUpdateTodoUseCase creates a new instance of UpdateTodoUseCase
func NewUpdateTodoUseCase(
	todoRepo repository.TodoRepository,
	todoService service.TodoService,
) UpdateTodoUseCase {
	return &updateTodoUseCase{
		todoRepo:    todoRepo,
		todoService: todoService,
	}
}

// Execute updates a todo
func (uc *updateTodoUseCase) Execute(
	ctx context.Context,
	userID, todoID int64,
	req dto.UpdateTodoRequest,
) (*dto.TodoResponse, error) {
	// Validate user ownership
	if err := uc.todoService.ValidateUserOwnership(ctx, todoID, userID); err != nil {
		return nil, err
	}

	// Get the todo
	todo, err := uc.todoRepo.FindByID(ctx, todoID)
	if err != nil {
		return nil, shared.ErrNotFound
	}

	// Update title if provided
	if req.Title != nil {
		title, err := vo.NewTodoTitle(*req.Title)
		if err != nil {
			return nil, err
		}
		todo.UpdateTitle(title)
	}

	// Update description if provided
	if req.Description != nil {
		description, err := vo.NewTodoDescription(*req.Description)
		if err != nil {
			return nil, err
		}
		todo.UpdateDescription(description)
	}

	// Update priority if provided
	if req.Priority != nil {
		priority, err := sharedvo.NewPriorityFromString(*req.Priority)
		if err != nil {
			return nil, err
		}
		todo.UpdatePriority(priority)
	}

	// Update due date if provided
	if req.DueDate != nil {
		if err := todo.UpdateDueDate(req.DueDate); err != nil {
			return nil, err
		}
	}

	// Update status if provided
	if req.Status != nil {
		status, err := vo.NewTodoStatusFromString(*req.Status)
		if err != nil {
			return nil, err
		}
		if err := todo.ChangeStatus(status); err != nil {
			return nil, err
		}
	}

	// Save updated todo
	if err := uc.todoRepo.Save(ctx, todo); err != nil {
		return nil, err
	}

	return toTodoResponse(todo), nil
}
