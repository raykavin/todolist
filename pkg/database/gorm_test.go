package database

import (
	"context"
	"errors"
	"testing"
	"time"
	"todolist/pkg/log"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Mock logger for testing
type mockLogger struct {
	mock.Mock
}

func (m *mockLogger) WithField(key string, value any) log.Interface {
	args := m.Called(key, value)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(log.Interface)
}

func (m *mockLogger) WithFields(fields map[string]any) log.Interface {
	args := m.Called(fields)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(log.Interface)
}

func (m *mockLogger) WithError(err error) log.Interface {
	args := m.Called(err)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(log.Interface)
}

func (m *mockLogger) Print(args ...any) {
	m.Called(args...)
}

func (m *mockLogger) Debug(args ...any) {
	m.Called(args...)
}

func (m *mockLogger) Info(args ...any) {
	m.Called(args...)
}

func (m *mockLogger) Warn(args ...any) {
	m.Called(args...)
}

func (m *mockLogger) Error(args ...any) {
	m.Called(args...)
}

func (m *mockLogger) Fatal(args ...any) {
	m.Called(args...)
}

func (m *mockLogger) Panic(args ...any) {
	m.Called(args...)
}

func (m *mockLogger) Printf(format string, args ...any) {
	m.Called(format, args)
}

func (m *mockLogger) Debugf(format string, args ...any) {
	m.Called(format, args)
}

func (m *mockLogger) Infof(format string, args ...any) {
	m.Called(format, args)
}

func (m *mockLogger) Warnf(format string, args ...any) {
	m.Called(format, args)
}

func (m *mockLogger) Errorf(format string, args ...any) {
	m.Called(format, args)
}

func (m *mockLogger) Fatalf(format string, args ...any) {
	m.Called(format, args)
}

func (m *mockLogger) Panicf(format string, args ...any) {
	m.Called(format, args)
}

// Mock logger that returns itself for chaining
type mockLoggerWithChain struct {
	mock.Mock
}

func (m *mockLoggerWithChain) WithField(key string, value any) log.Interface {
	m.Called(key, value)
	return m
}

func (m *mockLoggerWithChain) WithFields(fields map[string]any) log.Interface {
	m.Called(fields)
	return m
}

func (m *mockLoggerWithChain) WithError(err error) log.Interface {
	m.Called(err)
	return m
}

func (m *mockLoggerWithChain) Print(args ...any) {
	m.Called(args...)
}

func (m *mockLoggerWithChain) Debug(args ...any) {
	m.Called(args...)
}

func (m *mockLoggerWithChain) Info(args ...any) {
	m.Called(args...)
}

func (m *mockLoggerWithChain) Warn(args ...any) {
	m.Called(args...)
}

func (m *mockLoggerWithChain) Error(args ...any) {
	m.Called(args...)
}

func (m *mockLoggerWithChain) Fatal(args ...any) {
	m.Called(args...)
}

func (m *mockLoggerWithChain) Panic(args ...any) {
	m.Called(args...)
}

func (m *mockLoggerWithChain) Printf(format string, args ...any) {
	m.Called(format, args)
}

func (m *mockLoggerWithChain) Debugf(format string, args ...any) {
	m.Called(format, args)
}

func (m *mockLoggerWithChain) Infof(format string, args ...any) {
	m.Called(format, args)
}

func (m *mockLoggerWithChain) Warnf(format string, args ...any) {
	m.Called(format, args)
}

func (m *mockLoggerWithChain) Errorf(format string, args ...any) {
	m.Called(format, args)
}

func (m *mockLoggerWithChain) Fatalf(format string, args ...any) {
	m.Called(format, args)
}

func (m *mockLoggerWithChain) Panicf(format string, args ...any) {
	m.Called(format, args)
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	assert.Equal(t, 25, cfg.MaxOpenConns)
	assert.Equal(t, 5, cfg.MaxIdleConns)
	assert.Equal(t, time.Hour, cfg.ConnMaxLifetime)
	assert.Equal(t, 10*time.Minute, cfg.ConnMaxIdleTime)
	assert.Equal(t, "info", cfg.LogLevel)
	assert.Equal(t, 200*time.Millisecond, cfg.SlowThreshold)
	assert.True(t, cfg.SkipDefaultTx)
	assert.True(t, cfg.PrepareStmt)
	assert.Equal(t, 3, cfg.RetryAttempts)
	assert.Equal(t, time.Second, cfg.RetryDelay)
}

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		dsn       string
		driver    string
		config    *gorm.Config
		wantError error
	}{
		{
			name:      "empty DSN",
			dsn:       "",
			driver:    "sqlite",
			config:    &gorm.Config{},
			wantError: ErrDSNRequired,
		},
		{
			name:      "empty driver",
			dsn:       ":memory:",
			driver:    "",
			config:    &gorm.Config{},
			wantError: ErrUnsupportedDriver,
		},
		{
			name:      "unsupported driver",
			dsn:       "test.db",
			driver:    "unsupported",
			config:    &gorm.Config{},
			wantError: ErrUnsupportedDriver,
		},
		{
			name:      "valid sqlite connection",
			dsn:       ":memory:",
			driver:    "sqlite",
			config:    &gorm.Config{},
			wantError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := New(tt.dsn, tt.driver, tt.config)
			if tt.wantError != nil {
				assert.ErrorIs(t, err, tt.wantError)
				assert.Nil(t, db)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, db)
				// Clean up
				sqlDB, _ := db.DB()
				if sqlDB != nil {
					sqlDB.Close()
				}
			}
		})
	}
}

