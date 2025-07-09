package entity

import (
	"errors"
	"time"
	"todolist/internal/domain/shared"
	sharedvo "todolist/internal/domain/shared/valueobject"
	vo "todolist/internal/domain/todo/valueobject"
)

var (
	ErrInvalidTodoTitle        = errors.New("todo title is required")
	ErrInvalidUserID           = errors.New("user ID is required")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrTodoAlreadyCompleted    = errors.New("todo is already completed")
	ErrInvalidDueDate          = errors.New("due date cannot be in the past")
)

// Todo represents a todo item
type Todo struct {
	shared.Entity
	userID      string
	title       vo.TodoTitle
	description vo.TodoDescription
	status      vo.TodoStatus
	priority    sharedvo.Priority
	dueDate     *time.Time
	completedAt *time.Time
	tags        []string
}

// NewTodo creates a new Todo entity
func NewTodo(
	id uint64,
	userID string,
	title vo.TodoTitle,
	description vo.TodoDescription,
	priority sharedvo.Priority,
	dueDate *time.Time,
) (*Todo, error) {
	if userID == "" {
		return nil, ErrInvalidUserID
	}

	if dueDate != nil && dueDate.Before(time.Now()) {
		return nil, ErrInvalidDueDate
	}

	return &Todo{
		Entity:      shared.NewEntity(id),
		userID:      userID,
		title:       title,
		description: description,
		status:      vo.StatusPending,
		priority:    priority,
		dueDate:     dueDate,
		tags:        []string{},
	}, nil
}

// Getters
func (t *Todo) UserID() string {
	return t.userID
}

func (t *Todo) Title() vo.TodoTitle {
	return t.title
}

func (t *Todo) Description() vo.TodoDescription {
	return t.description
}

func (t *Todo) Status() vo.TodoStatus {
	return t.status
}

func (t *Todo) Priority() sharedvo.Priority {
	return t.priority
}

func (t *Todo) DueDate() *time.Time {
	if t.dueDate == nil {
		return nil
	}
	dueDateCopy := *t.dueDate
	return &dueDateCopy
}

func (t *Todo) CompletedAt() *time.Time {
	if t.completedAt == nil {
		return nil
	}
	completedAtCopy := *t.completedAt
	return &completedAtCopy
}

func (t *Todo) Tags() []string {
	// Return a copy to prevent external modification
	tagsCopy := make([]string, len(t.tags))
	copy(tagsCopy, t.tags)
	return tagsCopy
}

// Business methods
func (t *Todo) IsCompleted() bool {
	return t.status == vo.StatusCompleted
}

func (t *Todo) IsOverdue() bool {
	if t.dueDate == nil || t.IsCompleted() {
		return false
	}
	return t.dueDate.Before(time.Now())
}

func (t *Todo) DaysUntilDue() *int {
	if t.dueDate == nil {
		return nil
	}
	days := int(time.Until(*t.dueDate).Hours() / 24)
	return &days
}

// Update methods
func (t *Todo) UpdateTitle(title vo.TodoTitle) {
	t.title = title
	t.SetAsModified()
}

func (t *Todo) UpdateDescription(description vo.TodoDescription) {
	t.description = description
	t.SetAsModified()
}

func (t *Todo) UpdatePriority(priority sharedvo.Priority) {
	t.priority = priority
	t.SetAsModified()
}

func (t *Todo) UpdateDueDate(dueDate *time.Time) error {
	if dueDate != nil && dueDate.Before(time.Now()) {
		return ErrInvalidDueDate
	}
	t.dueDate = dueDate
	t.SetAsModified()
	return nil
}

func (t *Todo) ChangeStatus(newStatus vo.TodoStatus) error {
	if !t.status.CanTransitionTo(newStatus) {
		return ErrInvalidStatusTransition
	}

	t.status = newStatus

	// Set completed time when marking as completed
	if newStatus == vo.StatusCompleted {
		now := time.Now()
		t.completedAt = &now
	} else if t.status == vo.StatusCompleted {
		// Clear completed time when moving away from completed
		t.completedAt = nil
	}

	t.SetAsModified()
	return nil
}

func (t *Todo) Complete() error {
	if t.IsCompleted() {
		return ErrTodoAlreadyCompleted
	}
	return t.ChangeStatus(vo.StatusCompleted)
}

func (t *Todo) StartProgress() error {
	return t.ChangeStatus(vo.StatusInProgress)
}

func (t *Todo) Cancel() error {
	return t.ChangeStatus(vo.StatusCancelled)
}

func (t *Todo) Reopen() error {
	return t.ChangeStatus(vo.StatusPending)
}

// Tag management
func (t *Todo) AddTag(tag string) {
	// Check if tag already exists
	for _, existingTag := range t.tags {
		if existingTag == tag {
			return
		}
	}
	t.tags = append(t.tags, tag)
	t.SetAsModified()
}

func (t *Todo) RemoveTag(tag string) {
	newTags := make([]string, 0, len(t.tags))
	for _, existingTag := range t.tags {
		if existingTag != tag {
			newTags = append(newTags, existingTag)
		}
	}
	t.tags = newTags
	t.SetAsModified()
}

func (t *Todo) HasTag(tag string) bool {
	for _, existingTag := range t.tags {
		if existingTag == tag {
			return true
		}
	}
	return false
}
