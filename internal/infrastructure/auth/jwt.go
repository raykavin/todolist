package auth

/*
 * jwt.go
 *
 * This file defines the JWT (JSON Web Token) authentication provider implementation.
 *
 * It should contain logic for generating, signing, and validating JWT tokens,
 * including setting claims such as user ID, expiration, and issuer.
 *
 * This allows the application to support stateless, token-based authentication
 * between client and server.
 *
 * This implementation is used when you want lightweight, self-contained tokens
 * without external dependencies.
 */

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// jwtCustomClaims defines custom JWT claims embedding standard registered claims.
type jwtCustomClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

// jwtToken implements TokenService using JWT tokens.
type jwtToken struct {
	secretKey     []byte
	tokenDuration time.Duration
}

// NewJWTToken creates a new JWT token service with the provided secret and token expiration duration.
func NewJWTToken(secret string, duration time.Duration) (*jwtToken, error) {
	if len(secret) == 0 {
		return nil, fmt.Errorf("invalid secret key")
	}

	return &jwtToken{
		secretKey:     []byte(secret),
		tokenDuration: duration,
	}, nil
}

// Generate creates a signed JWT token containing the user ID and expiration time.
func (s *jwtToken) Generate(issuerName string, userID int64) (string, error) {
	claims := &jwtCustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    issuerName,
		},
	}

	// Create a new token object specifying signing method and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret key
	return token.SignedString(s.secretKey)
}

// Validate parses and validates the token string, returning the user ID if valid.
func (s *jwtToken) Validate(tokenStr string) (int64, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwtCustomClaims{}, func(token *jwt.Token) (any, error) {
		// Verify the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return s.secretKey, nil
	})

	if err != nil {
		return 0, err
	}

	// Extract the claims and verify token validity
	if claims, ok := token.Claims.(*jwtCustomClaims); ok && token.Valid {
		return claims.UserID, nil
	}

	return 0, errors.New("invalid token")
}
