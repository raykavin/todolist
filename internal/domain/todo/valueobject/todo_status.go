package valueobject

import "errors"

// TodoStatus represents the status of a todo item
type TodoStatus string

const (
	StatusPending    TodoStatus = "pending"
	StatusInProgress TodoStatus = "in_progress"
	StatusCompleted  TodoStatus = "completed"
	StatusCancelled  TodoStatus = "cancelled"
)

var (
	ErrInvalidStatus = errors.New("invalid todo status")
)

// IsValid validates if the status is valid
func (s TodoStatus) IsValid() bool {
	switch s {
	case StatusPending, StatusInProgress, StatusCompleted, StatusCancelled:
		return true
	default:
		return false
	}
}

// CanTransitionTo checks if status can transition to another status
func (s TodoStatus) CanTransitionTo(newStatus TodoStatus) bool {
	switch s {
	case StatusPending:
		return newStatus == StatusInProgress || newStatus == StatusCancelled || newStatus == StatusCompleted
	case StatusInProgress:
		return newStatus == StatusCompleted || newStatus == StatusCancelled || newStatus == StatusPending
	case StatusCompleted:
		return newStatus == StatusPending // Allow reopening
	case StatusCancelled:
		return newStatus == StatusPending // Allow reactivating
	default:
		return false
	}
}

// String returns the string representation
func (s TodoStatus) String() string { return string(s) }

// IsFinal checks if the status is a final state
func (s TodoStatus) IsFinal() bool { return s == StatusCompleted || s == StatusCancelled }

// NewTodoStatusFromString creates a TodoStatus from string
func NewTodoStatusFromString(status string) (TodoStatus, error) {
	s := TodoStatus(status)
	if !s.IsValid() {
		return "", ErrInvalidStatus
	}
	return s, nil
}
