package model

import (
	"time"
)

// AuditLog is the audit log table for tracking user activities
type AuditLog struct {
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement"`
	UserID     int64     `gorm:"column:user_id;index"`
	EntityID   int64     `gorm:"column:entity_id;not null;index"`
	EntityType string    `gorm:"column:entity_type;type:varchar(50);not null;index"`
	Action     string    `gorm:"column:action;type:varchar(50);not null;index"`
	OldValues  string    `gorm:"column:old_values;type:jsonb"`
	NewValues  string    `gorm:"column:new_values;type:jsonb"`
	IPAddress  string    `gorm:"column:ip_address;type:varchar(45)"`
	UserAgent  string    `gorm:"column:user_agent;type:varchar(255)"`
	CreatedAt  time.Time `gorm:"column:created_at;not null;index"`

	// Relationships
	User *User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

// TableName specifies the table name
func (AuditLog) TableName() string {
	return "audit_logs"
}
