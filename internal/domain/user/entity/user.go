package entity

import (
	"errors"
	"todolist/internal/domain/shared"
	"todolist/internal/domain/user/valueobject"
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
		role:     valueobject.RoleUser,
	}, nil
}

// Getters

// PersonID returns the ID of the person associated with the user
func (u User) PersonID() string { return u.personID }

// Username returns the username of the user
func (u User) Username() string { return u.username }

// Password returns the password of the user
func (u User) Password() valueobject.Password { return u.password }

// Status returns the current status of the user
func (u User) Status() valueobject.UserStatus { return u.status }

// Role returns the role of the user
func (u User) Role() valueobject.UserRole { return u.role }

// IsActive checks if the user is active
func (u User) IsActive() bool { return u.status == valueobject.StatusActive }

// Business methods
func (u User) CanPerformAction() error {
	if !u.IsActive() {
		return ErrUserNotActive
	}
	return nil
}

// Update methods

// ChangeUsername changes the username of the user
func (u *User) ChangePassword(newPassword valueobject.Password) {
	u.password = newPassword
	u.SetAsModified()
}

// ChangeUsername changes the username of the user
func (u *User) ChangeRole(newRole valueobject.UserRole) {
	u.role = newRole
	u.SetAsModified()
}

// ChangeUsername changes the username of the user
func (u *User) Activate() {
	u.status = valueobject.StatusActive
	u.SetAsModified()
}

// Deactivate sets the user status to inactive
func (u *User) Deactivate() {
	u.status = valueobject.StatusInactive
	u.SetAsModified()
}

// Block sets the user status to blocked
func (u *User) Block() {
	u.status = valueobject.StatusBlocked
	u.SetAsModified()
}
