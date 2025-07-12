package model

import (
	"time"

	"gorm.io/gorm"
)

// User is the user table
type User struct {
	ID           int64          `gorm:"type:uuid;primaryKey"`
	CreatedAt    time.Time      `gorm:"not null"`
	UpdatedAt    time.Time      `gorm:"not null"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	PersonID     int64          `gorm:"type:uuid;not null;uniqueIndex"`
	Username     string         `gorm:"type:varchar(50);not null;uniqueIndex"`
	PasswordHash string         `gorm:"type:varchar(255);not null"`
	Status       string         `gorm:"type:varchar(20);not null;default:'active'"`
	Role         string         `gorm:"type:varchar(20);not null;default:'user'"`
	LastLoginAt  *time.Time     `gorm:"type:timestamp"`

	// Relationships
	Person Person  `gorm:"foreignKey:PersonID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Todos  []*Todo `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	// Additional fields for security
	FailedLoginAttempts int        `gorm:"default:0"`
	LockedUntil         *time.Time `gorm:"type:timestamp"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Status == "" {
		u.Status = "active"
	}
	if u.Role == "" {
		u.Role = "user"
	}
	return nil
}
