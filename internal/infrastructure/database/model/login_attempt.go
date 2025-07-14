package model

import "time"

// LoginAttempt is the login attempts table for security monitoring
type LoginAttempt struct {
	ID         int64     `gorm:"primaryKey;autoIncrement"`
	CreatedAt  time.Time `gorm:"not null;index"`
	Username   string    `gorm:"type:varchar(50);not null;index"`
	UserID     *string   `gorm:"type:uuid;index"`
	Success    bool      `gorm:"not null;index"`
	IPAddress  string    `gorm:"type:varchar(45);not null;index"`
	UserAgent  string    `gorm:"type:varchar(255)"`
	FailReason string    `gorm:"type:varchar(255)"`

	// Relationships
	User *User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

func (LoginAttempt) TableName() string {
	return "login_attempts"
}
