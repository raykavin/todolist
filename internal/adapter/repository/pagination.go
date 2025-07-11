package repository

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"todolist/internal/domain/shared"
	"todolist/pkg"

	"gorm.io/gorm"
)

// PaginationAdapter handles database pagination operations
type PaginationAdapter struct {
	db *gorm.DB
}

// NewPaginationAdapter creates a new pagination adapter
func NewPaginationAdapter(db *gorm.DB) *PaginationAdapter {
	return &PaginationAdapter{db: db}
}

// Paginate applies pagination to a GORM query
func (a *PaginationAdapter) Paginate(
	tx *gorm.DB,
	req *pkg.Request,
	model any,
) (*gorm.DB, *pkg.Result, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, nil, err
	}

	// Clone the transaction to avoid modifying the original
	query := tx.Session(&gorm.Session{})

	// Apply filters
	if err := a.applyFilters(query, req.Filters); err != nil {
		return nil, nil, fmt.Errorf("failed to apply filters: %w", err)
	}

	// Apply global filter if provided
	if req.GlobalFilter != "" {
		if err := a.applyGlobalFilter(query, req.GlobalFilter, model); err != nil {
			return nil, nil, fmt.Errorf("failed to apply global filter: %w", err)
		}
	}

	// Apply sorting
	a.applySorting(query, req.SortFields)

	// Count total rows before pagination
	var totalRows int64
	if err := query.Model(model).Count(&totalRows).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to count total rows: %w", err)
	}

	// Calculate pagination result
	result := &pkg.Result{
		TotalRows:   totalRows,
		TotalPages:  pkg.CalculateTotalPages(totalRows, req.GetSize()),
		CurrentPage: req.GetPage(),
		PageSize:    req.GetSize(),
	}

	// Apply limit and offset
	query = query.Limit(req.GetSize()).Offset(req.GetOffset())

	return query, result, nil
}

// applyFilters applies filter conditions to the query
func (a *PaginationAdapter) applyFilters(tx *gorm.DB, filters []pkg.Filter) error {
	for _, filter := range filters {
		if filter.Field == "" {
			continue
		}

		operator := pkg.NormalizeOperator(filter.Operator)

		switch operator {
		case pkg.OperatorEqual:
			tx.Where(fmt.Sprintf("%s = ?", filter.Field), filter.Value)

		case pkg.OperatorNotEqual:
			tx.Where(fmt.Sprintf("%s != ?", filter.Field), filter.Value)

		case pkg.OperatorGreater:
			tx.Where(fmt.Sprintf("%s > ?", filter.Field), filter.Value)

		case pkg.OperatorLess:
			tx.Where(fmt.Sprintf("%s < ?", filter.Field), filter.Value)

		case pkg.OperatorGreaterOrEqual:
			tx.Where(fmt.Sprintf("%s >= ?", filter.Field), filter.Value)

		case pkg.OperatorLessOrEqual:
			tx.Where(fmt.Sprintf("%s <= ?", filter.Field), filter.Value)

		case pkg.OperatorLike:
			tx.Where(fmt.Sprintf("%s LIKE ?", filter.Field), fmt.Sprintf("%%%v%%", filter.Value))

		case pkg.OperatorIn:
			tx.Where(fmt.Sprintf("%s IN ?", filter.Field), filter.Value)

		case pkg.OperatorNotIn:
			tx.Where(fmt.Sprintf("%s NOT IN ?", filter.Field), filter.Value)

		case pkg.OperatorBetween:
			values, err := a.parseBetweenValues(filter.Value)
			if err != nil {
				return fmt.Errorf("invalid between values for field %s: %w", filter.Field, err)
			}
			tx.Where(fmt.Sprintf("%s BETWEEN ? AND ?", filter.Field), values[0], values[1])

		case pkg.OperatorRegex:
			if err := a.applyRegexFilter(tx, filter.Field, filter.Value); err != nil {
				return err
			}

		case pkg.OperatorIsNull:
			tx.Where(fmt.Sprintf("%s IS NULL", filter.Field))

		case pkg.OperatorIsNotNull:
			tx.Where(fmt.Sprintf("%s IS NOT NULL", filter.Field))

		default:
			return fmt.Errorf("unsupported operator '%s' for field '%s'", filter.Operator, filter.Field)
		}
	}

	return nil
}

