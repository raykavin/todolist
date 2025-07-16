package di

import (
	"context"
	"fmt"
	"todolist/internal/config"
	infraDB "todolist/internal/infrastructure/database"
	"todolist/pkg/database"
	"todolist/pkg/logger"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

// DatabaseParams defines the dependencies required to create the database engine
type DatabaseParams struct {
	fx.In
	Shutdowner            fx.Shutdowner
	DefaultDatabaseConfig config.DatabaseServiceProvider
	ApplicationConfig     config.ApplicationProvider
	Log                   logger.Interface
}

// DatabaseContainer groups all database implementations provided from Fx
type DatabaseContainer struct {
	fx.Out
	DefaultDatabase *gorm.DB
}

// validateDatabaseConfig validates the database configuration
func validateDatabaseConfig(provider config.DatabaseServiceProvider) error {
	if provider.GetDSN() == "" {
		return fmt.Errorf("database DSN cannot be empty")
	}

	if len(provider.GetDialector()) == 0 {
		return fmt.Errorf("database dialector cannot be nil")
	}

	if provider.GetMaxOpenConns() < 0 {
		return fmt.Errorf("max open connections cannot be negative")
	}

	if provider.GetIdleMaxConns() < 0 {
		return fmt.Errorf("max idle connections cannot be negative")
	}

	if provider.GetIdleConnsTime() < 0 {
		return fmt.Errorf("idle connection time cannot be negative")
	}

	return nil
}

// newDatabaseContainer creates all database implementations with proper validation
func newDatabaseContainer(p DatabaseParams) (DatabaseContainer, error) {
	// Validate configuration
	if err := validateDatabaseConfig(p.DefaultDatabaseConfig); err != nil {
		return DatabaseContainer{}, fmt.Errorf("invalid database configuration: %w", err)
	}

	// Create database configuration
	dbConfig := database.DefaultConfig()

	// Apply configuration values
	if idleTime := p.DefaultDatabaseConfig.GetIdleConnsTime(); idleTime > 0 {
		dbConfig.ConnMaxIdleTime = idleTime
	}

	if maxIdle := p.DefaultDatabaseConfig.GetIdleMaxConns(); maxIdle > 0 {
		dbConfig.MaxIdleConns = maxIdle
	}

	if maxOpen := p.DefaultDatabaseConfig.GetMaxOpenConns(); maxOpen > 0 {
		dbConfig.MaxOpenConns = maxOpen
	}

	dbConfig.DSN = p.DefaultDatabaseConfig.GetDSN()
	dbConfig.Dialector = p.DefaultDatabaseConfig.GetDialector()
	dbConfig.LogLevel = p.DefaultDatabaseConfig.GetLogLevel()
	dbConfig.Logger = p.Log

	// dbConfig.LogLevel =

	// Create database connection
	db, err := database.NewWithConfig(dbConfig)
	if err != nil {
		return DatabaseContainer{}, fmt.Errorf("failed to initialize default database: %w", err)
	}

	return DatabaseContainer{DefaultDatabase: db}, nil
}

// NewDatabases creates databases with proper lifecycle management
func NewDatabases(p DatabaseParams, lc fx.Lifecycle) (DatabaseContainer, error) {
	container, err := newDatabaseContainer(p)
	if err != nil {
		return DatabaseContainer{}, err
	}

	// Register lifecycle hooks
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := infraDB.MigrateDefault(container.DefaultDatabase); err != nil {
				return fmt.Errorf("default database migration failed: %w", err)
			}

			if err := infraDB.SeedDefault(container.DefaultDatabase); err != nil {
				return fmt.Errorf("default database seeding failed: %w", err)
			}

			return nil
		},
		OnStop: func(ctx context.Context) error {
			if sqlDB, err := container.DefaultDatabase.DB(); err == nil {
				return sqlDB.Close()
			}
			return nil
		},
	})

	return container, nil
}

// DatabasesModule returns the fx module with all database dependencies
func DatabasesModule() fx.Option {
	return fx.Module("databases",
		fx.Provide(NewDatabases),
	)
}
