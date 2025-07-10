package entity

import (
	"errors"
	"time"
	"todolist/internal/domain/shared"
	vo "todolist/internal/domain/user/valueobject"
)

var (
	ErrInvalidUsername = errors.New("invalid username")
	ErrUserNotActive   = errors.New("user is not active")
	ErrInvalidPersonID = errors.New("invalid person ID")
)

// User represents a system user
type User struct {
	shared.Entity
	personID           int64
	username           string
	loginAttempts      int
	lastLoginAttemptAt time.Time
	password           vo.Password
	status             vo.UserStatus
	role               vo.UserRole
}

// NewUser creates a new User entity
func NewUser(
	id int64,
	personID int64,
	username string,
	password vo.Password,
) (*User, error) {
	if personID == 0 {
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
		status:   vo.StatusActive,
		role:     vo.RoleUser,
	}, nil
}

// Getters

// PersonID returns the ID of the person associated with the user
func (u User) PersonID() int64 { return u.personID }

// Username returns the username of the user
func (u User) Username() string { return u.username }

// Password returns the password of the user
func (u User) Password() vo.Password { return u.password }

// Status returns the current status of the user
func (u User) Status() vo.UserStatus { return u.status }

// Role returns the role of the user
func (u User) Role() vo.UserRole { return u.role }

// IsActive checks if the user is active
func (u User) IsActive() bool { return u.status == vo.StatusActive }

// LastLoginAttemptAt returns the last login attemps of user
func (u User) LastLoginAttemptAt() time.Time { return u.lastLoginAttemptAt }

// FailedLoginAttempts returns the counter of login attempts
func (u User) FailedLoginAttempts() int { return u.loginAttempts }

// Business methods

// CanPerformAction check if user is allowed to perform action
func (u User) CanPerformAction() error {
	if !u.IsActive() {
		return ErrUserNotActive
	}
	return nil
}

// Update methods

// ChangeUsername changes the username of the user
func (u *User) ChangePassword(newPassword vo.Password) {
	u.password = newPassword
	u.SetAsModified()
}

// ChangeUsername changes the username of the user
func (u *User) ChangeRole(newRole vo.UserRole) {
	u.role = newRole
	u.SetAsModified()
}

// ChangeUsername changes the username of the user
func (u *User) Activate() {
	u.status = vo.StatusActive
	u.SetAsModified()
}

// Deactivate sets the user status to inactive
func (u *User) Deactivate() {
	u.status = vo.StatusInactive
	u.SetAsModified()
}

// Block sets the user status to blocked
func (u *User) Block() {
	u.status = vo.StatusBlocked
	u.SetAsModified()
}

// IncrementLoginAttempts increment login attempts
func (u *User) IncrementLoginAttempts() {
	u.loginAttempts++
	u.lastLoginAttemptAt = time.Now()
	u.SetAsModified()
}
