package http

import (
	"encoding/json"
	"strconv"
	"strings"
	"todolist/pkg"

	"github.com/gin-gonic/gin"
)

// PaginationAdapter handles HTTP-specific pagination conversion
type PaginationAdapter struct{}

// NewPaginationAdapter creates a new HTTP pagination adapter
func NewPaginationAdapter() *PaginationAdapter {
	return &PaginationAdapter{}
}

// ParseFromQuery parses pagination parameters from query string
func (a *PaginationAdapter) ParseFromQuery(c *gin.Context) *pkg.Request {
	req := pkg.DefaultRequest()

	// Parse page
	if page := c.Query("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil {
			req.Page = p
		}
	}

	// Parse size
	if size := c.Query("size"); size != "" {
		if s, err := strconv.Atoi(size); err == nil {
			req.Size = s
		}
	}

	// Parse sorting
	req.SortFields = a.parseSortFromQuery(c)

	// Parse global filter
	if search := c.Query("search"); search != "" {
		req.GlobalFilter = search
	}

	// Parse filters
	req.Filters = a.parseFiltersFromQuery(c)

	// Alternative: parse complete pagination from 'filter' query param as JSON
	if filterJSON := c.Query("filter"); filterJSON != "" {
		if parsed, err := pkg.ParseFromJSON(filterJSON); err == nil {
			return parsed
		}
	}

	return req
}

// ParseFromJSON parses pagination parameters from request body
func (a *PaginationAdapter) ParseFromJSON(c *gin.Context) (*pkg.Request, error) {
	var req pkg.Request

	// Try to bind JSON body
	if err := c.ShouldBindJSON(&req); err != nil {
		// If no body or error, return default
		return pkg.DefaultRequest(), nil
	}

	// Validate the request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	return &req, nil
}

// parseSortFromQuery parses sort parameters from query string
func (a *PaginationAdapter) parseSortFromQuery(c *gin.Context) []pkg.SortField {
	var sortFields []pkg.SortField

	if sort := c.Query("sort"); sort != "" {
		parts := strings.Split(sort, ",")
		for _, part := range parts {
			fieldParts := strings.Split(strings.TrimSpace(part), ":")
			if len(fieldParts) >= 1 {
				field := fieldParts[0]
				desc := false
				if len(fieldParts) >= 2 {
					desc = strings.ToLower(fieldParts[1]) == "desc"
				}
				sortFields = append(sortFields, pkg.SortField{
					Field: field,
					Desc:  desc,
				})
			}
		}
	}

	// Alternative format: sort_by=field&order=desc
	if sortBy := c.Query("sort_by"); sortBy != "" {
		order := c.DefaultQuery("order", "asc")
		sortFields = append(sortFields, pkg.SortField{
			Field: sortBy,
			Desc:  strings.ToLower(order) == "desc",
		})
	}

	return sortFields
}

// parseFiltersFromQuery parses filter parameters from query string
func (a *PaginationAdapter) parseFiltersFromQuery(c *gin.Context) []pkg.Filter {
	var filters []pkg.Filter

	// filter[field][op]=value
	// Example: filter[name][like]=John&filter[age][gte]=18
	for key, values := range c.Request.URL.Query() {
		if strings.HasPrefix(key, "filter[") && len(values) > 0 {
			parts := strings.Split(strings.TrimPrefix(key, "filter["), "][")
			if len(parts) == 2 {
				field := strings.TrimSuffix(parts[0], "]")
				op := strings.TrimSuffix(parts[1], "]")

				filter := pkg.Filter{
					Field:    field,
					Operator: op,
					Value:    a.parseFilterValue(values[0]),
				}

				filters = append(filters, filter)
			}
		}
	}

	// Simple key-value with operator suffix
	// Example: name_like=John&age_gte=18
	for key, values := range c.Request.URL.Query() {
		if len(values) == 0 || strings.HasPrefix(key, "filter") {
			continue
		}

		// Skip pagination params
		if key == "page" || key == "size" || key == "sort" || key == "sort_by" ||
			key == "order" || key == "search" || key == "filter" {
			continue
		}

		// Check for operator suffix
		field, op := a.parseFieldOperator(key)
		if field != "" {
			filter := pkg.Filter{
				Field:    field,
				Operator: op,
				Value:    a.parseFilterValue(values[0]),
			}
			filters = append(filters, filter)
		}
	}

	return filters
}

