package repository

import (
	"context"
	"errors"
	"todolist/internal/domain/shared"
	"todolist/internal/domain/todo/entity"
	"todolist/internal/domain/todo/repository"
	"todolist/internal/infrastructure/database/mapper"
	"todolist/internal/infrastructure/database/model"

	"gorm.io/gorm"
)

// todoRepository implements repository.TodoRepository
type todoRepository struct {
	db     *gorm.DB
	mapper *mapper.TodoMapper
}

// NewTodoRepository creates a new todo repository
func NewTodoRepository(db *gorm.DB) repository.TodoRepository {
	return &todoRepository{
		db:     db,
		mapper: mapper.NewTodoMapper(),
	}
}

// Save saves or updates a todo
func (r *todoRepository) Save(ctx context.Context, todo *entity.Todo) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		todoModel := r.mapper.ToModel(todo)

		// Save todo
		if err := tx.Save(todoModel).Error; err != nil {
			return err
		}

		// Handle tags
		// First, delete existing tag associations
		if err := tx.
			Where("todo_id = ?", todoModel.ID).
			Delete(&model.TodoTag{}).
			Error; err != nil {
			return err
		}

		// Then, create new tag associations
		for _, tagName := range todo.Tags() {
			tagModel := &model.Tag{}

			// Find or create tag
			if err := tx.
				Where("name = ?", tagName).
				FirstOrCreate(tagModel, model.Tag{Name: tagName}).
				Error; err != nil {
				return err
			}

			// Create association
			todoTag := model.TodoTag{
				TodoID: tagModel.ID,
				TagID:  tagModel.ID,
			}

			if err := tx.Create(&todoTag).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// Delete deletes a todo (soft delete)
func (r *todoRepository) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Delete(&model.Todo{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return shared.ErrNotFound
	}

	return nil
}

// FindByID finds a todo by ID
func (r *todoRepository) FindByID(ctx context.Context, id int64) (*entity.Todo, error) {
	var model model.Todo

	if err := r.db.WithContext(ctx).
		Preload("Tags").
		First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.ErrNotFound
		}
		return nil, err
	}

	return r.mapper.ToDomain(&model)
}

// FindByUserID finds all todos for a user
func (r *todoRepository) FindByUserID(ctx context.Context, userID int64) ([]*entity.Todo, error) {
	var model []*model.Todo

	if err := r.db.WithContext(ctx).
		Preload("Tags").
		Where("user_id = ?", userID).
		Find(&model).Error; err != nil {
		return nil, err
	}

	return r.mapper.ToDomainList(model)
}

// DeleteByUserID deletes all todos for a user
func (r *todoRepository) DeleteByUserID(ctx context.Context, userID int64) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&model.Todo{}).Error
}
