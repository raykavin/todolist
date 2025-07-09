package dto

import "time"

// TodoFilter represents filters for todo queries
type TodoFilter struct {
	UserID      int64
	Status      []string
	Priority    []string
	Tags        []string
	IsOverdue   *bool
	DueDateFrom *time.Time
	DueDateTo   *time.Time
	SearchTerm  string
}
