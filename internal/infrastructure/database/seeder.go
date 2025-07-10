package database

import (
	"fmt"
	"time"
	"todolist/internal/infrastructure/database/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Seeder handles database seeding
type Seeder struct {
	db *gorm.DB
}

// NewSeeder creates a new seeder instance
func NewSeeder(db *gorm.DB) *Seeder {
	return &Seeder{db: db}
}

// Seed runs all seeders
func (s *Seeder) Seed() error {
	// Seed in order of dependencies
	if err := s.seedTags(); err != nil {
		return fmt.Errorf("failed to seed tags: %w", err)
	}

	if err := s.seedPersonsAndUsers(); err != nil {
		return fmt.Errorf("failed to seed persons and users: %w", err)
	}

	if err := s.seedTodos(); err != nil {
		return fmt.Errorf("failed to seed todos: %w", err)
	}

	return nil
}

// seedTags seeds default tags
func (s *Seeder) seedTags() error {
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
		if err := s.db.FirstOrCreate(&tag, model.Tag{Name: tag.Name}).Error; err != nil {
			return err
		}
	}

	return nil
}

// seedPersonsAndUsers seeds test persons and users
func (s *Seeder) seedPersonsAndUsers() error {
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
		if err := s.db.Where("email = ?", data.person.Email).First(&existingPerson).Error; err == nil {
			continue // Person already exists
		}

		if err := s.db.Create(&data.person).Error; err != nil {
			return err
		}

		// Create user
		data.user.PersonID = data.person.ID
		data.user.PasswordHash = string(hashedPassword)

		if err := s.db.Create(&data.user).Error; err != nil {
			return err
		}
	}

	return nil
}

// seedTodos seeds sample todos
func (s *Seeder) seedTodos() error {
	// Get users
	var users []model.User
	if err := s.db.Find(&users).Error; err != nil {
		return err
	}

	if len(users) == 0 {
		return nil // No users to create todos for
	}

	// Get tags
	var tags []model.Tag
	if err := s.db.Find(&tags).Error; err != nil {
		return err
	}

	// Sample todos
	todoTemplates := []struct {
		title       string
		description string
		status      string
		priority    int8
		daysOffset  int
		tagNames    []string
	}{
		{
			title:       "Complete project documentation",
			description: "Write comprehensive documentation for the new feature",
			status:      "in_progress",
			priority:    3, // High
			daysOffset:  3,
			tagNames:    []string{"work", "urgent"},
		},
		{
			title:       "Buy groceries",
			description: "Milk, bread, eggs, vegetables",
			status:      "pending",
			priority:    2, // Medium
			daysOffset:  1,
			tagNames:    []string{"personal", "shopping"},
		},
		{
			title:       "Schedule dentist appointment",
			description: "Annual checkup",
			status:      "pending",
			priority:    2,
			daysOffset:  7,
			tagNames:    []string{"personal", "health"},
		},
		{
			title:       "Review Q4 budget",
			description: "Analyze expenses and prepare report",
			status:      "pending",
			priority:    3,
			daysOffset:  5,
			tagNames:    []string{"work", "finance"},
		},
		{
			title:       "Learn Go generics",
			description: "Study the new generics feature in Go",
			status:      "completed",
			priority:    1, // Low
			daysOffset:  -2,
			tagNames:    []string{"personal", "learning"},
		},
		{
			title:       "Fix kitchen sink",
			description: "Call plumber or try DIY fix",
			status:      "cancelled",
			priority:    2,
			daysOffset:  -5,
			tagNames:    []string{"home"},
		},
	}

	// Create todos for each user
	for _, user := range users {
		for _, template := range todoTemplates {
			todo := model.Todo{
				ID:          time.Now().Unix(),
				UserID:      user.ID,
				Title:       template.title,
				Description: template.description,
				Status:      template.status,
				Priority:    template.priority,
			}

			// Set due date
			if template.daysOffset != 0 {
				dueDate := time.Now().AddDate(0, 0, template.daysOffset)
				todo.DueDate = &dueDate
			}

			// Set completed date for completed todos
			if template.status == "completed" {
				completedAt := time.Now().AddDate(0, 0, -1)
				todo.CompletedAt = &completedAt
			}

			// Create todo
			if err := s.db.Create(&todo).Error; err != nil {
				return err
			}

			// Assign tags
			for _, tagName := range template.tagNames {
				for _, tag := range tags {
					if tag.Name == tagName {
						if err := s.db.Create(&model.TodoTag{
							TodoID: todo.ID,
							TagID:  tag.ID,
						}).Error; err != nil {
							return err
						}
						break
					}
				}
			}
		}
	}

	return nil
}

// Clean removes all seeded data (useful for testing)
func (s *Seeder) Clean() error {
	// Delete in reverse order of dependencies
	tables := []string{
		"todo_tags",
		"todos",
		"login_attempts",
		"audit_logs",
		"users",
		"persons",
		"tags",
		"todo_daily_statistics",
	}

	for _, table := range tables {
		if err := s.db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)).Error; err != nil {
			return err
		}
	}

	return nil
}
