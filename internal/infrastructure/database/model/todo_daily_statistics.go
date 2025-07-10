package model

import "time"

// TodoDailyStatistics is daily todo statistics
type TodoDailyStatistics struct {
	Date                  time.Time `gorm:"column:date;primaryKey"`
	UserID                int64     `gorm:"column:user_id;primaryKey"`
	Created               int64     `gorm:"column:created"`
	Completed             int64     `gorm:"column:completed"`
	Cancelled             int64     `gorm:"column:cancelled"`
	AverageTimeToComplete float64   `gorm:"column:avg_time_to_complete"`
}

func (TodoDailyStatistics) TableName() string {
	return "todo_daily_statistics"
}
