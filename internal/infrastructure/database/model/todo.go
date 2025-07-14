package model

import (
	"time"

	"gorm.io/gorm"
)

// Todo is the todo table
type Todo struct {
	ID          int64          `gorm:"primaryKey"`
	CreatedAt   time.Time      `gorm:"not null"`
	UpdatedAt   time.Time      `gorm:"not null"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	UserID      int64          `gorm:"not null;index"`
	Title       string         `gorm:"type:varchar(200);not null"`
	Description string         `gorm:"type:text"`
	Status      string         `gorm:"type:varchar(20);not null;default:'pending';index"`
	Priority    int8           `gorm:"type:int8;not null;default:1;index"`
	DueDate     *time.Time     `gorm:"type:timestamp;index"`
	CompletedAt *time.Time     `gorm:"type:timestamp"`

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
