package auth

/*
 * jwt_token.go
 *
 * This adapter connects the JWT service implementation to the domain TokenService interface.
 * It translates between the domain types and the JWT service types, keeping the JWT
 * implementation completely isolated from domain knowledge.
 */

import (
	"context"
	"errors"
	"time"
	"todolist/internal/service"
	"todolist/pkg/auth"
)

// JWTTokenAdapter adapts the JWT service to the domain TokenService interface
type JWTTokenAdapter struct {
	jwtService *auth.JWTToken
}

// NewJWTTokenAdapter creates a new adapter instance
func NewJWTTokenAdapter(secret string, accessDuration, refreshDuration time.Duration) (service.TokenService, error) {
	jwtService, err := auth.NewJWTToken(secret, accessDuration, refreshDuration)
	if err != nil {
		return nil, mapJWTError(err)
	}

	return &JWTTokenAdapter{
		jwtService: jwtService,
	}, nil
}

// GenerateTokens creates a new set of access and refresh tokens
func (a *JWTTokenAdapter) GenerateTokens(ctx context.Context, issuerName string, userID int64) (*service.AuthTokens, error) {
	accessToken, refreshToken, err := a.jwtService.GenerateTokens(ctx, issuerName, userID)
	if err != nil {
		return nil, &service.TokenServiceError{
			Code:    service.ErrCodeTokenGeneration,
			Message: "failed to generate tokens",
			Err:     err,
		}
	}

	// Get token information with full details
	_, _, accessExpiry, accessIssuedAt, accessID, _, _ := a.jwtService.GetTokenInfo(accessToken)
	_, _, refreshExpiry, refreshIssuedAt, refreshID, _, _ := a.jwtService.GetTokenInfo(refreshToken)

	return &service.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		AccessMeta: service.TokenMetadata{
			TokenID:   accessID,
			IssuedAt:  accessIssuedAt,
			ExpiresAt: accessExpiry,
			TokenType: service.TypeAccess,
		},
		RefreshMeta: service.TokenMetadata{
			TokenID:   refreshID,
			IssuedAt:  refreshIssuedAt,
			ExpiresAt: refreshExpiry,
			TokenType: service.TypeRefresh,
		},
	}, nil
}

// RefreshTokens generates new tokens from a valid refresh token
func (a *JWTTokenAdapter) RefreshTokens(ctx context.Context, refreshToken string) (*service.AuthTokens, error) {
	newAccessToken, newRefreshToken, err := a.jwtService.RefreshTokens(ctx, refreshToken)
	if err != nil {
		return nil, mapJWTError(err)
	}

	_, _, accessExpiry, accessIssuedAt, accessID, _, _ := a.jwtService.GetTokenInfo(newAccessToken)
	_, _, refreshExpiry, refreshIssuedAt, refreshID, _, _ := a.jwtService.GetTokenInfo(newRefreshToken)

	return &service.AuthTokens{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		AccessMeta: service.TokenMetadata{
			TokenID:   accessID,
			IssuedAt:  accessIssuedAt,
			ExpiresAt: accessExpiry,
			TokenType: service.TypeAccess,
		},
		RefreshMeta: service.TokenMetadata{
			TokenID:   refreshID,
			IssuedAt:  refreshIssuedAt,
			ExpiresAt: refreshExpiry,
			TokenType: service.TypeRefresh,
		},
	}, nil
}

// ValidateToken validates any type of token and returns its claims
func (a *JWTTokenAdapter) ValidateToken(ctx context.Context, token string) (*service.ValidationResult, error) {
	userID, err := a.jwtService.ValidateAccessToken(ctx, token)
	tokenType := service.TypeAccess

	if err != nil {
		if errors.Is(err, auth.ErrInvalidTokenType) {
			userID, err = a.jwtService.ValidateRefreshToken(ctx, token)
			tokenType = service.TypeRefresh
		}
		if err != nil {
			return nil, mapJWTError(err)
		}
	}

	_, issuer, expiresAt, issuedAt, tokenID, customClaims, err := a.jwtService.GetTokenInfo(token)
	if err != nil {
		return nil, &service.TokenServiceError{
			Code:    service.ErrCodeInvalidToken,
			Message: "failed to extract token information",
			Err:     err,
		}
	}

	return &service.ValidationResult{
		UserID:     userID,
		IssuerName: issuer,
		Metadata: service.TokenMetadata{
			TokenID:   tokenID,
			IssuedAt:  issuedAt,
			ExpiresAt: expiresAt,
			TokenType: tokenType,
		},
		Claims: customClaims,
	}, nil
}

// RevokeToken invalidates a specific token
func (a *JWTTokenAdapter) RevokeToken(ctx context.Context, token string) error {
	err := a.jwtService.RevokeToken(ctx, token)
	if err != nil {
		return mapJWTError(err)
	}
	return nil
}

// mapJWTError maps JWT service errors to domain errors
func mapJWTError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, auth.ErrInvalidToken):
		return &service.TokenServiceError{
			Code:    service.ErrCodeInvalidToken,
			Message: "invalid token",
			Err:     err,
		}
	case errors.Is(err, auth.ErrExpiredToken):
		return &service.TokenServiceError{
			Code:    service.ErrCodeExpiredToken,
			Message: "token has expired",
			Err:     err,
		}
	case errors.Is(err, auth.ErrRevokedToken):
		return &service.TokenServiceError{
			Code:    service.ErrCodeRevokedToken,
			Message: "token has been revoked",
			Err:     err,
		}
	case errors.Is(err, auth.ErrInvalidTokenType):
		return &service.TokenServiceError{
			Code:    service.ErrCodeInvalidTokenType,
			Message: "invalid token type",
			Err:     err,
		}
	default:
		return &service.TokenServiceError{
			Code:    service.ErrCodeTokenGeneration,
			Message: "token operation failed",
			Err:     err,
		}
	}
}
