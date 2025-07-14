package model

import "time"

// LoginAttempt is the login attempts table for security monitoring
type LoginAttempt struct {
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement"`
	UserID     *int64    `gorm:"column:user_id;index"`
	CreatedAt  time.Time `gorm:"column:created_at;not null;index"`
	Username   string    `gorm:"column:username;type:varchar(50);not null;index"`
	Success    bool      `gorm:"column:success;not null;index"`
	IPAddress  string    `gorm:"column:ip_address;type:varchar(45);not null;index"`
	UserAgent  string    `gorm:"column:user_agent;type:varchar(255)"`
	FailReason string    `gorm:"column:fail_reason;type:varchar(255)"`

	// Relationships
	User *User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

func (LoginAttempt) TableName() string {
	return "login_attempts"
}
