package model

import (
	"time"

	"gorm.io/gorm"
)

// User is the user table
type User struct {
	ID           int64          `gorm:"column:id;primaryKey"`
	CreatedAt    time.Time      `gorm:"column:created_at;not null"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;not null"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;index"`
	PersonID     int64          `gorm:"column:person_id;not null;uniqueIndex"`
	Username     string         `gorm:"column:username;type:varchar(50);not null;uniqueIndex"`
	PasswordHash string         `gorm:"column:password_hash;type:varchar(255);not null"`
	Status       string         `gorm:"column:status;type:varchar(20);not null;default:'active'"`
	Role         string         `gorm:"column:role;type:varchar(20);not null;default:'user'"`
	LastLoginAt  *time.Time     `gorm:"column:last_login_at;type:timestamp"`

	// Relationships
	Person Person  `gorm:"foreignKey:PersonID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Todos  []*Todo `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	// Additional fields for security
	FailedLoginAttempts int        `gorm:"column:failed_login_attempts;default:0"`
	LockedUntil         *time.Time `gorm:"column:locked_until;type:timestamp"`
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