// applyRegexFilter applies regex filter based on database dialect
func (a *PaginationAdapter) applyRegexFilter(tx *gorm.DB, field string, value any) error {
	dialect := tx.Dialector.Name()

	switch dialect {
	case "postgres":
		tx.Where(fmt.Sprintf("%s ~ ?", field), value)
	case "mysql":
		tx.Where(fmt.Sprintf("%s REGEXP ?", field), value)
	case "sqlite":
		return errors.New("regex not supported for SQLite")
	default:
		return fmt.Errorf("unsupported dialect for regex: %s", dialect)
	}

	return nil
}

// parseBetweenValues parses and validates between operator values
func (a *PaginationAdapter) parseBetweenValues(value any) ([]any, error) {
	// Try direct slice assertion
	if vals, ok := value.([]any); ok {
		if len(vals) != 2 {
			return nil, errors.New("between operator requires exactly 2 values")
		}
		return vals, nil
	}

	// Try reflection for other slice types
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return nil, errors.New("between operator value must be a slice")
	}

	if v.Len() != 2 {
		return nil, errors.New("between operator requires exactly 2 values")
	}

	result := make([]any, 2)
	result[0] = v.Index(0).Interface()
	result[1] = v.Index(1).Interface()

	return result, nil
}

// applyGlobalFilter applies a global search filter across searchable fields
func (a *PaginationAdapter) applyGlobalFilter(tx *gorm.DB, globalFilter string, model any) error {
	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	if modelType.Kind() != reflect.Struct {
		return errors.New("model must be a struct type")
	}

	var conditions []string
	var args []any

	// Build conditions for searchable fields
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)

		// Skip non-searchable fields
		if !a.isSearchableField(field) {
			continue
		}

		fieldName := a.getFieldName(field)

		switch field.Type.Kind() {
		case reflect.String:
			conditions = append(conditions, fmt.Sprintf("%s LIKE ?", fieldName))
			args = append(args, fmt.Sprintf("%%%s%%", globalFilter))

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			// Only add numeric fields if the filter looks numeric
			if _, err := fmt.Sscanf(globalFilter, "%d", new(int)); err == nil {
				conditions = append(conditions, fmt.Sprintf("%s = ?", fieldName))
				args = append(args, globalFilter)
			}
		}
	}

	// Apply the combined condition
	if len(conditions) > 0 {
		tx.Where(strings.Join(conditions, " OR "), args...)
	}

	return nil
}

// isSearchableField determines if a field should be included in global search
func (a *PaginationAdapter) isSearchableField(field reflect.StructField) bool {
	// Check for explicit nosearch tag
	if tag := field.Tag.Get("search"); tag == "false" || tag == "no" {
		return false
	}

	// Check JSON tag for nosearch
	if jsonTag := field.Tag.Get("json"); strings.Contains(jsonTag, "nosearch") {
		return false
	}

	// Skip private fields
	if !field.IsExported() {
		return false
	}

	// Skip common non-searchable fields
	fieldName := strings.ToLower(field.Name)
	nonSearchable := []string{"id", "createdat", "updatedat", "deletedat", "password"}
	for _, ns := range nonSearchable {
		if fieldName == ns {
			return false
		}
	}

	return true
}

