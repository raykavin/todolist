package valueobject

// UserRole represents the role of a user
type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
)

// IsValid validates if the role is valid
func (r UserRole) IsValid() bool {
	switch r {
	case RoleUser, RoleAdmin:
		return true
	default:
		return false
	}
}

// String returns the string representation
func (r UserRole) String() string {
	return string(r)
}

// HasPermission checks if the role has a specific permission
func (r UserRole) HasPermission(permission string) bool {
	switch r {
	case RoleAdmin:
		return true
	case RoleUser:
		userPermissions := map[string]bool{
			"todo:read":   true,
			"todo:create": false,
			"todo:update": false,
			"todo:delete": false,
		}
		return userPermissions[permission]
	default:
		return false
	}
}
