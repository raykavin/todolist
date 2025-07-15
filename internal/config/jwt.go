package config

import "time"

/*
 * jwt.go
 *
 * This file defines configuration settings for JWT token provider.
 *
 * Examples include the secret key and expiration time.
 *
 * These settings enable your application to generate JWT tokens.
 */

var _ JWTConfigProvider = (*jwtConfig)(nil)

type jwtConfig struct {
	SecretKey             string        `mapstructure:"secret_key"`
	ExpirationTime        time.Duration `mapstructure:"expiration_time"`
	RefreshExpirationTime time.Duration `mapstructure:"refresh_expiration_time"`
}

// GetExpirationTime implements JWTConfigProvider.
func (j *jwtConfig) GetExpirationTime() time.Duration { return j.ExpirationTime }

// GetRefreshExpiration implements JWTConfigProvider.
func (j *jwtConfig) GetRefreshExpirationTime() time.Duration { return j.RefreshExpirationTime }

// GetSecretKey implements JWTConfigProvider.
func (j *jwtConfig) GetSecretKey() string { return j.SecretKey }
