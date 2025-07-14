// Package gorm provides GORM database configuration and initialization
package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	stdlog "log"
	"os"
	"time"
	"todolist/pkg/log"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"gorm.io/plugin/dbresolver"
)

// Common errors
var (
	ErrUnsupportedDriver = errors.New("unsupported database dialector")
	ErrDSNRequired       = errors.New("DSN is required")
	ErrInvalidConfig     = errors.New("invalid configuration")
	ErrConnectionFailed  = errors.New("failed to connect to database")
)

// Config holds database configuration
type Config struct {
	DSN             string        // Data Source Name for the database connection
	Dialector       string        // Database dialector (e.g., "postgres", "mysql", "sqlite", "sqlserver")
	LogLevel        string        // Log level for GORM (e.g., "silent", "info", "error", "warning")
	MaxOpenConns    int           // Maximum number of open connections to the database
	MaxIdleConns    int           // Maximum number of connections in the idle connection pool
	SkipDefaultTx   bool          // Skip default transaction for single create, update, delete operations
	PrepareStmt     bool          // Executes the given query in cached statement
	DryRun          bool          // Generate SQL without executing
	ConnMaxLifetime time.Duration // Maximum amount of time a connection may be reused
	ConnMaxIdleTime time.Duration // Maximum amount of time a connection may be idle before being closed
	SlowThreshold   time.Duration // Threshold for logging slow queries

	// Custom logger
	Logger log.Interface

	// Read/Write splitting
	Replicas []ReplicaConfig // List of read replicas for load balancing

	// Connection retry
	RetryAttempts int           // Number of retry attempts for connection failures
	RetryDelay    time.Duration // Delay between retry attempts

	// GORM Config override
	GormConfig *gorm.Config // Custom GORM configuration, if nil defaults will be used
}

// ReplicaConfig holds configuration for read replicas
type ReplicaConfig struct {
	DSN string
}

// DefaultConfig returns a Config with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour,
		ConnMaxIdleTime: 10 * time.Minute,
		LogLevel:        "info",
		SlowThreshold:   200 * time.Millisecond,
		SkipDefaultTx:   true,
		PrepareStmt:     true,
		RetryAttempts:   3,
		RetryDelay:      time.Second,
	}
}

// New initializes and returns a new instance of *gorm.DB
func New(dsn, dialector string, config *gorm.Config) (*gorm.DB, error) {
	// Validate inputs
	if dsn == "" {
		return nil, ErrDSNRequired
	}
	if dialector == "" {
		return nil, ErrUnsupportedDriver
	}

	// Get the GORM dialector function based on the dialector string
	dialFn, err := getDriverDialectorFunc(dialector)
	if err != nil {
		return nil, err
	}

	// Open the database connection
	conn, err := gorm.Open(dialFn(dsn), config)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// NewWithConfig initializes database with detailed configuration
func NewWithConfig(cfg *Config) (*gorm.DB, error) {
	if cfg == nil {
		return nil, ErrInvalidConfig
	}

	// Validate configuration
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	// Use provided GORM config or build default
	gormConfig := cfg.GormConfig
	if gormConfig == nil {
		gormConfig = buildGormConfig(cfg)
	}

	// Get the GORM dialector function based on the dialector
	dialFn, err := getDriverDialectorFunc(cfg.Dialector)
	if err != nil {
		return nil, err
	}

	// Open the database connection with retry logic
	var conn *gorm.DB
	for i := 0; i <= cfg.RetryAttempts; i++ {
		conn, err = gorm.Open(dialFn(cfg.DSN), gormConfig)
		if err == nil {
			break
		}
		if i < cfg.RetryAttempts {
			time.Sleep(cfg.RetryDelay)
		}
	}
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConnectionFailed, err)
	}

	// Configure connection pool
	sqlDB, err := conn.DB()
	if err != nil {
		return nil, err
	}
	configureConnectionPool(sqlDB, cfg)

	// Setup read replicas if configured
	if len(cfg.Replicas) > 0 {
		if err := setupReplicas(conn, cfg); err != nil {
			return nil, err
		}
	}

	return conn, nil
}

// validateConfig validates the database configuration
func validateConfig(cfg *Config) error {
	if cfg.DSN == "" {
		return ErrDSNRequired
	}
	if cfg.Dialector == "" {
		return ErrUnsupportedDriver
	}

	// Validate if dialector is supported
	_, err := getDriverDialectorFunc(cfg.Dialector)
	if err != nil {
		return err
	}

	return nil
}

// buildGormConfig builds GORM configuration from our config
func buildGormConfig(cfg *Config) *gorm.Config {
	var gormLogger logger.Interface

	// Use custom logger if provided
	if cfg.Logger != nil {
		gormLogger = &customLoggerAdapter{
			log:           cfg.Logger,
			LogLevel:      ParseLoggerLevel(cfg.LogLevel),
			SlowThreshold: cfg.SlowThreshold,
		}
	} else {
		// Create default logger with slow threshold
		gormLogger = logger.New(
			stdlog.New(os.Stdout, "\r\n", stdlog.LstdFlags),
			logger.Config{
				SlowThreshold:             cfg.SlowThreshold,
				LogLevel:                  ParseLoggerLevel(cfg.LogLevel),
				IgnoreRecordNotFoundError: true,
				Colorful:                  false,
			},
		)
	}

	return &gorm.Config{
		Logger:                 gormLogger,
		SkipDefaultTransaction: cfg.SkipDefaultTx,
		PrepareStmt:            cfg.PrepareStmt,
		DryRun:                 cfg.DryRun,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}
}

