package repository

import (
	"context"
	"errors"
	"todolist/internal/domain/shared"
	"todolist/internal/domain/user/entity"
	"todolist/internal/domain/user/repository"
	"todolist/internal/infrastructure/database/mapper"
	"todolist/internal/infrastructure/database/model"

	"gorm.io/gorm"
)

// userRepository implements repository.UserRepository
type userRepository struct {
	db     *gorm.DB
	mapper *mapper.UserMapper
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{
		db:     db,
		mapper: mapper.NewUserMapper(),
	}
}

// Save saves or updates a user
func (r *userRepository) Save(ctx context.Context, user *entity.User) error {
	userModel := r.mapper.ToModel(user)

	err := r.db.WithContext(ctx).Save(userModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return shared.ErrDuplicateEntry
		}
		return err
	}

	return nil
}

// Delete deletes a user (soft delete)
func (r *userRepository) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Delete(&model.User{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return shared.ErrNotFound
	}

	return nil
}

// FindByID finds a user by ID
func (r *userRepository) FindByID(ctx context.Context, id int64) (*entity.User, error) {
	userModel := &model.User{}

	err := r.db.WithContext(ctx).
		Preload("Person").
		First(userModel, "id = ?", id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.ErrNotFound
		}
		return nil, err
	}

	return r.mapper.ToDomain(userModel)
}

// FindByUsername finds a user by username
func (r *userRepository) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	userModel := &model.User{}

	err := r.db.
		WithContext(ctx).
		Preload("Person").
		Where("username = ?", username).
		First(userModel).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.ErrNotFound
		}
		return nil, err
	}

	return r.mapper.ToDomain(userModel)
}

// FindByPersonID finds a user by person ID
func (r *userRepository) FindByPersonID(ctx context.Context, personID int64) (*entity.User, error) {
	userModel := &model.User{}

	err := r.db.WithContext(ctx).
		Preload("Person").
		Where("person_id = ?", personID).
		First(userModel).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.ErrNotFound
		}
		return nil, err
	}

	return r.mapper.ToDomain(userModel)
}

// ExistsByUsername checks if a user exists by username
func (r *userRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64

	if err := r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("username = ?", username).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// ExistsByPersonID checks if a user exists by person ID
func (r *userRepository) ExistsByPersonID(ctx context.Context, personID int64) (bool, error) {
	var count int64

	if err := r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("person_id = ?", personID).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
