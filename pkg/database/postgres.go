package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DatabaseConfig is a struct of database configuration
type DatabaseConfig struct {
	ApplicationName  string
	ConnectionString string
	MaxConns         int32
	MinConns         int32
	MaxConnLifetime  time.Duration
	MaxConnIdleTime  time.Duration
}

// PostgresDB is a concrete database implementation
type PostgresDB struct {
	pool *pgxpool.Pool
}

// Ensure if PostgresDB struct implements all interface methods
var _ Database = (*PostgresDB)(nil)

// NewPostresDB create a new instance of PostgreSQL database
func NewPostresDB(ctx context.Context, config *DatabaseConfig) (*PostgresDB, error) {
	// Pool configuration
	poolConfig, err := pgxpool.ParseConfig(config.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("postgres db config: %w", err)
	}

	poolConfig.MaxConns = config.MaxConns
	poolConfig.MinConns = config.MinConns
	poolConfig.MaxConnLifetime = config.MaxConnLifetime
	poolConfig.MaxConnIdleTime = config.MaxConnIdleTime

	// Connection configuration
	poolConfig.ConnConfig.ConnectTimeout = 10 * time.Second
	poolConfig.ConnConfig.RuntimeParams = map[string]string{
		"standard_conforming_strings": "on",
		"application_name":            config.ApplicationName,
	}

	// Create new pool using configuration
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("postgres db pool init: %w", err)
	}

	// Create instance
	db := &PostgresDB{pool: pool}

	// First health check to check if database connection is ok
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return db, db.HealthCheck(ctx)
}

// Close implements Database.
func (db *PostgresDB) Close() {
	if db.pool != nil {
		db.pool.Close()
	}
}

// Exec implements Database.
func (db *PostgresDB) Exec(ctx context.Context, sql string, args ...any) error {
	_, err := db.pool.Exec(ctx, sql, args...)
	return err
}

// GetPool implements Database.
func (db *PostgresDB) GetPool() *pgxpool.Pool {
	return db.pool
}

// Ping implements Database.
func (db *PostgresDB) Ping(ctx context.Context) error {
	return db.pool.Ping(ctx)
}

// Query implements Database.
func (db *PostgresDB) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return db.pool.Query(ctx, sql, args...)
}

// QueryRow implements Database.
func (db *PostgresDB) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return db.pool.QueryRow(ctx, sql, args...)
}

// Stats implements Database.
func (db *PostgresDB) Stats() *pgxpool.Stat {
	return db.pool.Stat()
}

// Transaction implements Database.
func (db *PostgresDB) Transaction(ctx context.Context, fn func(pgx.Tx) error) error {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("postgres db transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(ctx)
			panic(p)
		} else if err != nil {
			tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	err = fn(tx)
	return err
}

// HealthCheck implements Database.
func (db *PostgresDB) HealthCheck(ctx context.Context) error {
	// Ping
	if err := db.Ping(ctx); err != nil {
		return fmt.Errorf("postgres db connection ping: %w", err)
	}

	// Simple query database version to check if is responding
	var version string
	if err := db.
		QueryRow(ctx, "SELECT version()").
		Scan(&version); err != nil {
		return fmt.Errorf("postgres db get version: %w", err)
	}

	// Health check ok
	return nil
}
