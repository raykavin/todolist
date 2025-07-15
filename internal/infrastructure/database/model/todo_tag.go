package model

import (
	"time"
)

// TodoTag is the junction table for many-to-many relationship
type TodoTag struct {
	TodoID    int64     `gorm:"column:todo_id;primaryKey"`
	TagID     int64     `gorm:"column:tag_id;primaryKey"`
	CreatedAt time.Time `gorm:"column:created_at;not null"`

	// Relationships
	Todo Todo `gorm:"foreignKey:TodoID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Tag  Tag  `gorm:"foreignKey:TagID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (TodoTag) TableName() string {
	return "todo_tags"
}
