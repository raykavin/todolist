package service

/*
 * token_service.go
 *
 * This file defines the TokenService interface and related types for the domain layer.
 * These types represent the core authentication contracts of the application,
 * independent of any specific implementation.
 *
 * The interface supports token generation, validation, refresh flows, and revocation,
 * providing a complete token lifecycle management system.
 */

import (
	"context"
	"time"
)

// TokenType represents the type of authentication token
type TokenType int

const (
	// TypeAccess represents an access token used for API authentication
	TypeAccess TokenType = iota
	// TypeRefresh represents a refresh token used to obtain new access tokens
	TypeRefresh
)

// String returns the string representation of TokenType
func (t TokenType) String() string {
	switch t {
	case TypeAccess:
		return "access"
	case TypeRefresh:
		return "refresh"
	default:
		return "unknown"
	}
}

// TokenMetadata contains additional information about a token
type TokenMetadata struct {
	// TokenID is a unique identifier for the token (useful for revocation)
	TokenID string `json:"token_id"`
	// IssuedAt is the timestamp when the token was created
	IssuedAt time.Time `json:"issued_at"`
	// ExpiresAt is the timestamp when the token expires
	ExpiresAt time.Time `json:"expires_at"`
	// TokenType indicates whether this is an access or refresh token
	TokenType TokenType `json:"token_type"`
}

// IsExpired checks if the token has expired
func (m TokenMetadata) IsExpired() bool {
	return time.Now().After(m.ExpiresAt)
}

// TimeUntilExpiry returns the duration until the token expires
func (m TokenMetadata) TimeUntilExpiry() time.Duration {
	return time.Until(m.ExpiresAt)
}

// AuthTokens represents a complete set of authentication tokens
type AuthTokens struct {
	// AccessToken is the JWT access token string
	AccessToken string `json:"access_token"`
	// RefreshToken is the JWT refresh token string
	RefreshToken string `json:"refresh_token"`
	// AccessMeta contains metadata about the access token
	AccessMeta TokenMetadata `json:"access_meta"`
	// RefreshMeta contains metadata about the refresh token
	RefreshMeta TokenMetadata `json:"refresh_meta"`
}

// ValidationResult contains the result of token validation
type ValidationResult struct {
	// UserID is the identifier of the authenticated user
	UserID int64 `json:"user_id"`
	// IssuerName is the name of the token issuer
	IssuerName string `json:"issuer_name"`
	// Metadata contains additional token information
	Metadata TokenMetadata `json:"metadata"`
	// Claims contains any custom claims from the token
	Claims map[string]interface{} `json:"claims,omitempty"`
}

// TokenService defines the interface for token generation and validation
type TokenService interface {
	// GenerateTokens creates a new set of access and refresh tokens
	// ctx: context for cancellation and timeout control
	// issuerName: identifier for the token issuer (e.g., application name)
	// userID: unique identifier for the user
	// Returns: AuthTokens containing both access and refresh tokens with metadata
	GenerateTokens(ctx context.Context, issuerName string, userID int64) (*AuthTokens, error)

	// RefreshTokens generates new tokens from a valid refresh token
	// ctx: context for cancellation and timeout control
	// refreshToken: the refresh token string to validate and use for renewal
	// Returns: new AuthTokens pair if the refresh token is valid
	RefreshTokens(ctx context.Context, refreshToken string) (*AuthTokens, error)

	// ValidateToken validates any type of token and returns its claims
	// ctx: context for cancellation and timeout control
	// token: the token string to validate (can be access or refresh token)
	// Returns: ValidationResult with user information and metadata if valid
	ValidateToken(ctx context.Context, token string) (*ValidationResult, error)

	// RevokeToken invalidates a specific token (optional but useful for logout)
	// ctx: context for cancellation and timeout control
	// token: the token string to revoke
	// Returns: error if revocation fails
	RevokeToken(ctx context.Context, token string) error
}

// TokenServiceError represents domain-specific errors for token operations
type TokenServiceError struct {
	Code    string
	Message string
	Err     error
}

func (e *TokenServiceError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *TokenServiceError) Unwrap() error {
	return e.Err
}

// Common error codes for token operations
const (
	ErrCodeInvalidToken     = "INVALID_TOKEN"
	ErrCodeExpiredToken     = "EXPIRED_TOKEN"
	ErrCodeRevokedToken     = "REVOKED_TOKEN"
	ErrCodeInvalidTokenType = "INVALID_TOKEN_TYPE"
	ErrCodeTokenGeneration  = "TOKEN_GENERATION_ERROR"
)