// getFieldName extracts the database field name from struct field
func (a *PaginationAdapter) getFieldName(field reflect.StructField) string {
	// Check GORM column tag
	if gormTag := field.Tag.Get("gorm"); gormTag != "" {
		parts := strings.Split(gormTag, ";")
		for _, part := range parts {
			if strings.HasPrefix(part, "column:") {
				return strings.TrimPrefix(part, "column:")
			}
		}
	}

	// Default to snake_case field name
	return toSnakeCase(field.Name)
}

// applySorting applies sorting to the query
func (a *PaginationAdapter) applySorting(tx *gorm.DB, sortFields []pkg.SortField) {
	for _, sortField := range sortFields {
		if sortField.Field != "" {
			tx.Order(sortField.String())
		}
	}

	// Apply default sorting if none specified
	if len(sortFields) == 0 {
		tx.Order("id DESC")
	}
}

// toSnakeCase converts PascalCase to snake_case
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

// ConvertToSharedFilters converts pagination filters to shared domain filters
func ConvertToSharedFilters(filters []pkg.Filter) []shared.Filter {
	result := make([]shared.Filter, 0, len(filters))

	for _, f := range filters {
		sharedFilter := shared.Filter{
			Field: f.Field,
			Value: f.Value,
		}

		// Map operators to shared.FilterOperator
		switch pkg.NormalizeOperator(f.Operator) {
		case pkg.OperatorEqual:
			sharedFilter.Operator = shared.FilterOperatorEqual
		case pkg.OperatorNotEqual:
			sharedFilter.Operator = shared.FilterOperatorNotEqual
		case pkg.OperatorGreater:
			sharedFilter.Operator = shared.FilterOperatorGreater
		case pkg.OperatorLess:
			sharedFilter.Operator = shared.FilterOperatorLess
		case pkg.OperatorGreaterOrEqual:
			sharedFilter.Operator = shared.FilterOperatorGreaterOrEqual
		case pkg.OperatorLessOrEqual:
			sharedFilter.Operator = shared.FilterOperatorLessOrEqual
		case pkg.OperatorLike:
			sharedFilter.Operator = shared.FilterOperatorLike
		case pkg.OperatorIn:
			sharedFilter.Operator = shared.FilterOperatorIn
		case pkg.OperatorBetween:
			sharedFilter.Operator = shared.FilterOperatorBetween
		default:
			// For extended operators, use the string value
			sharedFilter.Operator = shared.FilterOperator(f.Operator)
		}

		result = append(result, sharedFilter)
	}

	return result
}

// ConvertFromSharedFilters converts shared domain filters to pagination filters
func ConvertFromSharedFilters(filters []shared.Filter) []pkg.Filter {
	result := make([]pkg.Filter, 0, len(filters))

	for _, f := range filters {
		paginationFilter := pkg.Filter{
			Field: f.Field,
			Value: f.Value,
		}

		// Map shared operators to pagination operators
		switch f.Operator {
		case shared.FilterOperatorEqual:
			paginationFilter.Operator = pkg.OperatorEqual
		case shared.FilterOperatorNotEqual:
			paginationFilter.Operator = pkg.OperatorNotEqual
		case shared.FilterOperatorGreater:
			paginationFilter.Operator = pkg.OperatorGreater
		case shared.FilterOperatorLess:
			paginationFilter.Operator = pkg.OperatorLess
		case shared.FilterOperatorGreaterOrEqual:
			paginationFilter.Operator = pkg.OperatorGreaterOrEqual
		case shared.FilterOperatorLessOrEqual:
			paginationFilter.Operator = pkg.OperatorLessOrEqual
		case shared.FilterOperatorLike:
			paginationFilter.Operator = pkg.OperatorLike
		case shared.FilterOperatorIn:
			paginationFilter.Operator = pkg.OperatorIn
		case shared.FilterOperatorBetween:
			paginationFilter.Operator = pkg.OperatorBetween
		default:
			// For any other operators, use the string value
			paginationFilter.Operator = string(f.Operator)
		}

		result = append(result, paginationFilter)
	}

	return result
}
