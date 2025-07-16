package database

import (
	"context"
	"errors"
	"time"
	pkgLogger "todolist/pkg/logger"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

// dbLogger adapts the custom log.Interface to GORM's logger.Interface
type dbLogger struct {
	log           pkgLogger.Interface
	LogLevel      logger.LogLevel
	SlowThreshold time.Duration
}

// LogMode implements logger.Interface
func (l *dbLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info implements logger.Interface
func (l *dbLogger) Info(ctx context.Context, msg string, data ...any) {
	if l.LogLevel >= logger.Info {
		l.log.WithField("source", utils.FileWithLineNum()).Infof(msg, data...)
	}
}

// Warn implements logger.Interface
func (l *dbLogger) Warn(ctx context.Context, msg string, data ...any) {
	if l.LogLevel >= logger.Warn {
		l.log.WithField("source", utils.FileWithLineNum()).Warnf(msg, data...)
	}
}

// Error implements logger.Interface
func (l *dbLogger) Error(ctx context.Context, msg string, data ...any) {
	if l.LogLevel >= logger.Error {
		l.log.WithField("source", utils.FileWithLineNum()).Errorf(msg, data...)
	}
}

// Trace implements logger.Interface
func (l *dbLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
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
		l.log.WithFields(fields).Infof("SQL: %s", sql)
	}
}

// WithLog creates a new database connection with a custom logger
func WithLog(db *gorm.DB, customLogger pkgLogger.Interface) *gorm.DB {
	// Create GORM logger adapter
	gormLogger := &dbLogger{
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
func WithLogConfig(db *gorm.DB, logger pkgLogger.Interface, logLevel logger.LogLevel, slowThreshold time.Duration) *gorm.DB {
	// Create GORM logger adapter with custom config
	gormLogger := &dbLogger{
		log:           logger,
		LogLevel:      logLevel,
		SlowThreshold: slowThreshold,
	}

	// Return new DB instance with custom logger
	return db.Session(&gorm.Session{
		Logger: gormLogger,
	})
}
