package config

/*
 * application.go
 *
 * This file defines application-level configuration settings.
 *
 * Examples include the application name, version, environment (development, staging, production),
 * or other global flags and constants.
 *
 * These values help structure environment-specific behaviors
 * and centralized configuration management.
 */

type application struct {
	Name        string     `mapstructure:"name"`
	Description string     `mapstructure:"description"`
	Version     string     `mapstructure:"version"`
	LoggerLevel string     `mapstructure:"logger_level"`
	Web         *webConfig `mapstructure:"web"`
}

// GetName returns the name of the application.
func (a application) GetName() string { return a.Name }

// GetDescription returns the description of the application.
func (a application) GetDescription() string { return a.Description }

// GetVersion returns the version of the application.
func (a application) GetVersion() string { return a.Version }

// GetLoggerLevel returns the logging level of the application.
func (a application) GetLoggerLevel() string { return a.LoggerLevel }

// GetWeb returns the web server configuration of the application.
func (a application) GetWeb() WebConfigProvider { return a.Web }
 