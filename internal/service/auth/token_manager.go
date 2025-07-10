package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenService defines an interface for token generation and validation.
type TokenService interface {
	Generate(userID int64) (string, error)
	Validate(token string) (int64, error)
}

// jwtTokenService implements TokenService using JWT tokens.
type jwtTokenService struct {
	secretKey     []byte
	tokenDuration time.Duration
}

// NewJWTTokenService creates a new JWT token service with the provided secret and token expiration duration.
func NewJWTTokenService(secret string, duration time.Duration) TokenService {
	return &jwtTokenService{
		secretKey:     []byte(secret),
		tokenDuration: duration,
	}
}

// jwtCustomClaims defines custom JWT claims embedding standard registered claims.
type jwtCustomClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

// Generate creates a signed JWT token containing the user ID and expiration time.
func (s *jwtTokenService) Generate(userID int64) (string, error) {
	claims := &jwtCustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "your-app-name", // replace with your app's name
		},
	}

	// Create a new token object specifying signing method and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret key
	return token.SignedString(s.secretKey)
}

// Validate parses and validates the token string, returning the user ID if valid.
func (s *jwtTokenService) Validate(tokenStr string) (int64, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
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
