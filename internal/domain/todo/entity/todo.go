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
	ErrTodoNotFound            = errors.New("todo not found")
	ErrUnauthorizedTodoAccess  = errors.New("unauthorized to access this todo")
)

// Todo represents a todo item
type Todo struct {
	shared.Entity
	userID      int64
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
	id int64,
	userID int64,
	title vo.TodoTitle,
	description vo.TodoDescription,
	priority sharedvo.Priority,
	dueDate *time.Time,
) (*Todo, error) {
	if userID == 0 {
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

// ID returns the todo's ID
func (t *Todo) UserID() int64 { return t.userID }

// UserID returns the todo's user ID
func (t *Todo) Title() vo.TodoTitle { return t.title }

// Description returns the todo's description
func (t *Todo) Description() vo.TodoDescription { return t.description }

// Status returns the todo's status
func (t *Todo) Status() vo.TodoStatus { return t.status }

// Priority returns the todo's priority
func (t *Todo) Priority() sharedvo.Priority { return t.priority }

// DueDate returns a copy of the todo's due date
func (t *Todo) DueDate() *time.Time {
	if t.dueDate == nil {
		return nil
	}
	dueDateCopy := *t.dueDate
	return &dueDateCopy
}

// CompletedAt returns a copy of the todo's completed timestamp
func (t *Todo) CompletedAt() *time.Time {
	if t.completedAt == nil {
		return nil
	}
	completedAtCopy := *t.completedAt
	return &completedAtCopy
}

// Tags returns a copy of the todo's tags
func (t *Todo) Tags() []string {
	// Return a copy to prevent external modification
	tagsCopy := make([]string, len(t.tags))
	copy(tagsCopy, t.tags)
	return tagsCopy
}

// Business methods

// IsCompleted checks if the todo is completed
func (t *Todo) IsCompleted() bool {
	return t.status == vo.StatusCompleted
}

// IsPending checks if the todo is pending
func (t *Todo) IsOverdue() bool {
	if t.dueDate == nil || t.IsCompleted() {
		return false
	}
	return t.dueDate.Before(time.Now())
}

// IsOverdue checks if the todo is overdue
func (t *Todo) DaysUntilDue() *int {
	if t.dueDate == nil {
		return nil
	}
	days := int(time.Until(*t.dueDate).Hours() / 24)
	return &days
}

// Update methods

// UpdateTitle updates the todo's title
func (t *Todo) UpdateTitle(title vo.TodoTitle) {
	t.title = title
	t.SetAsModified()
}

// UpdateDescription updates the todo's description
func (t *Todo) UpdateDescription(description vo.TodoDescription) {
	t.description = description
	t.SetAsModified()
}

// UpdatePriority updates the todo's priority
func (t *Todo) UpdatePriority(priority sharedvo.Priority) {
	t.priority = priority
	t.SetAsModified()
}

// UpdateDueDate updates the todo's due date
func (t *Todo) UpdateDueDate(dueDate *time.Time) error {
	if dueDate != nil && dueDate.Before(time.Now()) {
		return ErrInvalidDueDate
	}
	t.dueDate = dueDate
	t.SetAsModified()
	return nil
}

// ChangeStatus changes the todo's status with validation
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

// Complete marks the todo as completed
func (t *Todo) Complete() error {
	if t.IsCompleted() {
		return ErrTodoAlreadyCompleted
	}
	return t.ChangeStatus(vo.StatusCompleted)
}

// Progress management

// StartProgress marks the todo as in progress
func (t *Todo) StartProgress() error {
	return t.ChangeStatus(vo.StatusInProgress)
}

// Cancel marks the todo as cancelled
func (t *Todo) Cancel() error {
	return t.ChangeStatus(vo.StatusCancelled)
}

// Reopen marks the todo as pending again
func (t *Todo) Reopen() error {
	return t.ChangeStatus(vo.StatusPending)
}

// Tag management

// AddTag adds a tag to the todo, ensuring no duplicates
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

// RemoveTag removes a tag from the todo
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

// HasTag checks if the todo has a specific tag
func (t *Todo) HasTag(tag string) bool {
	for _, existingTag := range t.tags {
		if existingTag == tag {
			return true
		}
	}
	return false
}
