package mapper

import (
	"todolist/internal/domain/user/entity"
	vo "todolist/internal/domain/user/valueobject"
	"todolist/internal/infrastructure/database/model"
)

// UserMapper handles conversion between domain entity and database model
type UserMapper struct{}

// NewUserMapper creates a new UserMapper
func NewUserMapper() *UserMapper {
	return &UserMapper{}
}

// ToModel converts domain entity to database model
func (m *UserMapper) ToModel(user *entity.User) *model.User {
	return &model.User{
		ID:           user.ID(),
		PersonID:     user.PersonID(),
		Username:     user.Username(),
		PasswordHash: user.Password().Hash(),
		Status:       string(user.Status()),
		Role:         string(user.Role()),
		CreatedAt:    user.CreatedAt(),
		UpdatedAt:    user.UpdatedAt(),
	}
}

// ToDomain converts database model to domain entity
func (m *UserMapper) ToDomain(model *model.User) (*entity.User, error) {
	password := vo.NewPasswordFromHash(model.PasswordHash)

	status := vo.UserStatus(model.Status)
	if !status.IsValid() {
		status = vo.StatusActive
	}

	role := vo.UserRole(model.Role)
	if !role.IsValid() {
		role = vo.RoleUser
	}

	user, err := entity.NewUser(
		model.ID,
		model.PersonID,
		model.Username,
		password,
		role,
	)
	if err != nil {
		return nil, err
	}

	// Set status and timestamps from database
	if model.Status != string(vo.StatusActive) {
		switch vo.UserStatus(model.Status) {
		case vo.StatusInactive:
			user.Deactivate()
		case vo.StatusBlocked:
			user.Block()
		}
	}

	user.Entity.SetCreatedAt(model.CreatedAt)
	user.Entity.SetUpdatedAt(model.UpdatedAt)

	return user, nil
}

// ToDomainList converts a list of model to domain entities
func (m *UserMapper) ToDomainList(model []*model.User) ([]*entity.User, error) {
	users := make([]*entity.User, 0, len(model))

	for _, model := range model {
		user, err := m.ToDomain(model)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
