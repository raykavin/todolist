package repository

import (
	"context"
	"time"
	"todolist/internal/domain/shared"
	"todolist/internal/domain/user/entity"
	"todolist/internal/domain/user/repository"
	vo "todolist/internal/domain/user/valueobject"
	"todolist/internal/infrastructure/database"
	"todolist/internal/infrastructure/database/mapper"
	"todolist/internal/infrastructure/database/model"

	"gorm.io/gorm"
)

// userQueryRepository implements repository.UserQueryRepository
type userQueryRepository struct {
	db     *gorm.DB
	mapper *mapper.UserMapper
}

// NewUserQueryRepository creates a new user query repository
func NewUserQueryRepository(db *gorm.DB) repository.UserQueryRepository {
	return &userQueryRepository{
		db:     db,
		mapper: mapper.NewUserMapper(),
	}
}

// FindAll finds all users with pagination
func (r *userQueryRepository) FindAll(ctx context.Context, options shared.QueryOptions) ([]*entity.User, error) {
	users := []*model.User{}

	query := r.db.WithContext(ctx).Model(&model.User{}).Preload("Person")
	query = database.ApplyQueryOptions(query, options)

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	return r.mapper.ToDomainList(users)
}

// FindByStatus finds users by status
func (r *userQueryRepository) FindByStatus(ctx context.Context, status vo.UserStatus, options shared.QueryOptions) ([]*entity.User, error) {
	users := []*model.User{}

	query := r.db.
		WithContext(ctx).
		Model(&model.User{}).
		Preload("Person").
		Where("status = ?", string(status))

	query = database.ApplyQueryOptions(query, options)

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	return r.mapper.ToDomainList(users)
}

// FindByRole finds users by role
func (r *userQueryRepository) FindByRole(ctx context.Context, role vo.UserRole, options shared.QueryOptions) ([]*entity.User, error) {
	people := []*model.User{}

	query := r.db.
		WithContext(ctx).
		Model(&model.User{}).
		Preload("Person").
		Where("role = ?", string(role))

	query = database.ApplyQueryOptions(query, options)
	if err := query.Find(&people).Error; err != nil {
		return nil, err
	}

	return r.mapper.ToDomainList(people)
}

// FindActiveUsers finds all active users
func (r *userQueryRepository) FindActiveUsers(ctx context.Context, options shared.QueryOptions) ([]*entity.User, error) {
	return r.FindByStatus(ctx, vo.StatusActive, options)
}

// FindInactiveUsers finds users who haven't logged in for specified days
func (r *userQueryRepository) FindInactiveUsers(ctx context.Context, days int, options shared.QueryOptions) ([]*entity.User, error) {
	users := []*model.User{}

	threshold := time.Now().AddDate(0, 0, -days)

	query := r.db.
		WithContext(ctx).
		Model(&model.User{}).
		Preload("Person").
		Where("last_login_at < ? OR last_login_at IS NULL", threshold).
		Where("status = ?", string(vo.StatusActive))

	query = database.ApplyQueryOptions(query, options)

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	return r.mapper.ToDomainList(users)
}

// FindBlockedUsers finds all blocked users
func (r *userQueryRepository) FindBlockedUsers(ctx context.Context, options shared.QueryOptions) ([]*entity.User, error) {
	return r.FindByStatus(ctx, vo.StatusBlocked, options)
}

// Count counts users with filters
func (r *userQueryRepository) Count(ctx context.Context, filters []shared.Filter) (int64, error) {
	var count int64

	query := r.db.WithContext(ctx).Model(&model.User{})
	query = database.ApplyFilters(query, filters)

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// CountByStatus counts users grouped by status
func (r *userQueryRepository) CountByStatus(ctx context.Context) (map[vo.UserStatus]int64, error) {
	var results []struct {
		Status string
		Count  int64
	}

	if err := r.db.WithContext(ctx).
		Model(&model.User{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	counts := make(map[vo.UserStatus]int64)
	for _, result := range results {
		status := vo.UserStatus(result.Status)
		if status.IsValid() {
			counts[status] = result.Count
		}
	}

	return counts, nil
}

// CountByRole counts users grouped by role
func (r *userQueryRepository) CountByRole(ctx context.Context) (map[vo.UserRole]int64, error) {
	var results []struct {
		Role  string
		Count int64
	}

	if err := r.db.WithContext(ctx).
		Model(&model.User{}).
		Select("role, COUNT(*) as count").
		Group("role").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	counts := make(map[vo.UserRole]int64)
	for _, result := range results {
		role := vo.UserRole(result.Role)
		if role.IsValid() {
			counts[role] = result.Count
		}
	}

	return counts, nil
}
