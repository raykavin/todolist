package database

import (
	"context"
	"errors"
	"time"
	"todolist/pkg/log"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

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
		l.log.WithFields(fields).Infof("SQL: %s", sql)
	}
}
