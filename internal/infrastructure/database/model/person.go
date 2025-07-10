package model

import (
	"time"

	"gorm.io/gorm"
)

// Person is the person table
type Person struct {
	ID        int64          `gorm:"primaryKey"`
	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string         `gorm:"type:varchar(100);not null"`
	Email     string         `gorm:"type:varchar(255);not null;uniqueIndex"`
	Phone     string         `gorm:"type:varchar(20);not null"`
	TxID      string         `gorm:"type:varchar(20);not null"`
	BirthDate *time.Time     `gorm:"type:date"`

	// Relationships
	User *User `gorm:"foreignKey:PersonID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (Person) TableName() string {
	return "people"
}

func (p *Person) BeforeCreate(_ *gorm.DB) error {
	if p.BirthDate == nil {
		return nil
	}

	// Truncate reset the time and keep only the date
	*p.BirthDate = p.BirthDate.Truncate(24 * time.Hour).UTC()
	return nil
}
