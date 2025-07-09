package valueobject

import (
	"errors"
	"strings"
)

var (
	ErrDescriptionTooLong = errors.New("description must not exceed 1000 characters")
)

// TodoDescription represents the description of a todo item
type TodoDescription struct {
	value string
}

// NewTodoDescription creates a new TodoDescription with validation
func NewTodoDescription(description string) (TodoDescription, error) {
	description = strings.TrimSpace(description)

	if len(description) > 1000 {
		return TodoDescription{}, ErrDescriptionTooLong
	}

	return TodoDescription{value: description}, nil
}

// Value returns the description value
func (d TodoDescription) Value() string { return d.value }

// String returns the string representation
func (d TodoDescription) String() string { return d.value }

// IsEmpty checks if the description is empty
func (d TodoDescription) IsEmpty() bool { return d.value == "" }

// WordCount returns the number of words in the description
func (d TodoDescription) WordCount() int {
	if d.IsEmpty() {
		return 0
	}
	return len(strings.Fields(d.value))
}
