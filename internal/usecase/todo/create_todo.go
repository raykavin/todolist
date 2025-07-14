package usecase

import (
	"context"
	"time"
	sharedvo "todolist/internal/domain/shared/valueobject"
	"todolist/internal/domain/todo/entity"
	"todolist/internal/domain/todo/repository"
	"todolist/internal/domain/todo/service"
	vo "todolist/internal/domain/todo/valueobject"
	"todolist/internal/dto"
)

// CreateTodoUseCase handles the creation of new todos
type CreateTodoUseCase interface {
	Execute(ctx context.Context, userID int64, input dto.CreateTodoRequest) (*dto.TodoResponse, error)
}

type createTodoUseCase struct {
	todoRepository repository.TodoRepository
	todoService    service.TodoService
}

// NewCreateTodoUseCase creates a new instance of CreateTodoUseCase
func NewCreateTodoUseCase(
	todoRepository repository.TodoRepository,
	todoService service.TodoService,
) CreateTodoUseCase {
	return &createTodoUseCase{
		todoRepository: todoRepository,
		todoService:    todoService,
	}
}

// Execute creates a new todo
func (uc *createTodoUseCase) Execute(
	ctx context.Context,
	userID int64,
	input dto.CreateTodoRequest,
) (*dto.TodoResponse, error) {
	// Create value objects
	title, err := vo.NewTodoTitle(input.Title)
	if err != nil {
		return nil, err
	}

	description, err := vo.NewTodoDescription(input.Description)
	if err != nil {
		return nil, err
	}

	priority, err := sharedvo.NewPriorityFromString(input.Priority)
	if err != nil {
		return nil, err
	}

	// Suggest due date if not provided
	dueDate := input.DueDate
	if dueDate == nil && input.Priority != "low" {
		// Get user workload for smart suggestion
		suggestedDate, _ := uc.todoService.SuggestDueDate(ctx, input.Priority, nil)
		dueDate = suggestedDate
	}

	// Create todo entity
	todo, err := entity.NewTodo(
		time.Now().Unix(),
		userID,
		title,
		description,
		priority,
		dueDate,
	)
	if err != nil {
		return nil, err
	}

	// Add tags
	for _, tag := range input.Tags {
		todo.AddTag(tag)
	}

	// Save todo
	if err := uc.todoRepository.Save(ctx, todo); err != nil {
		return nil, err
	}

	// Convert to response
	return toTodoResponse(todo), nil
}

// Helper function to convert entity to DTO
func toTodoResponse(todo *entity.Todo) *dto.TodoResponse {
	return &dto.TodoResponse{
		ID:          todo.ID(),
		UserID:      todo.UserID(),
		Title:       todo.Title().Value(),
		Description: todo.Description().Value(),
		Status:      todo.Status().String(),
		Priority:    todo.Priority().String(),
		DueDate:     todo.DueDate(),
		CompletedAt: todo.CompletedAt(),
		Tags:        todo.Tags(),
		IsOverdue:   todo.IsOverdue(),
		CreatedAt:   todo.CreatedAt(),
		UpdatedAt:   todo.UpdatedAt(),
	}
}
