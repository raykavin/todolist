package valueobject

import (
	"errors"
	"net/mail"
	"strings"
)

var (
	ErrInvalidEmail = errors.New("invalid email address")
	ErrEmptyEmail   = errors.New("email cannot be empty")
)

// Email represents a validated email address
type Email struct {
	value      string
	localPart  string
	domainPart string
}

// NewEmail creates a new Email with validation
func NewEmail(email string) (Email, error) {
	email = strings.TrimSpace(email)
	if email == "" {
		return Email{}, ErrEmptyEmail
	}

	// Normalize to lowercase
	normalized := strings.ToLower(email)

	address, err := mail.ParseAddress(normalized)
	if err != nil {
		return Email{}, ErrInvalidEmail
	}

	// Extract just the email part (removes display name if present)
	emailOnly := address.Address

	parts := strings.Split(emailOnly, "@")
	if len(parts) != 2 {
		return Email{}, ErrInvalidEmail
	}

	return Email{
		value:      emailOnly,
		localPart:  parts[0],
		domainPart: parts[1],
	}, nil
}

// Value returns the email value
func (e Email) Value() string { return e.value }

// LocalPart returns the local part of the email
func (e Email) LocalPart() string { return e.localPart }

// Domain returns the domain part of the email
func (e Email) Domain() string { return e.domainPart }

// String returns the string representation
func (e Email) String() string { return e.value }

// Equals compares two emails
func (e Email) Equals(other Email) bool { return e.value == other.value }

// IsBusinessEmail checks if the email belongs to a common free email provider
func (e Email) IsBusinessEmail() bool {
	freeProviders := map[string]bool{
		"gmail.com":      true,
		"yahoo.com":      true,
		"hotmail.com":    true,
		"outlook.com":    true,
		"live.com":       true,
		"aol.com":        true,
		"icloud.com":     true,
		"protonmail.com": true,
	}

	return !freeProviders[e.domainPart]
}

// Mask returns a masked version of the email for display purposes
func (e Email) Mask() string {
	if len(e.localPart) <= 2 {
		return "**@" + e.domainPart
	}

	visibleChars := 2
	if len(e.localPart) > 6 {
		visibleChars = 3
	}

	masked := e.localPart[:visibleChars] + strings.Repeat("*", len(e.localPart)-visibleChars)
	return masked + "@" + e.domainPart
}
