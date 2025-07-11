package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// Common errors
var (
	ErrMigrationFailed = errors.New("migration failed")
	ErrRecordNotFound  = gorm.ErrRecordNotFound
)

// Database represents a database connection with additional features
type Database struct {
	*gorm.DB
}

// New creates a new Database instance from a gorm.DB connection
func New(db *gorm.DB) *Database {
	return &Database{DB: db}
}

// AutoMigrate runs auto migration for given models
func (db *Database) Migrate(models ...any) error {
	if err := db.DB.AutoMigrate(models...); err != nil {
		return fmt.Errorf("%w: %v", ErrMigrationFailed, err)
	}
	return nil
}

// Close closes the database connection
func (db *Database) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Ping verifies the database connection
func (db *Database) Ping(ctx context.Context) error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

// Stats returns database statistics
func (db *Database) Stats() (sql.DBStats, error) {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return sql.DBStats{}, err
	}
	return sqlDB.Stats(), nil
}

// WithContext returns a new database instance with context
func (db *Database) WithContext(ctx context.Context) *Database {
	return &Database{
		DB: db.DB.WithContext(ctx),
	}
}

// Transaction executes a function within a database transaction
func (db *Database) Transaction(fn func(*gorm.DB) error, opts ...*sql.TxOptions) error {
	return db.DB.Transaction(fn, opts...)
}

// Repository provides common database operations
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new repository instance
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create creates a new record
func (r *Repository) Create(ctx context.Context, value any) error {
	return r.db.WithContext(ctx).Create(value).Error
}

// CreateInBatches creates records in batches
func (r *Repository) CreateInBatches(ctx context.Context, value any, batchSize int) error {
	return r.db.WithContext(ctx).CreateInBatches(value, batchSize).Error
}

// FindByID finds a record by ID
func (r *Repository) FindByID(ctx context.Context, dest any, id any) error {
	return r.db.WithContext(ctx).First(dest, id).Error
}

// FindOne finds a single record matching conditions
func (r *Repository) FindOne(ctx context.Context, dest any, conditions ...any) error {
	query := r.db.WithContext(ctx)
	if len(conditions) > 0 {
		query = query.Where(conditions[0], conditions[1:]...)
	}
	return query.First(dest).Error
}

// FindAll finds all records
func (r *Repository) FindAll(ctx context.Context, dest any, conditions ...any) error {
	query := r.db.WithContext(ctx)
	if len(conditions) > 0 {
		query = query.Where(conditions[0], conditions[1:]...)
	}
	return query.Find(dest).Error
}

// Update updates a record
func (r *Repository) Update(ctx context.Context, value any) error {
	return r.db.WithContext(ctx).Save(value).Error
}

// UpdateColumns updates specific columns
func (r *Repository) UpdateColumns(ctx context.Context, model any, columns map[string]any) error {
	return r.db.WithContext(ctx).Model(model).Updates(columns).Error
}

// Delete deletes a record
func (r *Repository) Delete(ctx context.Context, value any, conditions ...any) error {
	query := r.db.WithContext(ctx)
	if len(conditions) > 0 {
		query = query.Where(conditions[0], conditions[1:]...)
	}
	return query.Delete(value).Error
}

// Count counts records
func (r *Repository) Count(ctx context.Context, model any, conditions ...any) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(model)
	if len(conditions) > 0 {
		query = query.Where(conditions[0], conditions[1:]...)
	}
	err := query.Count(&count).Error
	return count, err
}

// Exists checks if a record exists
func (r *Repository) Exists(ctx context.Context, model any, conditions ...any) (bool, error) {
	count, err := r.Count(ctx, model, conditions...)
	return count > 0, err
}

// Preload preloads associations
func (r *Repository) Preload(query string, args ...any) *Repository {
	return &Repository{
		db: r.db.Preload(query, args...),
	}
}

// Joins adds joins to the query
func (r *Repository) Joins(query string, args ...any) *Repository {
	return &Repository{
		db: r.db.Joins(query, args...),
	}
}

// Order adds order to the query
func (r *Repository) Order(value any) *Repository {
	return &Repository{
		db: r.db.Order(value),
	}
}

// Limit adds limit to the query
func (r *Repository) Limit(limit int) *Repository {
	return &Repository{
		db: r.db.Limit(limit),
	}
}

