package di

import (
	"fmt"
	"slices"
	"strings"
	"todolist/internal/config"
	infraLog "todolist/internal/infrastructure/logger"
	"todolist/pkg/log"

	"go.uber.org/fx"
)

// LoggerParams defines the dependencies to global logger
type LoggerParams struct {
	fx.In
	Config config.ApplicationProvider
}

// LoggerConfig holds logger configuration constants
type LoggerConfig struct {
	TimeFormat    string
	ColorEnabled  bool
	CallerEnabled bool
}

// DefaultLoggerConfig returns default logger configuration
func DefaultLoggerConfig() LoggerConfig {
	return LoggerConfig{
		TimeFormat:    "2006-01-02 15:04:05",
		ColorEnabled:  true,
		CallerEnabled: false,
	}
}

// validateLoggerLevel validates the logger level
func validateLoggerLevel(level string) error {
	if level == "" {
		return fmt.Errorf("logger level cannot be empty")
	}

	validLevels := []string{"debug", "info", "warn", "error", "fatal", "panic"}
	level = strings.ToLower(level)

	if slices.Contains(validLevels, level) {
		return nil
	}

	return fmt.Errorf("invalid logger level '%s', must be one of: %v", level, validLevels)
}

// NewLogger creates a new logger instance with validation
func NewLogger(params LoggerParams) (log.ExtendedLog, error) {
	// Validate logger level
	logLevel := params.Config.GetLogLevel()
	if err := validateLoggerLevel(logLevel); err != nil {
		return nil, fmt.Errorf("invalid logger configuration: %w", err)
	}

	// Get default configuration
	config := DefaultLoggerConfig()

	// Create smart logger with validated configuration
	logger, err := log.NewSmartLog(
		logLevel,
		config.TimeFormat,
		config.ColorEnabled,
		config.CallerEnabled,
		false,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize smart logger: %w", err)
	}

	// Wrap with adapter
	return &infraLog.SmartLogAdapter{Logger: logger}, nil
}

// LoggerModule returns the fx module with logger dependencies
func LoggerModule() fx.Option {
	return fx.Module("logger",
		fx.Provide(NewLogger),
	)
}