// parseFieldOperator extracts field and operator from key
func (a *PaginationAdapter) parseFieldOperator(key string) (field string, operator string) {
	// Known operator suffixes
	operatorSuffixes := map[string]string{
		"_eq":       "=",
		"_ne":       "!=",
		"_gt":       ">",
		"_gte":      ">=",
		"_lt":       "<",
		"_lte":      "<=",
		"_like":     "like",
		"_in":       "in",
		"_between":  "between",
		"_null":     "is_null",
		"_not_null": "is_not_null",
	}

	for suffix, op := range operatorSuffixes {
		if strings.HasSuffix(key, suffix) {
			return strings.TrimSuffix(key, suffix), op
		}
	}

	// Default to equal operator
	return key, "="
}

// parseFilterValue attempts to parse the string value to appropriate type
func (a *PaginationAdapter) parseFilterValue(value string) any {
	// Try to parse as number
	if i, err := strconv.Atoi(value); err == nil {
		return i
	}

	if f, err := strconv.ParseFloat(value, 64); err == nil {
		return f
	}

	// Try to parse as boolean
	if b, err := strconv.ParseBool(value); err == nil {
		return b
	}

	// Try to parse as JSON array (for IN operator)
	if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
		var arr []any
		if err := json.Unmarshal([]byte(value), &arr); err == nil {
			return arr
		}
	}

	// Check for comma-separated values (alternative for IN operator)
	if strings.Contains(value, ",") {
		parts := strings.Split(value, ",")
		result := make([]any, len(parts))
		for i, part := range parts {
			result[i] = a.parseFilterValue(strings.TrimSpace(part))
		}
		return result
	}

	// Return as string
	return value
}

// WriteResponse writes paginated response with proper headers
func (a *PaginationAdapter) WriteResponse(c *gin.Context, data any, result *pkg.Result) {
	// Set pagination headers
	c.Header("X-Total-Count", strconv.FormatInt(result.TotalRows, 10))
	c.Header("X-Total-Pages", strconv.Itoa(result.TotalPages))
	c.Header("X-Current-Page", strconv.Itoa(result.CurrentPage))
	c.Header("X-Page-Size", strconv.Itoa(result.PageSize))

	// Build response
	response := pkg.NewResponse(data, result)

	c.JSON(200, response)
}

// BuildFilterFromParams builds filters from common search parameters
func (a *PaginationAdapter) BuildFilterFromParams(c *gin.Context, fieldMappings map[string]string) []pkg.Filter {
	var filters []pkg.Filter

	for param, field := range fieldMappings {
		if value := c.Query(param); value != "" {
			filters = append(filters, pkg.Filter{
				Field:    field,
				Operator: pkg.OperatorEqual,
				Value:    a.parseFilterValue(value),
			})
		}
	}

	return filters
}

// Middleware creates a pagination middleware that parses and validates pagination
func (a *PaginationAdapter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req *pkg.Request

		// For POST/PUT/PATCH, try to parse from body
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			if parsed, err := a.ParseFromJSON(c); err == nil {
				req = parsed
			}
		}

		// For all methods, also check query parameters
		if req == nil {
			req = a.ParseFromQuery(c)
		}

		// Validate pagination
		if err := req.Validate(); err != nil {
			c.JSON(400, gin.H{"error": "Invalid pagination parameters", "details": err.Error()})
			c.Abort()
			return
		}

		// Store in context
		c.Set("pagination", req)

		c.Next()
	}
}

// GetPaginationFromContext retrieves pagination from gin context
func GetPaginationFromContext(c *gin.Context) (*pkg.Request, bool) {
	if value, exists := c.Get("pagination"); exists {
		if req, ok := value.(*pkg.Request); ok {
			return req, true
		}
	}
	return nil, false
}
