package config

import "time"

/*
 *database.go
 *
 *This file defines database-specific configuration settings.
 *
 *Examples include database driver, connection strings, credentials,
 *connection pool settings, and timeouts.
 *
 *These settings are used by repository adapters to initialize and manage database connections.
 */

type databaseService struct {
	DSN                 string        `mapstructure:"dsn"`
	Dialector           string        `mapstructure:"dialector"`
	IdleConnectionsTime time.Duration `mapstructure:"idle_connections_time"`
	IdleMaxConnections  int           `mapstructure:"idle_max_connections"`
	MaxOpenConnections  int           `mapstructure:"max_open_connections"`
}

// GetConnectionString returns the database connection string.
func (d databaseService) GetDSN() string { return d.DSN }

// GetDriver returns the database dialector name.
func (d databaseService) GetDialector() string { return d.Dialector }

// GetIdleConnsTime returns the idle connections time duration.
func (c databaseService) GetIdleConnsTime() time.Duration { return c.IdleConnectionsTime }

// GetIdleMaxConns returns the maximum number of idle connections.
func (c databaseService) GetIdleMaxConns() int { return c.IdleMaxConnections }

// GetMaxOpenConns returns the maximum number of open connections.
func (c databaseService) GetMaxOpenConns() int { return c.MaxOpenConnections }
