package gorm

import (
	"context"
	"ecommerce/internal/domain/shared/valueobjects"
	"ecommerce/internal/domain/user/entity"
	"ecommerce/internal/domain/user/repository"
	"ecommerce/internal/infrastructure/model"
	"errors"
	"strconv"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new GORM user repository.
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	userModel := model.FromDomain(user)
	return r.db.WithContext(ctx).Create(userModel).Error
}

func (r *userRepository) GetByID(ctx context.Context, id valueobjects.ID) (*entity.User, error) {
	var userModel model.User
	// This conversion is risky and assumes the ID can be represented as a uint.
	// A better approach would be to use the string representation of the ID if the DB supports it.
	uid, err := strconv.ParseUint(id.String(), 10, 64)
	if err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).First(&userModel, uid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entity.ErrUserNotFound
		}
		return nil, err
	}
	return userModel.ToDomain()
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var userModel model.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&userModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entity.ErrUserNotFound
		}
		return nil, err
	}
	return userModel.ToDomain()
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	userModel := model.FromDomain(user)
	return r.db.WithContext(ctx).Save(userModel).Error
}

func (r *userRepository) Delete(ctx context.Context, id valueobjects.ID) error {
	uid, err := strconv.ParseUint(id.String(), 10, 64)
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).Delete(&model.User{}, uid).Error
}
