package database

/*
 * seeder.go
 *
 * This file provides database seeding logic.
 *
 * Use it to populate the database with initial or sample data,
 * which can help in development, testing, or demo environments.
 *
 * Examples include inserting default roles, admin accounts, or static lookup tables.
 */

import (
	"fmt"
	"time"
	"todolist/internal/infrastructure/database/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Seed runs all seeders
func Seed(db *gorm.DB) error {
	// Seed in order of dependencies
	if err := seedTags(db); err != nil {
		return fmt.Errorf("failed to seed tags: %w", err)
	}

	if err := seedPersonsAndUsers(db); err != nil {
		return fmt.Errorf("failed to seed persons and users: %w", err)
	}

	return nil
}

// seedTags seeds default tags
func seedTags(db *gorm.DB) error {
	tags := []model.Tag{
		{Name: "work"},
		{Name: "personal"},
		{Name: "urgent"},
		{Name: "shopping"},
		{Name: "health"},
		{Name: "finance"},
		{Name: "learning"},
		{Name: "home"},
		{Name: "family"},
		{Name: "project"},
	}

	for _, tag := range tags {
		if err := db.FirstOrCreate(&tag, model.Tag{Name: tag.Name}).Error; err != nil {
			return err
		}
	}

	return nil
}

// seedPersonsAndUsers seeds test persons and users
func seedPersonsAndUsers(db *gorm.DB) error {
	testData := []struct {
		person model.Person
		user   model.User
	}{
		{
			person: model.Person{
				ID:    time.Now().Unix(),
				Name:  "Administrador",
				Email: "soge@fibralink.net.br",
				Phone: "94999990001",
			},
			user: model.User{
				ID:       time.Now().Unix(),
				Username: "admin",
				Role:     "admin",
				Status:   "active",
			},
		},
	}

	// Hash password for all users
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	for _, data := range testData {
		// Create person
		var existingPerson model.Person
		if err := db.Where("email = ?", data.person.Email).First(&existingPerson).Error; err == nil {
			continue // Person already exists
		}

		if err := db.Create(&data.person).Error; err != nil {
			return err
		}

		// Create user
		data.user.PersonID = data.person.ID
		data.user.PasswordHash = string(hashedPassword)

		if err := db.Create(&data.user).Error; err != nil {
			return err
		}
	}

	return nil
}
