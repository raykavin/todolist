package valueobject

import (
	"errors"
	"strings"
)

var (
	ErrTitleEmpty    = errors.New("title cannot be empty")
	ErrTitleTooShort = errors.New("title must be at least 3 characters long")
	ErrTitleTooLong  = errors.New("title must not exceed 200 characters")
)

// TodoTitle represents the title of a todo item
type TodoTitle struct {
	value string
}

// NewTodoTitle creates a new TodoTitle with validation
func NewTodoTitle(title string) (TodoTitle, error) {
	title = strings.TrimSpace(title)

	if title == "" {
		return TodoTitle{}, ErrTitleEmpty
	}

	if len(title) < 3 {
		return TodoTitle{}, ErrTitleTooShort
	}

	if len(title) > 200 {
		return TodoTitle{}, ErrTitleTooLong
	}

	return TodoTitle{value: title}, nil
}

// Value returns the title value
func (t TodoTitle) Value() string { return t.value }

// String returns the string representation
func (t TodoTitle) String() string { return t.value }

// Equals compares two titles
func (t TodoTitle) Equals(other TodoTitle) bool { return t.value == other.value }
