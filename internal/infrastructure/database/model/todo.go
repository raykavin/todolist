package model

import (
	"time"

	"gorm.io/gorm"
)

// Todo is the todo table
type Todo struct {
	ID          int64          `gorm:"column:id;primaryKey"`
	CreatedAt   time.Time      `gorm:"column:created_at;not null"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;not null"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index"`
	UserID      int64          `gorm:"column:user_id;not null;index"`
	Title       string         `gorm:"column:title;type:varchar(200);not null"`
	Description string         `gorm:"column:description;type:text"`
	Status      string         `gorm:"column:status;type:varchar(20);not null;default:'pending';index"`
	Priority    int8           `gorm:"column:priority;type:int8;not null;default:1;index"`
	DueDate     *time.Time     `gorm:"column:due_date;type:timestamp;index"`
	CompletedAt *time.Time     `gorm:"column:completed_at;type:timestamp"`

	// Relationships
	User User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Tags []*Tag `gorm:"many2many:todo_tags;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (Todo) TableName() string {
	return "todos"
}

func (t *Todo) BeforeCreate(_ *gorm.DB) error {
	if t.Status == "" {
		t.Status = "pending"
	}
	if t.Priority == 0 {
		t.Priority = 1
	}
	return nil
}
