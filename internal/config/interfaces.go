package config

import (
	"errors"
	"time"
)

/*
 * interfaces.go
 *
 * This file defines interfaces for the configuration layer.
 *
 * Example: an `Config` interface that exposes strongly typed accessors
 * for retrieving configuration values in a consistent way.
 *
 * This abstraction makes it easy to swap the underlying implementation
 * or mock configuration in tests.
 */

var (
	ErrDatabaseNotFound = errors.New("database not found")
)

// ConfigProvider is the main interface for accessing application configuration.
type ConfigProvider interface {
	GetApplication() ApplicationProvider                      // Returns the application configuration
	GetDatabase(name string) (DatabaseServiceProvider, error) // Returns a database service by name
	// GetCache(name string) (CacheServiceProvider, error)               // Returns the cache service
	// GetMessageQueue(name string) (MessageQueueServiceProvider, error) // Returns the message queue service
	// GetStorage(name string) (StorageServiceProvider, error)           // Returns the storage service by name
}

// ApplicationProvider represents the main application configuration.
type ApplicationProvider interface {
	GetName() string           // Name of the application
	GetDescription() string    // Description of the application
	GetVersion() string        // Version of the application
	GetLoggerLevel() string    // Logging level (e.g., "debug", "info", "warn", "error")
	GetWeb() WebConfigProvider // Web server settings
	// GetOIDC() OIDCConfigProvider // OIDC settings
}

// WebConfigProvider defines the configuration for the web server
type WebConfigProvider interface {
	GetListen() uint16              // Port to listen on
	GetReadTimeout() time.Duration  // Duration for reading request
	GetWriteTimeout() time.Duration // Duration for writing response
	GetSSLCert() string             // Path to SSL certificate file
	GetSSLKey() string              // Path to SSL key file
	GetNoRouteTo() string           // Custom 404 page path
	GetCORS() map[string]string     // CORS headers
}

// DatabaseServiceProvider defines the interface for a database service
type DatabaseServiceProvider interface {
	GetDriver() string               // Returns the database driver name
	GetConnString() string           // Returns the database connection string
	GetIdleConnsTime() time.Duration // Returns the idle connections time duration
	GetIdleMaxConns() int            // Returns the maximum number of idle connections
	GetMaxOpenConns() int            // Returns the maximum number of open connections
}
