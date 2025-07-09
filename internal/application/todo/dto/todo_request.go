package dto

import "time"

// CreateTodoRequest represents the request to create a todo
type CreateTodoRequest struct {
	Title       string     `json:"title" validate:"required,min=3,max=200"`
	Description string     `json:"description,omitempty" validate:"max=1000"`
	Priority    string     `json:"priority" validate:"required,oneof=low medium high critical"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	Tags        []string   `json:"tags,omitempty"`
}

// UpdateTodoRequest represents the request to update a todo
type UpdateTodoRequest struct {
	Title       *string    `json:"title,omitempty" validate:"omitempty,min=3,max=200"`
	Description *string    `json:"description,omitempty" validate:"omitempty,max=1000"`
	Priority    *string    `json:"priority,omitempty" validate:"omitempty,oneof=low medium high critical"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	Status      *string    `json:"status,omitempty" validate:"omitempty,oneof=pending in_progress completed cancelled"`
}