// configureConnectionPool configures the database connection pool
func configureConnectionPool(sqlDB *sql.DB, cfg *Config) {
	if cfg.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	}
	if cfg.ConnMaxIdleTime > 0 {
		sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	}
}

// setupReplicas configures read replicas
func setupReplicas(db *gorm.DB, cfg *Config) error {
	dialFn, err := getDriverDialectorFunc(cfg.Dialector)
	if err != nil {
		return err
	}

	replicas := make([]gorm.Dialector, len(cfg.Replicas))
	for i, replica := range cfg.Replicas {
		replicas[i] = dialFn(replica.DSN)
	}

	return db.Use(
		dbresolver.Register(dbresolver.Config{
			Replicas: replicas,
			Policy:   dbresolver.RandomPolicy{},
		}),
	)
}

// ParseLoggerLevel maps the logger level string to gorm's LogLevel
func ParseLoggerLevel(levelStr string) logger.LogLevel {
	levels := map[string]logger.LogLevel{
		"silent":  logger.Silent,
		"info":    logger.Info,
		"error":   logger.Error,
		"err":     logger.Error,
		"warning": logger.Warn,
		"warn":    logger.Warn,
	}

	if logLevel, found := levels[levelStr]; found {
		return logLevel
	}
	return logger.Info
}

// UpdateConnectionPool updates the connection pool settings for an existing connection
func UpdateConnectionPool(db *gorm.DB, cfg *Config) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	configureConnectionPool(sqlDB, cfg)
	return nil
}

// GetConnectionStats returns current connection pool statistics
func GetConnectionStats(db *gorm.DB) (sql.DBStats, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return sql.DBStats{}, err
	}

	return sqlDB.Stats(), nil
}

// getDriverDialectorFunc returns a function that provides the appropriate GORM dialector based on the dialector
func getDriverDialectorFunc(dialector string) (func(string) gorm.Dialector, error) {
	dbDrivers := map[string]func(string) gorm.Dialector{
		"postgres":  postgres.Open,
		"mysql":     mysql.Open,
		"mariadb":   mysql.Open,
		"sqlite":    sqlite.Open,
		"sqlserver": sqlserver.Open,
		"mssql":     sqlserver.Open,
	}

	if dialFn, exists := dbDrivers[dialector]; exists {
		return dialFn, nil
	}

	return nil, fmt.Errorf("%w: %s", ErrUnsupportedDriver, dialector)
}

// WithLog creates a new database connection with a custom logger
func WithLog(db *gorm.DB, customLogger log.Interface) *gorm.DB {
	// Create GORM logger adapter
	gormLogger := &customLoggerAdapter{
		log:           customLogger,
		LogLevel:      logger.Info,
		SlowThreshold: 200 * time.Millisecond,
	}

	// Return new DB instance with custom logger
	return db.Session(&gorm.Session{
		Logger: gormLogger,
	})
}

// WithLogConfig creates a new database connection with a custom logger and configuration
func WithLogConfig(db *gorm.DB, customLogger log.Interface, logLevel logger.LogLevel, slowThreshold time.Duration) *gorm.DB {
	// Create GORM logger adapter with custom config
	gormLogger := &customLoggerAdapter{
		log:           customLogger,
		LogLevel:      logLevel,
		SlowThreshold: slowThreshold,
	}

	// Return new DB instance with custom logger
	return db.Session(&gorm.Session{
		Logger: gormLogger,
	})
}

// customLoggerAdapter adapts the custom log.Interface to GORM's logger.Interface
type customLoggerAdapter struct {
	log           log.Interface
	LogLevel      logger.LogLevel
	SlowThreshold time.Duration
}

// LogMode implements logger.Interface
func (l *customLoggerAdapter) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info implements logger.Interface
func (l *customLoggerAdapter) Info(ctx context.Context, msg string, data ...any) {
	if l.LogLevel >= logger.Info {
		l.log.WithField("source", utils.FileWithLineNum()).Infof(msg, data...)
	}
}

// Warn implements logger.Interface
func (l *customLoggerAdapter) Warn(ctx context.Context, msg string, data ...any) {
	if l.LogLevel >= logger.Warn {
		l.log.WithField("source", utils.FileWithLineNum()).Warnf(msg, data...)
	}
}

// Error implements logger.Interface
func (l *customLoggerAdapter) Error(ctx context.Context, msg string, data ...any) {
	if l.LogLevel >= logger.Error {
		l.log.WithField("source", utils.FileWithLineNum()).Errorf(msg, data...)
	}
}

// Trace implements logger.Interface
func (l *customLoggerAdapter) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()
	fields := map[string]any{
		"duration": elapsed.Milliseconds(),
		"rows":     rows,
		"source":   utils.FileWithLineNum(),
	}

	switch {
	case err != nil && l.LogLevel >= logger.Error && !errors.Is(err, gorm.ErrRecordNotFound):
		fields["error"] = err
		l.log.WithFields(fields).Errorf("SQL Error: %s", sql)

	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		fields["slow"] = true
		fields["threshold"] = l.SlowThreshold.Milliseconds()
		l.log.WithFields(fields).Warnf("SLOW SQL: %s", sql)

	case l.LogLevel == logger.Info:
		l.log.WithFields(fields).Debugf("SQL: %s", sql)
	}
}
