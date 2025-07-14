package config

import "fmt"

/*
 * config.go
 *
 * This file provides the core logic for loading and managing configuration.
 *
 * It should handle loading settings from files (e.g., YAML, JSON, ENV),
 * parsing environment variables, and binding values to configuration structs.
 *
 * This keeps the rest of the application decoupled from configuration file formats.
 */

type Config struct {
	Application application                `mapstructure:"application"`
	Databases   map[string]databaseService `mapstructure:"databases"`
	// Cache        map[string]CacheService        `mapstructure:"cache"`         // Map of cache services by name
	// MessageQueue map[string]MessageQueueService `mapstructure:"message_queue"` // Map of message queue services by name
	// Storage      map[string]StorageService      `mapstructure:"storage"`       // Map of storage services by name
}

// GetApplication returns the application configuration.
func (c Config) GetApplication() ApplicationProvider { return c.Application }

// GetDatabase returns a database service by name.
func (c Config) GetDatabase(name string) (DatabaseServiceProvider, error) {
	if db, ok := c.Databases[name]; ok {
		return db, nil
	}
	return nil, fmt.Errorf("%s: %w", name, ErrDatabaseNotFound)
}
