package model

import (
	"time"

	"gorm.io/gorm"
)

// Tag is the tag table
type Tag struct {
	ID        int64          `gorm:"column:id;primaryKey;autoIncrement"`
	Name      string         `gorm:"column:name;type:varchar(50);not null;uniqueIndex"`
	CreatedAt time.Time      `gorm:"column:created_at;not null"`
	UpdatedAt time.Time      `gorm:"column:updated_at;not null"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`

	// Relationships
	Todos []*Todo `gorm:"many2many:todo_tags;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (Tag) TableName() string {
	return "tags"
}
