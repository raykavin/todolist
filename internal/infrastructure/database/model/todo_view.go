package model

import "time"

// TodoView is a database view for optimized todo queries
type TodoView struct {
	ID          int64      `gorm:"column:id"`
	UserID      int64      `gorm:"column:user_id"`
	Username    string     `gorm:"column:username"`
	PersonName  string     `gorm:"column:person_name"`
	Title       string     `gorm:"column:title"`
	Description string     `gorm:"column:description"`
	Status      string     `gorm:"column:status"`
	Priority    int        `gorm:"column:priority"`
	DueDate     *time.Time `gorm:"column:due_date"`
	CompletedAt *time.Time `gorm:"column:completed_at"`
	IsOverdue   bool       `gorm:"column:is_overdue"`
	Tags        string     `gorm:"column:tags"`
	CreatedAt   time.Time  `gorm:"column:created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at"`
}

func (TodoView) TableName() string {
	return "todo_view"
}
