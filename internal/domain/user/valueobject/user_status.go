package valueobject

// UserStatus represents the status of a user
type UserStatus string

const (
	StatusActive   UserStatus = "active"
	StatusInactive UserStatus = "inactive"
	StatusBlocked  UserStatus = "blocked"
	StatusPending  UserStatus = "pending"
)

// IsValid validates if the status is valid
func (s UserStatus) IsValid() bool {
	switch s {
	case StatusActive, StatusInactive, StatusBlocked, StatusPending:
		return true
	default:
		return false
	}
}

// String returns the string representation
func (s UserStatus) String() string {
	return string(s)
}
