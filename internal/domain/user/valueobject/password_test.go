package valueobject

import (
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestNewPassword(t *testing.T) {
	tests := []struct {
		name          string
		plainPassword string
		wantErr       error
	}{
		{
			name:          "valid password with minimum length",
			plainPassword: "12345678",
			wantErr:       nil,
		},
		{
			name:          "valid password with more than minimum length",
			plainPassword: "thisIsAVerySecurePassword123!",
			wantErr:       nil,
		},
		{
			name:          "password too short - empty",
			plainPassword: "",
			wantErr:       ErrPasswordTooShort,
		},
		{
			name:          "password too short - 7 characters",
			plainPassword: "1234567",
			wantErr:       ErrPasswordTooShort,
		},
		{
			name:          "valid password with special characters",
			plainPassword: "P@ssw0rd!123",
			wantErr:       nil,
		},
		{
			name:          "valid password with unicode characters",
			plainPassword: "пароль123",
			wantErr:       nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			password, err := NewPassword(tt.plainPassword)

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("NewPassword() error = %v, wantErr %v", err, tt.wantErr)
				}
				if password.hash != "" {
					t.Errorf("NewPassword() returned non-empty hash on error")
				}
				return
			}

			if err != nil {
				t.Errorf("NewPassword() unexpected error = %v", err)
				return
			}

			if password.hash == "" {
				t.Errorf("NewPassword() returned empty hash")
			}

			// Verify that the hash is a valid bcrypt hash
			if err := bcrypt.CompareHashAndPassword([]byte(password.hash), []byte(tt.plainPassword)); err != nil {
				t.Errorf("NewPassword() generated invalid hash: %v", err)
			}
		})
	}
}

func TestNewPasswordFromHash(t *testing.T) {
	tests := []struct {
		name string
		hash string
	}{
		{
			name: "valid bcrypt hash",
			hash: "$2a$10$N9qo8uLOickgx2ZMRZoMye",
		},
		{
			name: "empty hash",
			hash: "",
		},
		{
			name: "arbitrary string",
			hash: "not-a-real-hash",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			password := NewPasswordFromHash(tt.hash)

			if password.hash != tt.hash {
				t.Errorf("NewPasswordFromHash() = %v, want %v", password.hash, tt.hash)
			}
		})
	}
}

func TestPassword_Hash(t *testing.T) {
	plainPassword := "testPassword123"
	password, err := NewPassword(plainPassword)
	if err != nil {
		t.Fatalf("Failed to create password: %v", err)
	}

	hash := password.Hash()
	if hash == "" {
		t.Errorf("Hash() returned empty string")
	}

	if hash != password.hash {
		t.Errorf("Hash() = %v, want %v", hash, password.hash)
	}
}

func TestPassword_Matches(t *testing.T) {
	plainPassword := "correctPassword123"
	password, err := NewPassword(plainPassword)
	if err != nil {
		t.Fatalf("Failed to create password: %v", err)
	}

	tests := []struct {
		name          string
		plainPassword string
		want          bool
	}{
		{
			name:          "correct password",
			plainPassword: "correctPassword123",
			want:          true,
		},
		{
			name:          "incorrect password",
			plainPassword: "wrongPassword123",
			want:          false,
		},
		{
			name:          "empty password",
			plainPassword: "",
			want:          false,
		},
		{
			name:          "password with similar characters",
			plainPassword: "correctPassword124",
			want:          false,
		},
		{
			name:          "password with different case",
			plainPassword: "CORRECTPASSWORD123",
			want:          false,
		},
		{
			name:          "password with extra whitespace",
			plainPassword: "correctPassword123 ",
			want:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := password.Matches(tt.plainPassword); got != tt.want {
				t.Errorf("Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPassword_MatchesWithInvalidHash(t *testing.T) {
	// Test with an invalid hash
	password := NewPasswordFromHash("invalid-hash")

	if password.Matches("anyPassword") {
		t.Errorf("Matches() returned true for invalid hash")
	}
}

func TestPasswordHashUniqueness(t *testing.T) {
	plainPassword := "samePassword123"

	// Create multiple passwords from the same plain text
	password1, err1 := NewPassword(plainPassword)
	if err1 != nil {
		t.Fatalf("Failed to create password1: %v", err1)
	}

	password2, err2 := NewPassword(plainPassword)
	if err2 != nil {
		t.Fatalf("Failed to create password2: %v", err2)
	}

	// Hashes should be different due to bcrypt's salt
	if password1.hash == password2.hash {
		t.Errorf("Two passwords created from same plain text have identical hashes")
	}

	// Both should match the original plain password
	if !password1.Matches(plainPassword) {
		t.Errorf("password1 does not match original plain password")
	}

	if !password2.Matches(plainPassword) {
		t.Errorf("password2 does not match original plain password")
	}
}

func TestPasswordRoundTrip(t *testing.T) {
	// Test that we can create a password, get its hash, create a new password from that hash,
	// and it still matches the original plain password
	plainPassword := "roundTripTest123!"

	password1, err := NewPassword(plainPassword)
	if err != nil {
		t.Fatalf("Failed to create password: %v", err)
	}

	hash := password1.Hash()
	password2 := NewPasswordFromHash(hash)

	if !password2.Matches(plainPassword) {
		t.Errorf("Password created from hash does not match original plain password")
	}
}

func TestPasswordWithLongInput(t *testing.T) {
	// Test with a password at the bcrypt limit (72 bytes)
	passwordAt72Bytes := strings.Repeat("a", 72)

	password1, err := NewPassword(passwordAt72Bytes)
	if err != nil {
		t.Fatalf("Failed to create password with 72 bytes: %v", err)
	}

	if !password1.Matches(passwordAt72Bytes) {
		t.Errorf("72-byte password does not match itself")
	}

	// Test with a password exceeding the bcrypt limit
	passwordOver72Bytes := strings.Repeat("a", 73)

	_, err = NewPassword(passwordOver72Bytes)
	if err == nil {
		t.Errorf("Expected error for password exceeding 72 bytes, but got none")
	}

	// The error should mention the bcrypt limitation
	if err != nil && !strings.Contains(err.Error(), "72 bytes") {
		t.Errorf("Expected error mentioning 72 bytes limit, got: %v", err)
	}

	// Test with a much longer password
	veryLongPassword := strings.Repeat("a", 100)

	_, err = NewPassword(veryLongPassword)
	if err == nil {
		t.Errorf("Expected error for 100-byte password, but got none")
	}
}

// Benchmark tests
func BenchmarkNewPassword(b *testing.B) {
	plainPassword := "benchmarkPassword123!"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = NewPassword(plainPassword)
	}
}

func BenchmarkPasswordMatches(b *testing.B) {
	plainPassword := "benchmarkPassword123!"
	password, err := NewPassword(plainPassword)
	if err != nil {
		b.Fatalf("Failed to create password: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = password.Matches(plainPassword)
	}
}

func BenchmarkPasswordMatchesFail(b *testing.B) {
	plainPassword := "benchmarkPassword123!"
	password, err := NewPassword(plainPassword)
	if err != nil {
		b.Fatalf("Failed to create password: %v", err)
	}

	wrongPassword := "wrongPassword123!"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = password.Matches(wrongPassword)
	}
}
