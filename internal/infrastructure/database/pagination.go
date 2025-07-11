package database

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/raykavin/khambalia/pkg/pagination"
	"gorm.io/gorm"
)

// PaginationRepository handles database operations related to pagination
type PaginationRepository struct {
	db *gorm.DB
}

// NewPaginationRepository creates a new pagination repository
func NewPaginationRepository(db *gorm.DB) *PaginationRepository {
	return &PaginationRepository{
		db: db,
	}
}

// ApplyPagination applies pagination parameters to the query
func (r *PaginationRepository) ApplyPagination(tx *gorm.DB, pagination *pagination.Pagination, dst any) (*gorm.DB, error) {
	// Clone the DB to not affect the original transaction
	query := tx.Session(&gorm.Session{})

	// Apply filters
	if err := r.applyFilters(query, pagination.Filters); err != nil {
		return nil, err
	}

	// Apply global filter if provided
	if pagination.GlobalFilter != nil && *pagination.GlobalFilter != "" {
		if err := r.applyGlobalFilter(query, *pagination.GlobalFilter, dst); err != nil {
			return nil, err
		}
	}

	// Apply sorting
	r.applySorting(query, pagination.SortFields)

	// Count total rows before pagination
	var totalRows int64
	err := query.Model(dst).Count(&totalRows).Error
	if err != nil {
		return nil, fmt.Errorf("failed to count total rows: %w", err)
	}
	pagination.TotalRows = totalRows

	// Calculate total pages
	pagination.CalculateTotalPages()

	// Apply limit and offset
	query = query.Limit(pagination.GetSize()).Offset(pagination.GetOffset())

	return query, nil
}

// applyFilters applies filter conditions to the query
func (r *PaginationRepository) applyFilters(tx *gorm.DB, filters []pagination.FilterPagination) error {
	dialector := tx.Dialector.Name()

	for _, filter := range filters {
		if filter.Field == "" || filter.Op == "" {
			continue
		}

		// Clean and prepare operator
		filter.Op = strings.TrimSpace(strings.ToLower(filter.Op))

		// Handle regex operator based on database dialect
		if filter.Op == "regex" {
			switch dialector {
			case "postgres":
				filter.Op = "pg_regex"
			case "mysql":
				filter.Op = "my_regex"
			case "sqlite":
				return errors.New("regex not supported for SQLite")
			default:
				return fmt.Errorf("unsupported dialect for regex: %s", dialector)
			}
		}

		// Apply filter based on operator
		switch filter.Op {
		case "=", "<", ">", "<=", ">=", "!=", "<>":
			tx.Where(fmt.Sprintf("%s %s ?", filter.Field, filter.Op), filter.Value)

		case "like":
			tx.Where(fmt.Sprintf("%s LIKE ?", filter.Field), fmt.Sprintf("%%%v%%", filter.Value))

		case "pg_regex":
			tx.Where(fmt.Sprintf("%s ~ ?", filter.Field), filter.Value)

		case "my_regex":
			tx.Where(fmt.Sprintf("%s REGEXP ?", filter.Field), filter.Value)

		case "in":
			tx.Where(fmt.Sprintf("%s IN ?", filter.Field), filter.Value)

		case "not in":
			tx.Where(fmt.Sprintf("%s NOT IN ?", filter.Field), filter.Value)

		case "between":
			vals, ok := assertBetweenValues(filter.Value)
			if !ok {
				return fmt.Errorf("invalid between values for field %s: %v", filter.Field, filter.Value)
			}
			if len(vals) != 2 {
				return fmt.Errorf("between operator requires exactly 2 values for field %s", filter.Field)
			}
			tx.Where(fmt.Sprintf("%s BETWEEN ? AND ?", filter.Field), vals[0], vals[1])

		case "is null":
			tx.Where(fmt.Sprintf("%s IS NULL", filter.Field))

		case "is not null":
			tx.Where(fmt.Sprintf("%s IS NOT NULL", filter.Field))

		default:
			return fmt.Errorf("unsupported operator '%s' for field '%s'", filter.Op, filter.Field)
		}
	}

	return nil
}

// assertBetweenValues attempts to convert the any to a slice for between operation
func assertBetweenValues(value any) ([]any, bool) {
	// Try as []any
	if vals, ok := value.([]any); ok {
		return vals, true
	}

	// Try using reflection for other slice types
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		result := make([]any, v.Len())
		for i := 0; i < v.Len(); i++ {
			result[i] = v.Index(i).Interface()
		}
		return result, true
	}

	return nil, false
}

// applyGlobalFilter applies a global search filter across multiple columns
func (r *PaginationRepository) applyGlobalFilter(tx *gorm.DB, globalFilter string, dest any) error {
	// Get model's field names using reflection
	modelType := reflect.TypeOf(dest)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	// We need a struct type
	if modelType.Kind() != reflect.Struct {
		return errors.New("destination must be a struct type")
	}

	var conditions []string
	var values []any

	// Loop through fields to build conditions
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)

		// Get field name from GORM tag or field name
		fieldName := field.Name
		gormTag := field.Tag.Get("gorm")
		if gormTag != "" {
			parts := strings.Split(gormTag, ";")
			for _, part := range parts {
				if strings.HasPrefix(part, "column:") {
					fieldName = strings.TrimPrefix(part, "column:")
					break
				}
			}
		}

		// Skip if field is not searchable
		jsonTag := field.Tag.Get("json")
		if strings.Contains(jsonTag, "nosearch") {
			continue
		}

		// Only add searchable field types (string, etc.)
		switch field.Type.Kind() {
		case reflect.String:
			conditions = append(conditions, fmt.Sprintf("%s LIKE ?", fieldName))
			values = append(values, fmt.Sprintf("%%%s%%", globalFilter))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			// Only add numeric fields if the global filter looks like a number
			if _, err := fmt.Sscanf(globalFilter, "%d", new(int)); err == nil {
				conditions = append(conditions, fmt.Sprintf("%s = ?", fieldName))
				values = append(values, globalFilter)
			}
		}
	}

	// Apply the global filter
	if len(conditions) > 0 {
		query := strings.Join(conditions, " OR ")
		args := make([]any, len(values))
		copy(args, values)
		tx.Where(query, args...)
	}

	return nil
}

// applySorting applies sorting to the query
func (r *PaginationRepository) applySorting(tx *gorm.DB, sortFields []pagination.SortField) {
	for _, sortField := range sortFields {
		if sortField.Field != "" {
			tx.Order(sortField.String())
		}
	}
}