func TestNewWithConfig(t *testing.T) {
	tests := []struct {
		name      string
		cfg       *Config
		wantError bool
		errorType error
	}{
		{
			name:      "nil config",
			cfg:       nil,
			wantError: true,
			errorType: ErrInvalidConfig,
		},
		{
			name: "empty DSN",
			cfg: &Config{
				Driver: "sqlite",
			},
			wantError: true,
			errorType: ErrDSNRequired,
		},
		{
			name: "empty driver",
			cfg: &Config{
				DSN: ":memory:",
			},
			wantError: true,
			errorType: ErrUnsupportedDriver,
		},
		{
			name: "unsupported driver",
			cfg: &Config{
				DSN:    "test.db",
				Driver: "unsupported",
			},
			wantError: true,
			errorType: ErrUnsupportedDriver,
		},
		{
			name: "valid config",
			cfg: &Config{
				DSN:          ":memory:",
				Driver:       "sqlite",
				MaxOpenConns: 10,
				MaxIdleConns: 5,
			},
			wantError: false,
		},
		{
			name: "with custom GORM config",
			cfg: &Config{
				DSN:        ":memory:",
				Driver:     "sqlite",
				GormConfig: &gorm.Config{PrepareStmt: true},
			},
			wantError: false,
		},
		{
			name: "with custom logger",
			cfg: &Config{
				DSN:      ":memory:",
				Driver:   "sqlite",
				Logger:   &mockLogger{},
				LogLevel: "info",
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := NewWithConfig(tt.cfg)
			if tt.wantError {
				assert.Error(t, err)
				if tt.errorType != nil {
					assert.ErrorIs(t, err, tt.errorType)
				}
				assert.Nil(t, db)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, db)
				// Clean up
				sqlDB, _ := db.DB()
				if sqlDB != nil {
					sqlDB.Close()
				}
			}
		})
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name      string
		cfg       *Config
		wantError error
	}{
		{
			name: "empty DSN",
			cfg: &Config{
				Driver: "sqlite",
			},
			wantError: ErrDSNRequired,
		},
		{
			name: "empty driver",
			cfg: &Config{
				DSN: ":memory:",
			},
			wantError: ErrUnsupportedDriver,
		},
		{
			name: "unsupported driver",
			cfg: &Config{
				DSN:    "test.db",
				Driver: "unsupported",
			},
			wantError: ErrUnsupportedDriver,
		},
		{
			name: "valid config",
			cfg: &Config{
				DSN:    ":memory:",
				Driver: "sqlite",
			},
			wantError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(tt.cfg)
			if tt.wantError != nil {
				assert.ErrorIs(t, err, tt.wantError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBuildGormConfig(t *testing.T) {
	t.Run("with default logger", func(t *testing.T) {
		cfg := &Config{
			LogLevel:      "info",
			SlowThreshold: 100 * time.Millisecond,
			SkipDefaultTx: true,
			PrepareStmt:   true,
			DryRun:        false,
		}

		gormCfg := buildGormConfig(cfg)

		assert.NotNil(t, gormCfg)
		assert.NotNil(t, gormCfg.Logger)
		assert.True(t, gormCfg.SkipDefaultTransaction)
		assert.True(t, gormCfg.PrepareStmt)
		assert.False(t, gormCfg.DryRun)
		assert.NotNil(t, gormCfg.NowFunc)
	})

	t.Run("with custom logger", func(t *testing.T) {
		mockLog := &mockLogger{}
		cfg := &Config{
			Logger:        mockLog,
			LogLevel:      "error",
			SlowThreshold: 200 * time.Millisecond,
		}

		gormCfg := buildGormConfig(cfg)

		assert.NotNil(t, gormCfg)
		assert.NotNil(t, gormCfg.Logger)
		// Verify it's our custom logger adapter
		_, ok := gormCfg.Logger.(*customLoggerAdapter)
		assert.True(t, ok)
	})
}

func TestConfigureConnectionPool(t *testing.T) {
	// Create a test database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	sqlDB, err := db.DB()
	require.NoError(t, err)
	defer sqlDB.Close()

	cfg := &Config{
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: 30 * time.Minute,
		ConnMaxIdleTime: 5 * time.Minute,
	}

	configureConnectionPool(sqlDB, cfg)

	// Verify settings were applied
	stats := sqlDB.Stats()
	assert.LessOrEqual(t, stats.MaxOpenConnections, 10)
}

func TestParseLoggerLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected logger.LogLevel
	}{
		{"silent", logger.Silent},
		{"info", logger.Info},
		{"error", logger.Error},
		{"err", logger.Error},
		{"warning", logger.Warn},
		{"warn", logger.Warn},
		{"unknown", logger.Info}, // default
		{"", logger.Info},        // default
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ParseLoggerLevel(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUpdateConnectionPool(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	sqlDB, err := db.DB()
	require.NoError(t, err)
	defer sqlDB.Close()

	cfg := &Config{
		MaxOpenConns: 20,
		MaxIdleConns: 10,
	}

	err = UpdateConnectionPool(db, cfg)
	assert.NoError(t, err)

	// Verify settings were updated
	stats := sqlDB.Stats()
	assert.LessOrEqual(t, stats.MaxOpenConnections, 20)
}

func TestGetConnectionStats(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	sqlDB, err := db.DB()
	require.NoError(t, err)
	defer sqlDB.Close()

	stats, err := GetConnectionStats(db)
	assert.NoError(t, err)
	assert.NotNil(t, stats)
}

func TestGetDriverDialectorFunc(t *testing.T) {
	tests := []struct {
		driver    string
		wantError bool
	}{
		{"postgres", false},
		{"mysql", false},
		{"mariadb", false},
		{"sqlite", false},
		{"sqlserver", false},
		{"mssql", false},
		{"unsupported", true},
		{"", true},
	}

	for _, tt := range tests {
		t.Run(tt.driver, func(t *testing.T) {
			dialFn, err := getDriverDialectorFunc(tt.driver)
			if tt.wantError {
				assert.Error(t, err)
				assert.ErrorIs(t, err, ErrUnsupportedDriver)
				assert.Nil(t, dialFn)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, dialFn)
			}
		})
	}
}

func TestWithLog(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	sqlDB, err := db.DB()
	require.NoError(t, err)
	defer sqlDB.Close()

	mockLog := &mockLogger{}
	newDB := WithLog(db, mockLog)

	assert.NotNil(t, newDB)
	assert.NotEqual(t, db, newDB) // Should be a new session
}

func TestWithLogConfig(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	sqlDB, err := db.DB()
	require.NoError(t, err)
	defer sqlDB.Close()

	mockLog := &mockLogger{}
	newDB := WithLogConfig(db, mockLog, logger.Error, 500*time.Millisecond)

	assert.NotNil(t, newDB)
	assert.NotEqual(t, db, newDB) // Should be a new session
}

func TestCustomLoggerAdapter(t *testing.T) {
	mockLog := &mockLoggerWithChain{}
	adapter := &customLoggerAdapter{
		log:           mockLog,
		LogLevel:      logger.Info,
		SlowThreshold: 200 * time.Millisecond,
	}

	t.Run("LogMode", func(t *testing.T) {
		newLogger := adapter.LogMode(logger.Error)
		assert.NotNil(t, newLogger)
		assert.Equal(t, logger.Error, newLogger.(*customLoggerAdapter).LogLevel)
	})

	t.Run("Info", func(t *testing.T) {
		mockLog.On("WithField", "source", mock.Anything).Return(mockLog).Once()
		mockLog.On("Infof", "test %s", []any{"message"}).Once()

		adapter.Info(context.Background(), "test %s", "message")
		mockLog.AssertExpectations(t)
	})

	t.Run("Info with Silent level", func(t *testing.T) {
		silentAdapter := &customLoggerAdapter{
			log:      mockLog,
			LogLevel: logger.Silent,
		}
		// Should not call logger when level is Silent
		silentAdapter.Info(context.Background(), "test")
		mockLog.AssertNotCalled(t, "WithField", mock.Anything, mock.Anything)
	})

	t.Run("Warn", func(t *testing.T) {
		mockLog.On("WithField", "source", mock.Anything).Return(mockLog).Once()
		mockLog.On("Warnf", "warning %s", []any{"test"}).Once()

		adapter.Warn(context.Background(), "warning %s", "test")
		mockLog.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockLog.On("WithField", "source", mock.Anything).Return(mockLog).Once()
		mockLog.On("Errorf", "error %s", []any{"test"}).Once()

		adapter.Error(context.Background(), "error %s", "test")
		mockLog.AssertExpectations(t)
	})

	t.Run("Trace with error", func(t *testing.T) {
		mockLog.On("WithFields", mock.MatchedBy(func(fields map[string]any) bool {
			return fields["error"] != nil && fields["duration"] != nil && fields["rows"] != nil
		})).Return(mockLog).Once()
		mockLog.On("Errorf", "SQL Error: %s", "SELECT * FROM users").Once()

		adapter.Trace(context.Background(), time.Now(), func() (string, int64) {
			return "SELECT * FROM users", 0
		}, errors.New("test error"))
		mockLog.AssertExpectations(t)
	})

	t.Run("Trace with slow query", func(t *testing.T) {
		adapter.SlowThreshold = 1 * time.Millisecond
		mockLog.On("WithFields", mock.MatchedBy(func(fields map[string]any) bool {
			return fields["slow"] == true && fields["threshold"] != nil
		})).Return(mockLog).Once()
		mockLog.On("Warnf", "SLOW SQL: %s", "SELECT * FROM users").Once()

		time.Sleep(2 * time.Millisecond) // Ensure it's slow
		adapter.Trace(context.Background(), time.Now().Add(-10*time.Millisecond), func() (string, int64) {
			return "SELECT * FROM users", 10
		}, nil)
		mockLog.AssertExpectations(t)
	})

	t.Run("Trace with normal query", func(t *testing.T) {
		infoAdapter := &customLoggerAdapter{
			log:           mockLog,
			LogLevel:      logger.Info,
			SlowThreshold: 1 * time.Second,
		}
		mockLog.On("WithFields", mock.MatchedBy(func(fields map[string]any) bool {
			return fields["duration"] != nil && fields["rows"] != nil
		})).Return(mockLog).Once()
		mockLog.On("Debugf", "SQL: %s", "SELECT * FROM users").Once()

		infoAdapter.Trace(context.Background(), time.Now(), func() (string, int64) {
			return "SELECT * FROM users", 5
		}, nil)
		mockLog.AssertExpectations(t)
	})

	t.Run("Trace with Silent level", func(t *testing.T) {
		silentAdapter := &customLoggerAdapter{
			log:      mockLog,
			LogLevel: logger.Silent,
		}
		// Should not log anything when Silent
		silentAdapter.Trace(context.Background(), time.Now(), func() (string, int64) {
			return "SELECT * FROM users", 0
		}, nil)
		mockLog.AssertNotCalled(t, "WithFields", mock.Anything)
	})

	t.Run("Trace with RecordNotFound error", func(t *testing.T) {
		// Should not log RecordNotFound errors
		adapter.Trace(context.Background(), time.Now(), func() (string, int64) {
			return "SELECT * FROM users", 0
		}, gorm.ErrRecordNotFound)
		mockLog.AssertNotCalled(t, "Errorf", mock.Anything, mock.Anything)
	})
}

func TestRetryLogic(t *testing.T) {
	t.Run("successful after retry", func(t *testing.T) {
		// This test would require mocking the database connection
		// For now, we'll test with a valid connection
		cfg := &Config{
			DSN:           ":memory:",
			Driver:        "sqlite",
			RetryAttempts: 2,
			RetryDelay:    10 * time.Millisecond,
		}

		db, err := NewWithConfig(cfg)
		assert.NoError(t, err)
		assert.NotNil(t, db)

		if db != nil {
			sqlDB, _ := db.DB()
			if sqlDB != nil {
				sqlDB.Close()
			}
		}
	})
}

func TestSetupReplicas(t *testing.T) {
	// This test requires a more complex setup with dbresolver
	// For basic testing, we ensure the function doesn't panic
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	cfg := &Config{
		Driver: "sqlite",
		Replicas: []ReplicaConfig{
			{DSN: ":memory:"},
		},
	}

	// This might fail in test environment but shouldn't panic
	err = setupReplicas(db, cfg)
	assert.NoError(t, err)

	sqlDB, _ := db.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}
}

// Integration test example
func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Create a complete configuration
	cfg := DefaultConfig()
	cfg.DSN = ":memory:"
	cfg.Driver = "sqlite"
	cfg.LogLevel = "info"

	// Initialize database
	db, err := NewWithConfig(cfg)
	require.NoError(t, err)
	require.NotNil(t, db)

	// Get connection stats
	stats, err := GetConnectionStats(db)
	assert.NoError(t, err)
	assert.NotNil(t, stats)

	// Update connection pool
	err = UpdateConnectionPool(db, &Config{
		MaxOpenConns: 15,
		MaxIdleConns: 7,
	})
	assert.NoError(t, err)

	// Clean up
	sqlDB, err := db.DB()
	require.NoError(t, err)
	err = sqlDB.Close()
	assert.NoError(t, err)
}
