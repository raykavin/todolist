package model

import "time"

// UserStatisticsView is a database view for user statistics
type UserStatisticsView struct {
	UserID          int64     `gorm:"column:user_id"`
	Username        string    `gorm:"column:username"`
	PersonName      string    `gorm:"column:person_name"`
	TotalTodos      int64     `gorm:"column:total_todos"`
	CompletedTodos  int64     `gorm:"column:completed_todos"`
	PendingTodos    int64     `gorm:"column:pending_todos"`
	InProgressTodos int64     `gorm:"column:in_progress_todos"`
	CancelledTodos  int64     `gorm:"column:cancelled_todos"`
	OverdueTodos    int64     `gorm:"column:overdue_todos"`
	CompletionRate  float64   `gorm:"column:completion_rate"`
	LastActivityAt  time.Time `gorm:"column:last_activity_at"`
}

// TableName specifies the view name
func (UserStatisticsView) TableName() string {
	return "user_statistics_view"
}
