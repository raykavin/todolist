package di

import (
	"fmt"
	"slices"
	"strings"
	adptLogger "todolist/internal/adapter/logger"
	"todolist/internal/config"
	"todolist/pkg/logger"

	"go.uber.org/fx"
)

// LoggerParams defines the dependencies to global logger
type LoggerParams struct {
	fx.In
	Config config.ApplicationProvider
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
func NewLogger(params LoggerParams) (logger.ExtendedLog, error) {
	// Validate logger level
	logLevel := params.Config.GetLogLevel()
	if err := validateLoggerLevel(logLevel); err != nil {
		return nil, fmt.Errorf("invalid logger configuration: %w", err)
	}

	// Create smart logger with validated configuration
	logger, err := logger.New(&logger.Config{
		Level:          logLevel,
		DateTimeLayout: "2006-01-02 15:04:05",
		Colored:        true,
		JSONFormat:     false,
		UseEmoji:       false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize smart logger: %w", err)
	}

	// Wrap with adapter
	return &adptLogger.SmartLogAdapter{Logger: logger}, nil
}

// LoggerModule returns the fx module with logger dependencies
func LoggerModule() fx.Option {
	return fx.Module("logger",
		fx.Provide(NewLogger),
	)
}
