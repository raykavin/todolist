package entity

import (
	"errors"
	"todolist/internal/domain/shared"
)

var (
	ErrInvalidUsername = errors.New("invalid username")
	ErrUserNotActive   = errors.New("user is not active")
	ErrInvalidPersonID = errors.New("invalid person ID")
)

// User represents a system user
type User struct {
	shared.Entity
	personID string
	username string
	password valueobject.Password
	status   valueobject.UserStatus
	role     valueobject.UserRole
}

// NewUser creates a new User entity
func NewUser(
	id uint64,
	personID string,
	username string,
	password valueobject.Password,
	role valueobject.UserRole,
) (*User, error) {
	if personID == "" {
		return nil, ErrInvalidPersonID
	}

	if username == "" {
		return nil, ErrInvalidUsername
	}

	return &User{
		Entity:   shared.NewEntity(id),
		personID: personID,
		username: username,
		password: password,
		status:   valueobject.StatusActive,
		role:     role,
	}, nil
}

// Getters
func (u *User) PersonID() string {
	return u.personID
}

func (u *User) Username() string {
	return u.username
}

func (u *User) Password() valueobject.Password {
	return u.password
}

func (u *User) Status() valueobject.UserStatus {
	return u.status
}

func (u *User) Role() valueobject.UserRole {
	return u.role
}

// Business methods
func (u *User) IsActive() bool {
	return u.status == valueobject.StatusActive
}

func (u *User) CanPerformAction() error {
	if !u.IsActive() {
		return ErrUserNotActive
	}
	return nil
}

// Update methods
func (u *User) ChangePassword(newPassword valueobject.Password) {
	u.password = newPassword
	u.SetUpdatedAt()
}

func (u *User) ChangeRole(newRole valueobject.UserRole) {
	u.role = newRole
	u.SetUpdatedAt()
}

func (u *User) Activate() {
	u.status = valueobject.StatusActive
	u.SetUpdatedAt()
}

func (u *User) Deactivate() {
	u.status = valueobject.StatusInactive
	u.SetUpdatedAt()
}

func (u *User) Block() {
	u.status = valueobject.StatusBlocked
	u.SetUpdatedAt()
}
