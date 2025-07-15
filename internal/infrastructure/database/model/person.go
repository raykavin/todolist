package model

import (
	"time"

	"gorm.io/gorm"
)

// Person is the person table
type Person struct {
	ID        int64          `gorm:"column:id;primaryKey"`
	CreatedAt time.Time      `gorm:"column:created_at;not null"`
	UpdatedAt time.Time      `gorm:"column:updated_at;not null"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
	Name      string         `gorm:"column:name;type:varchar(100);not null"`
	Email     string         `gorm:"column:email;type:varchar(255);not null;uniqueIndex"`
	Phone     string         `gorm:"column:phone;type:varchar(20);not null"`
	TaxID     string         `gorm:"column:tax_id;type:varchar(20);uniqueIndex;not null"`
	BirthDate *time.Time     `gorm:"column:birth_date;type:date"`

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
