package valueobject

import (
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPasswordTooShort = errors.New("password must be at least 8 characters long")
	ErrPasswordTooWeak  = errors.New("password is too weak")
)

// Password represents an encrypted password
type Password struct {
	hash string
}

// NewPassword creates a new Password value object
func NewPassword(plainPassword string) (Password, error) {
	if len(plainPassword) < 8 {
		return Password{}, ErrPasswordTooShort
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return Password{}, err
	}

	return Password{hash: string(hash)}, nil
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return ErrPasswordTooShort
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasDigit = true
		case strings.ContainsRune("!@#$%^&*()_+-=[]{}|;:,.<>?", char):
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasDigit || !hasSpecial {
		return ErrPasswordTooWeak
	}

	return nil
}

// NewPasswordFromHash creates a Password from an existing hash
func NewPasswordFromHash(hash string) Password {
	return Password{hash: hash}
}

// Hash returns the password hash
func (p Password) Hash() string {
	return p.hash
}

// Matches checks if the plain password matches the hashed password
func (p Password) Matches(plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p.hash), []byte(plainPassword))
	return err == nil
}
