package pkg

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
)

// Request represents generic pagination request parameters
type Request struct {
	Page         int         `json:"page,omitempty"`
	Size         int         `json:"size,omitempty"`
	Filters      []Filter    `json:"filters,omitempty"`
	GlobalFilter string      `json:"global_filter,omitempty"`
	SortFields   []SortField `json:"sort_fields,omitempty"`
}

// Result represents generic pagination result
type Result struct {
	TotalRows   int64 `json:"total_rows"`
	TotalPages  int   `json:"total_pages"`
	CurrentPage int   `json:"current_page"`
	PageSize    int   `json:"page_size"`
}

// Filter represents a generic filter condition
type Filter struct {
	Field    string      `json:"field,omitempty"`
	Operator string      `json:"op,omitempty"`
	Value    interface{} `json:"value,omitempty"`
}

// SortField represents a field to sort by
type SortField struct {
	Field string `json:"field,omitempty"`
	Desc  bool   `json:"desc,omitempty"`
}

// String returns the sort field as a SQL-compatible string
func (s SortField) String() string {
	if s.Desc {
		return fmt.Sprintf("%s DESC", s.Field)
	}
	return fmt.Sprintf("%s ASC", s.Field)
}

// GetOffset calculates the offset for the current page
func (r *Request) GetOffset() int {
	return (r.GetPage() - 1) * r.GetSize()
}

// GetSize returns the page size with fallback to default
func (r *Request) GetSize() int {
	if r.Size <= 0 {
		return DefaultPageSize
	}
	if r.Size > MaxPageSize {
		return MaxPageSize
	}
	return r.Size
}

// GetPage returns the current page with fallback to default
func (r *Request) GetPage() int {
	if r.Page <= 0 {
		return DefaultPage
	}
	return r.Page
}

// CalculateTotalPages calculates total pages based on total rows
func CalculateTotalPages(totalRows int64, pageSize int) int {
	if pageSize <= 0 {
		return 0
	}
	return int(math.Ceil(float64(totalRows) / float64(pageSize)))
}

// Response represents the complete pagination response
type Response struct {
	Data       interface{} `json:"data"`
	Pagination Meta        `json:"pagination"`
}

// Meta represents pagination metadata
type Meta struct {
	Page       int   `json:"page"`
	Size       int   `json:"size"`
	TotalRows  int64 `json:"total_rows"`
	TotalPages int   `json:"total_pages"`
}

// Constants for pagination defaults and limits
const (
	DefaultPage     = 1
	DefaultPageSize = 25
	MaxPageSize     = 100
	MinPageSize     = 1
)

// DefaultRequest returns a pagination request with default values
func DefaultRequest() *Request {
	return &Request{
		Page: DefaultPage,
		Size: DefaultPageSize,
	}
}

// NewResponse creates a new pagination response
func NewResponse(data interface{}, result *Result) Response {
	return Response{
		Data: data,
		Pagination: Meta{
			Page:       result.CurrentPage,
			Size:       result.PageSize,
			TotalRows:  result.TotalRows,
			TotalPages: result.TotalPages,
		},
	}
}

// Operators for filtering
const (
	OperatorEqual          = "="
	OperatorNotEqual       = "!="
	OperatorGreater        = ">"
	OperatorLess           = "<"
	OperatorGreaterOrEqual = ">="
	OperatorLessOrEqual    = "<="
	OperatorLike           = "like"
	OperatorIn             = "in"
	OperatorNotIn          = "not_in"
	OperatorBetween        = "between"
	OperatorRegex          = "regex"
	OperatorIsNull         = "is_null"
	OperatorIsNotNull      = "is_not_null"
)

// NormalizeOperator converts various operator representations to standard form
func NormalizeOperator(op string) string {
	normalized := strings.ToLower(strings.TrimSpace(op))

	switch normalized {
	case "=", "eq", "equal", "equals":
		return OperatorEqual
	case "!=", "<>", "ne", "neq", "not_equal":
		return OperatorNotEqual
	case ">", "gt":
		return OperatorGreater
	case "<", "lt":
		return OperatorLess
	case ">=", "gte":
		return OperatorGreaterOrEqual
	case "<=", "lte":
		return OperatorLessOrEqual
	case "like", "contains":
		return OperatorLike
	case "in":
		return OperatorIn
	case "not_in", "notin":
		return OperatorNotIn
	case "between":
		return OperatorBetween
	case "regex", "regexp":
		return OperatorRegex
	case "is_null", "isnull", "null":
		return OperatorIsNull
	case "is_not_null", "isnotnull", "not_null":
		return OperatorIsNotNull
	default:
		return normalized
	}
}

// SortOrder represents sort direction
type SortOrder string

const (
	SortOrderAsc  SortOrder = "ASC"
	SortOrderDesc SortOrder = "DESC"
)

// ParseSortOrder parses string to SortOrder
func ParseSortOrder(s string) SortOrder {
	upper := strings.ToUpper(strings.TrimSpace(s))
	if upper == "DESC" {
		return SortOrderDesc
	}
	return SortOrderAsc
}

// Validation errors
var (
	ErrInvalidPage     = fmt.Errorf("page must be greater than 0")
	ErrInvalidPageSize = fmt.Errorf("page size must be between %d and %d", MinPageSize, MaxPageSize)
	ErrInvalidOperator = fmt.Errorf("invalid filter operator")
)

// Validate validates the pagination request
func (r *Request) Validate() error {
	if r.Page < 0 {
		return ErrInvalidPage
	}

	if r.Size < 0 || (r.Size > MaxPageSize && r.Size != 0) {
		return ErrInvalidPageSize
	}

	return nil
}

// Builder provides a fluent interface for building pagination requests
type Builder struct {
	request *Request
}

// NewBuilder creates a new pagination request builder
func NewBuilder() *Builder {
	return &Builder{
		request: DefaultRequest(),
	}
}

// WithPage sets the page number
func (b *Builder) WithPage(page int) *Builder {
	b.request.Page = page
	return b
}

// WithSize sets the page size
func (b *Builder) WithSize(size int) *Builder {
	b.request.Size = size
	return b
}

// WithFilter adds a filter
func (b *Builder) WithFilter(field, operator string, value interface{}) *Builder {
	b.request.Filters = append(b.request.Filters, Filter{
		Field:    field,
		Operator: NormalizeOperator(operator),
		Value:    value,
	})
	return b
}

// WithSort adds a sort field
func (b *Builder) WithSort(field string, desc bool) *Builder {
	b.request.SortFields = append(b.request.SortFields, SortField{
		Field: field,
		Desc:  desc,
	})
	return b
}

// WithGlobalFilter sets the global filter
func (b *Builder) WithGlobalFilter(filter string) *Builder {
	b.request.GlobalFilter = filter
	return b
}

// Build returns the pagination request
func (b *Builder) Build() (*Request, error) {
	if err := b.request.Validate(); err != nil {
		return nil, err
	}
	return b.request, nil
}

// JSON helpers

// ParseFromJSON parses pagination request from JSON string
func ParseFromJSON(jsonStr string) (*Request, error) {
	if jsonStr == "" {
		return DefaultRequest(), nil
	}

	var req Request
	if err := json.Unmarshal([]byte(jsonStr), &req); err != nil {
		return nil, fmt.Errorf("failed to parse pagination JSON: %w", err)
	}

	if err := req.Validate(); err != nil {
		return nil, err
	}

	return &req, nil
}