// Offset adds offset to the query
func (r *Repository) Offset(offset int) *Repository {
	return &Repository{
		db: r.db.Offset(offset),
	}
}

// Where adds conditions to the query
func (r *Repository) Where(query any, args ...any) *Repository {
	return &Repository{
		db: r.db.Where(query, args...),
	}
}

// DB returns the underlying gorm.DB instance
func (r *Repository) DB() *gorm.DB {
	return r.db
}

// WithDB creates a new repository with a different database connection
func (r *Repository) WithDB(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Scopes applies scopes to the repository
func (r *Repository) Scopes(funcs ...func(*gorm.DB) *gorm.DB) *Repository {
	return &Repository{
		db: r.db.Scopes(funcs...),
	}
}

// QueryBuilder provides a fluent interface for building queries
type QueryBuilder struct {
	db *gorm.DB
}

// NewQueryBuilder creates a new query builder
func NewQueryBuilder(db *gorm.DB) *QueryBuilder {
	return &QueryBuilder{db: db}
}

// Model specifies the model for the query
func (qb *QueryBuilder) Model(value any) *QueryBuilder {
	qb.db = qb.db.Model(value)
	return qb
}

// Select specifies fields to select
func (qb *QueryBuilder) Select(query any, args ...any) *QueryBuilder {
	qb.db = qb.db.Select(query, args...)
	return qb
}

// Where adds WHERE conditions
func (qb *QueryBuilder) Where(query any, args ...any) *QueryBuilder {
	qb.db = qb.db.Where(query, args...)
	return qb
}

// Or adds OR conditions
func (qb *QueryBuilder) Or(query any, args ...any) *QueryBuilder {
	qb.db = qb.db.Or(query, args...)
	return qb
}

// Not adds NOT conditions
func (qb *QueryBuilder) Not(query any, args ...any) *QueryBuilder {
	qb.db = qb.db.Not(query, args...)
	return qb
}

// Joins adds JOIN clauses
func (qb *QueryBuilder) Joins(query string, args ...any) *QueryBuilder {
	qb.db = qb.db.Joins(query, args...)
	return qb
}

// Preload preloads associations
func (qb *QueryBuilder) Preload(query string, args ...any) *QueryBuilder {
	qb.db = qb.db.Preload(query, args...)
	return qb
}

// Order adds ORDER BY clause
func (qb *QueryBuilder) Order(value any) *QueryBuilder {
	qb.db = qb.db.Order(value)
	return qb
}

// Limit sets the limit
func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	qb.db = qb.db.Limit(limit)
	return qb
}

// Offset sets the offset
func (qb *QueryBuilder) Offset(offset int) *QueryBuilder {
	qb.db = qb.db.Offset(offset)
	return qb
}

// Group adds GROUP BY clause
func (qb *QueryBuilder) Group(name string) *QueryBuilder {
	qb.db = qb.db.Group(name)
	return qb
}

// Having adds HAVING clause
func (qb *QueryBuilder) Having(query any, args ...any) *QueryBuilder {
	qb.db = qb.db.Having(query, args...)
	return qb
}

// Distinct adds DISTINCT clause
func (qb *QueryBuilder) Distinct(args ...any) *QueryBuilder {
	qb.db = qb.db.Distinct(args...)
	return qb
}

// Find executes the query and stores result in dest
func (qb *QueryBuilder) Find(dest any) error {
	return qb.db.Find(dest).Error
}

// First finds the first record
func (qb *QueryBuilder) First(dest any) error {
	return qb.db.First(dest).Error
}

// Last finds the last record
func (qb *QueryBuilder) Last(dest any) error {
	return qb.db.Last(dest).Error
}

// Count counts the records
func (qb *QueryBuilder) Count(count *int64) error {
	return qb.db.Count(count).Error
}

// Pluck queries single column into slice
func (qb *QueryBuilder) Pluck(column string, dest any) error {
	return qb.db.Pluck(column, dest).Error
}

// Scan scans result into dest
func (qb *QueryBuilder) Scan(dest any) error {
	return qb.db.Scan(dest).Error
}

// DB returns the underlying gorm.DB
func (qb *QueryBuilder) DB() *gorm.DB {
	return qb.db
}
