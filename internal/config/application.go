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

var _ ApplicationProvider = (*application)(nil)

type application struct {
	Name        string     `mapstructure:"name"`
	Description string     `mapstructure:"description"`
	Version     string     `mapstructure:"version"`
	LogLevel    string     `mapstructure:"log_level"`
	Web         *webConfig `mapstructure:"web"`
	JWT         *jwtConfig `mapstructure:"jwt"`
}

// GetName returns the name of the application.
func (a application) GetName() string { return a.Name }

// GetDescription returns the description of the application.
func (a application) GetDescription() string { return a.Description }

// GetVersion returns the version of the application.
func (a application) GetVersion() string { return a.Version }

// GetLogLevel returns the logging level of the application.
func (a application) GetLogLevel() string { return a.LogLevel }

// GetWeb returns the web server configuration of the application.
func (a application) GetWeb() WebConfigProvider { return a.Web }

// GetJWT implements ApplicationProvider.
func (a application) GetJWT() JWTConfigProvider { return a.JWT }
