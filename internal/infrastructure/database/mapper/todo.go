package mapper

import (
	sharedvo "todolist/internal/domain/shared/valueobject"
	"todolist/internal/domain/todo/entity"
	vo "todolist/internal/domain/todo/valueobject"
	"todolist/internal/infrastructure/database/model"
)

// TodoMapper handles conversion between domain entity and database model
type TodoMapper struct{}

// NewTodoMapper creates a new TodoMapper
func NewTodoMapper() *TodoMapper {
	return &TodoMapper{}
}

// ToModel converts domain entity to database model
func (m *TodoMapper) ToModel(todo *entity.Todo) *model.Todo {
	mdl := &model.Todo{
		ID:          todo.ID(),
		UserID:      todo.UserID(),
		Title:       todo.Title().Value(),
		Description: todo.Description().Value(),
		Status:      string(todo.Status()),
		Priority:    int8(todo.Priority()),
		DueDate:     todo.DueDate(),
		CompletedAt: todo.CompletedAt(),
		CreatedAt:   todo.CreatedAt(),
		UpdatedAt:   todo.UpdatedAt(),
	}

	// Convert tags
	tags := todo.Tags()
	if len(tags) > 0 {
		mdl.Tags = make([]*model.Tag, 0, len(tags))
		for _, tagName := range tags {
			mdl.Tags = append(mdl.Tags, &model.Tag{Name: tagName})
		}
	}

	return mdl
}

// ToDomain converts database model to domain entity
func (m *TodoMapper) ToDomain(model *model.Todo) (*entity.Todo, error) {
	title, err := vo.NewTodoTitle(model.Title)
	if err != nil {
		return nil, err
	}

	description, err := vo.NewTodoDescription(model.Description)
	if err != nil {
		return nil, err
	}

	priority := sharedvo.Priority(model.Priority)
	if !priority.IsValid() {
		priority = sharedvo.PriorityLow
	}

	todo, err := entity.NewTodo(
		model.ID,
		model.UserID,
		title,
		description,
		priority,
		model.DueDate,
	)
	if err != nil {
		return nil, err
	}

	// Set status
	status := vo.TodoStatus(model.Status)
	if status != vo.StatusPending && status.IsValid() {
		if err := todo.ChangeStatus(status); err != nil {
			// If status change fails, just keep the default
		}
	}

	// Set completed time if exists
	if model.CompletedAt != nil && todo.Status() == vo.StatusCompleted {
		// Use reflection or add a method to set CompletedAt directly
		// For now, the ChangeStatus method already handles this
	}

	// Set tags
	if len(model.Tags) > 0 {
		for _, tag := range model.Tags {
			todo.AddTag(tag.Name)
		}
	}

	// Set timestamps from database
	todo.Entity.SetCreatedAt(model.CreatedAt)
	todo.Entity.SetUpdatedAt(model.UpdatedAt)

	return todo, nil
}

// ToDomainList converts a list of model to domain entities
func (m *TodoMapper) ToDomainList(model []*model.Todo) ([]*entity.Todo, error) {
	todos := make([]*entity.Todo, 0, len(model))

	for _, model := range model {
		todo, err := m.ToDomain(model)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	return todos, nil
}
