package database

import (
	"fmt"
	"strings"
	"todolist/internal/domain/shared"

	"gorm.io/gorm"
)

// ApplyFilters applies domain filters to GORM query
func ApplyFilters(query *gorm.DB, filters []shared.Filter) *gorm.DB {
	for _, filter := range filters {
		query = applyFilter(query, filter)
	}
	return query
}

// applyFilter applies a single filter to the query
func applyFilter(query *gorm.DB, filter shared.Filter) *gorm.DB {
	switch filter.Operator {
	case shared.FilterOperatorEqual:
		return query.Where(fmt.Sprintf("%s = ?", filter.Field), filter.Value)
	case shared.FilterOperatorNotEqual:
		return query.Where(fmt.Sprintf("%s != ?", filter.Field), filter.Value)
	case shared.FilterOperatorGreater:
		return query.Where(fmt.Sprintf("%s > ?", filter.Field), filter.Value)
	case shared.FilterOperatorLess:
		return query.Where(fmt.Sprintf("%s < ?", filter.Field), filter.Value)
	case shared.FilterOperatorGreaterOrEqual:
		return query.Where(fmt.Sprintf("%s >= ?", filter.Field), filter.Value)
	case shared.FilterOperatorLessOrEqual:
		return query.Where(fmt.Sprintf("%s <= ?", filter.Field), filter.Value)
	case shared.FilterOperatorLike:
		return query.Where(fmt.Sprintf("%s LIKE ?", filter.Field), fmt.Sprintf("%%%v%%", filter.Value))
	case shared.FilterOperatorIn:
		return query.Where(fmt.Sprintf("%s IN ?", filter.Field), filter.Value)
	case shared.FilterOperatorBetween:
		if values, ok := filter.Value.([]any); ok && len(values) == 2 {
			return query.Where(fmt.Sprintf("%s BETWEEN ? AND ?", filter.Field), values[0], values[1])
		}
	}
	return query
}

// ApplyQueryOptions applies pagination and sorting options
func ApplyQueryOptions(query *gorm.DB, options shared.QueryOptions) *gorm.DB {
	// Apply sorting
	if options.OrderBy != "" {
		order := options.OrderBy
		if options.OrderDesc {
			order += " DESC"
		} else {
			order += " ASC"
		}
		query = query.Order(order)
	}

	// Apply pagination
	if options.Limit > 0 {
		query = query.Limit(options.Limit)
	}
	if options.Offset > 0 {
		query = query.Offset(options.Offset)
	}

	return query
}

// BuildSearchQuery builds a search query for multiple fields
func BuildSearchQuery(query *gorm.DB, searchTerm string, fields ...string) *gorm.DB {
	if searchTerm == "" || len(fields) == 0 {
		return query
	}

	searchTerm = strings.TrimSpace(searchTerm)
	conditions := make([]string, len(fields))
	values := make([]any, len(fields))

	for i, field := range fields {
		conditions[i] = fmt.Sprintf("%s ILIKE ?", field)
		values[i] = fmt.Sprintf("%%%s%%", searchTerm)
	}

	whereClause := strings.Join(conditions, " OR ")
	return query.Where(whereClause, values...)
}

// Preload applies preloading with conditions
func Preload(query *gorm.DB, associations ...string) *gorm.DB {
	for _, association := range associations {
		query = query.Preload(association)
	}
	return query
}

// SoftDelete applies soft delete filter
func SoftDelete(query *gorm.DB) *gorm.DB {
	return query.Where("deleted_at IS NULL")
}

// Paginate creates a paginated response
type PaginatedResult struct {
	Data       any   `json:"data"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}

// Paginate applies pagination and returns paginated result
func Paginate(db *gorm.DB, page, pageSize int, result any) (*PaginatedResult, error) {
	var total int64

	// Count total records
	countQuery := db.Session(&gorm.Session{})
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, err
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Get paginated data
	if err := db.Limit(pageSize).Offset(offset).Find(result).Error; err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &PaginatedResult{
		Data:       result,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// BatchInsert performs batch insert with chunk size
func BatchInsert(db *gorm.DB, records any, chunkSize int) error {
	return db.CreateInBatches(records, chunkSize).Error
}

// Transaction wraps a function in a database transaction
func Transaction(db *gorm.DB, fn func(*gorm.DB) error) error {
	return db.Transaction(fn)
}
