package model

import (
	"time"

	"gorm.io/gorm"
)

// Tag is the tag table
type Tag struct {
	ID        int64          `gorm:"primaryKey;autoIncrement"`
	Name      string         `gorm:"type:varchar(50);not null;uniqueIndex"`
	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Relationships
	Todos []*Todo `gorm:"many2many:todo_tags;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (Tag) TableName() string {
	return "tags"
}
