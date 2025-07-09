package shared

import (
	"errors"
)

// Common repository errors
var (
	ErrNotFound         = errors.New("entity not found")
	ErrDuplicateEntry   = errors.New("duplicate entry")
	ErrInvalidOperation = errors.New("invalid operation")
	ErrOptimisticLock   = errors.New("optimistic lock error")
)

// QueryOptions represents query options
type QueryOptions struct {
	Limit     int
	Offset    int
	OrderBy   string
	OrderDesc bool
}

// Filter represents a generic filter
type Filter struct {
	Field    string
	Operator FilterOperator
	Value    any
}

// FilterOperator types of operators
type FilterOperator string

const (
	FilterOperatorEqual          FilterOperator = "="
	FilterOperatorNotEqual       FilterOperator = "!="
	FilterOperatorGreater        FilterOperator = ">"
	FilterOperatorLess           FilterOperator = "<"
	FilterOperatorGreaterOrEqual FilterOperator = ">="
	FilterOperatorLessOrEqual    FilterOperator = "<="
	FilterOperatorLike           FilterOperator = "LIKE"
	FilterOperatorIn             FilterOperator = "IN"
	FilterOperatorBetween        FilterOperator = "BETWEEN"
)

// Transaction interface for transactions
type Transaction interface {
	Commit() error
	Rollback() error
}
