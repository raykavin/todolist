package model

import (
	"time"
)

// AuditLog is the audit log table for tracking user activities
type AuditLog struct {
	ID         int64     `gorm:"primaryKey;autoIncrement"`
	UserID     string    `gorm:"type:uuid;index"`
	EntityType string    `gorm:"type:varchar(50);not null;index"`
	EntityID   string    `gorm:"type:varchar(255);not null;index"`
	Action     string    `gorm:"type:varchar(50);not null;index"`
	OldValues  string    `gorm:"type:jsonb"`
	NewValues  string    `gorm:"type:jsonb"`
	IPAddress  string    `gorm:"type:varchar(45)"`
	UserAgent  string    `gorm:"type:varchar(255)"`
	CreatedAt  time.Time `gorm:"not null;index"`

	// Relationships
	User *User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

// TableName specifies the table name
func (AuditLog) TableName() string {
	return "audit_logs"
}
