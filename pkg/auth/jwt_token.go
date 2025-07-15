package auth

/*
 * jwt.go
 *
 * This file implements a JWT (JSON Web Token) authentication service.
 * It provides token generation, validation, refresh flows, and revocation
 * without any knowledge of the domain layer.
 *
 * This implementation uses HMAC-SHA256 for token signing and includes
 * in-memory revocation tracking (should be replaced with Redis or similar in production).
 */

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Common errors returned by the service
var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrExpiredToken     = errors.New("token has expired")
	ErrRevokedToken     = errors.New("token has been revoked")
	ErrInvalidTokenType = errors.New("invalid token type")
	ErrInvalidSecret    = errors.New("invalid secret key")
	ErrInvalidDuration  = errors.New("invalid token duration")
)

// tokenType represents internal token types
type tokenType string

const (
	tokenTypeAccess  tokenType = "access"
	tokenTypeRefresh tokenType = "refresh"
)

// jwtClaims defines the JWT claims structure
type jwtClaims struct {
	UserID    int64          `json:"user_id"`
	TokenID   string         `json:"token_id"`
	TokenType tokenType      `json:"token_type"`
	Custom    map[string]any `json:"custom,omitempty"`
	jwt.RegisteredClaims
}

// JWTToken implements JWT token operations
type JWTToken struct {
	secretKey            []byte
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
	revokedTokens        sync.Map
}

// NewJWTToken creates a new JWT token instance
func NewJWTToken(secret string, accessDuration, refreshDuration time.Duration) (*JWTToken, error) {
	if len(secret) == 0 {
		return nil, ErrInvalidSecret
	}

	if accessDuration <= 0 || refreshDuration <= 0 {
		return nil, ErrInvalidDuration
	}

	return &JWTToken{
		secretKey:            []byte(secret),
		accessTokenDuration:  accessDuration,
		refreshTokenDuration: refreshDuration,
	}, nil
}

// GenerateTokens creates a new access and refresh token pair
func (s *JWTToken) GenerateTokens(ctx context.Context, issuerName string, userID int64) (accessToken, refreshToken string, err error) {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return "", "", fmt.Errorf("context error: %w", err)
	}

	// Generate access token
	accessToken, err = s.generateToken(issuerName, userID, s.accessTokenDuration, tokenTypeAccess)
	if err != nil {
		return "", "", fmt.Errorf("generate access token: %w", err)
	}

	// Generate refresh token
	refreshToken, err = s.generateToken(issuerName, userID, s.refreshTokenDuration, tokenTypeRefresh)
	if err != nil {
		return "", "", fmt.Errorf("generate refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

// RefreshTokens validates a refresh token and generates a new token pair
func (s *JWTToken) RefreshTokens(ctx context.Context, refreshToken string) (newAccessToken, newRefreshToken string, err error) {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return "", "", fmt.Errorf("context error: %w", err)
	}

	// Parse and validate the refresh token
	claims, err := s.validateToken(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("validate refresh token: %w", err)
	}

	// Ensure it's a refresh token
	if claims.TokenType != tokenTypeRefresh {
		return "", "", ErrInvalidTokenType
	}

	// Generate new token pair
	return s.GenerateTokens(ctx, claims.Issuer, claims.UserID)
}

// ValidateAccessToken validates an access token and returns the user ID
func (s *JWTToken) ValidateAccessToken(ctx context.Context, token string) (userID int64, err error) {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return 0, fmt.Errorf("context error: %w", err)
	}

	claims, err := s.validateToken(token)
	if err != nil {
		return 0, err
	}

	if claims.TokenType != tokenTypeAccess {
		return 0, ErrInvalidTokenType
	}

	return claims.UserID, nil
}

// ValidateRefreshToken validates a refresh token and returns the user ID
func (s *JWTToken) ValidateRefreshToken(ctx context.Context, token string) (userID int64, err error) {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return 0, fmt.Errorf("context error: %w", err)
	}

	claims, err := s.validateToken(token)
	if err != nil {
		return 0, err
	}

	if claims.TokenType != tokenTypeRefresh {
		return 0, ErrInvalidTokenType
	}

	return claims.UserID, nil
}

// RevokeToken marks a token as revoked
func (s *JWTToken) RevokeToken(ctx context.Context, token string) error {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("context error: %w", err)
	}

	// Parse token to get its ID
	claims, err := s.validateToken(token)
	if err != nil {
		// If token is already expired, consider it successfully revoked
		if errors.Is(err, ErrExpiredToken) {
			return nil
		}
		return err
	}

	// Store the token ID as revoked
	s.revokedTokens.Store(claims.TokenID, time.Now())

	return nil
}

// GetTokenInfo extracts all useful information from a token
func (s *JWTToken) GetTokenInfo(token string) (
	userID int64,
	issuer string,
	expiresAt, issuedAt time.Time,
	tokenID string,
	customClaims map[string]any,
	err error,
) {
	claims := &jwtClaims{}

	_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return s.secretKey, nil
	})

	if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
		return 0, "", time.Time{}, time.Time{}, "", nil, ErrInvalidToken
	}

	return claims.UserID, claims.Issuer, claims.ExpiresAt.Time, claims.IssuedAt.Time, claims.TokenID, claims.Custom, nil
}

// generateToken creates a JWT token with the specified parameters
func (s *JWTToken) generateToken(issuerName string, userID int64, duration time.Duration, tType tokenType) (string, error) {
	tokenID := generateTokenID()
	now := time.Now()
	expiresAt := now.Add(duration)

	claims := &jwtClaims{
		UserID:    userID,
		TokenID:   tokenID,
		TokenType: tType,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuerName,
			Subject:   fmt.Sprintf("%d", userID),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        tokenID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

// validateToken parses and validates a token
func (s *JWTToken) validateToken(tokenString string) (*jwtClaims, error) {
	claims := &jwtClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return s.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	// Check if token is revoked
	if _, revoked := s.revokedTokens.Load(claims.TokenID); revoked {
		return nil, ErrRevokedToken
	}

	return claims, nil
}

// generateTokenID creates a unique token identifier
func generateTokenID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}

// CleanupRevokedTokens removes expired tokens from the revocation list
// This should be called periodically in production
func (s *JWTToken) CleanupRevokedTokens(before time.Time) int {
	count := 0
	s.revokedTokens.Range(func(key, value any) bool {
		if revokedAt, ok := value.(time.Time); ok {
			if revokedAt.Before(before) {
				s.revokedTokens.Delete(key)
				count++
			}
		}
		return true
	})
	return count
}
