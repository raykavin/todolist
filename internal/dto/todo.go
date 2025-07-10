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

// TodoResponse represents a todo in API responses
type TodoResponse struct {
	ID          int64      `json:"id"`
	UserID      int64      `json:"user_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Tags        []string   `json:"tags"`
	IsOverdue   bool       `json:"is_overdue"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// TodoListResponse represents a list of todos
type TodoListResponse struct {
	Todos      []*TodoResponse `json:"todos"`
	TotalCount int64           `json:"total_count"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
}
