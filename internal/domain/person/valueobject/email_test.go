package valueobject

import (
	"errors"
	"testing"
)

func TestNewEmail(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantErr    error
		wantValue  string
		wantLocal  string
		wantDomain string
	}{
		{
			name:       "valid simple email",
			input:      "user@example.com",
			wantErr:    nil,
			wantValue:  "user@example.com",
			wantLocal:  "user",
			wantDomain: "example.com",
		},
		{
			name:       "valid email with uppercase",
			input:      "User@Example.COM",
			wantErr:    nil,
			wantValue:  "user@example.com",
			wantLocal:  "user",
			wantDomain: "example.com",
		},
		{
			name:       "valid email with dots",
			input:      "first.last@example.com",
			wantErr:    nil,
			wantValue:  "first.last@example.com",
			wantLocal:  "first.last",
			wantDomain: "example.com",
		},
		{
			name:       "valid email with plus",
			input:      "user+tag@example.com",
			wantErr:    nil,
			wantValue:  "user+tag@example.com",
			wantLocal:  "user+tag",
			wantDomain: "example.com",
		},
		{
			name:       "valid email with subdomain",
			input:      "user@mail.example.com",
			wantErr:    nil,
			wantValue:  "user@mail.example.com",
			wantLocal:  "user",
			wantDomain: "mail.example.com",
		},
		{
			name:       "valid email with spaces around",
			input:      "  user@example.com  ",
			wantErr:    nil,
			wantValue:  "user@example.com",
			wantLocal:  "user",
			wantDomain: "example.com",
		},
		{
			name:    "empty email",
			input:   "",
			wantErr: ErrEmptyEmail,
		},
		{
			name:    "only spaces",
			input:   "   ",
			wantErr: ErrEmptyEmail,
		},
		{
			name:    "missing @",
			input:   "userexample.com",
			wantErr: ErrInvalidEmail,
		},
		{
			name:    "multiple @",
			input:   "user@@example.com",
			wantErr: ErrInvalidEmail,
		},
		{
			name:    "missing domain",
			input:   "user@",
			wantErr: ErrInvalidEmail,
		},
		{
			name:    "missing local part",
			input:   "@example.com",
			wantErr: ErrInvalidEmail,
		},
		{
			name:    "invalid characters",
			input:   "user name@example.com",
			wantErr: ErrInvalidEmail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email, err := NewEmail(tt.input)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			if email.Value() != tt.wantValue {
				t.Errorf("Value() = %v, want %v", email.Value(), tt.wantValue)
			}

			if email.LocalPart() != tt.wantLocal {
				t.Errorf("LocalPart() = %v, want %v", email.LocalPart(), tt.wantLocal)
			}

			if email.Domain() != tt.wantDomain {
				t.Errorf("Domain() = %v, want %v", email.Domain(), tt.wantDomain)
			}

			if email.String() != tt.wantValue {
				t.Errorf("String() = %v, want %v", email.String(), tt.wantValue)
			}
		})
	}
}

func TestEmail_Equals(t *testing.T) {
	email1, _ := NewEmail("user@example.com")
	email2, _ := NewEmail("USER@EXAMPLE.COM")
	email3, _ := NewEmail("other@example.com")

	if !email1.Equals(email2) {
		t.Error("Expected email1 and email2 to be equal")
	}

	if email1.Equals(email3) {
		t.Error("Expected email1 and email3 to be different")
	}
}

func TestEmail_IsBusinessEmail(t *testing.T) {
	tests := []struct {
		email      string
		isBusiness bool
	}{
		{"user@company.com", true},
		{"user@gmail.com", false},
		{"user@yahoo.com", false},
		{"user@hotmail.com", false},
		{"user@outlook.com", false},
		{"user@customdomain.org", true},
		{"user@university.edu", true},
		{"user@protonmail.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			email, _ := NewEmail(tt.email)
			if email.IsBusinessEmail() != tt.isBusiness {
				t.Errorf("IsBusinessEmail() = %v, want %v", email.IsBusinessEmail(), tt.isBusiness)
			}
		})
	}
}

func TestEmail_Mask(t *testing.T) {
	tests := []struct {
		email    string
		expected string
	}{
		{"user@example.com", "us**@example.com"},
		{"a@example.com", "**@example.com"},
		{"ab@example.com", "**@example.com"},
		{"abc@example.com", "ab*@example.com"},
		{"longusername@example.com", "lon*********@example.com"},
		{"verylongusername@example.com", "ver*************@example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			email, _ := NewEmail(tt.email)
			masked := email.Mask()
			if masked != tt.expected {
				t.Errorf("Mask() = %v, want %v", masked, tt.expected)
			}
		})
	}
}
