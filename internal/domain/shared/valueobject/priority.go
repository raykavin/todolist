package valueobject

import "errors"

// Priority represents the priority level
type Priority int8

const (
	PriorityLow Priority = iota + 1
	PriorityMedium
	PriorityHigh
	PriorityCritical
)

var (
	ErrInvalidPriority = errors.New("invalid priority")
)

// String returns the string representation of priority
func (p Priority) String() string {
	switch p {
	case PriorityLow:
		return "low"
	case PriorityMedium:
		return "medium"
	case PriorityHigh:
		return "high"
	case PriorityCritical:
		return "critical"
	default:
		return "unknown"
	}
}

// IsValid checks if the priority is valid
func (p Priority) IsValid() bool {
	return p >= PriorityLow && p <= PriorityCritical
}

// NewPriorityFromString creates a Priority from string
func NewPriorityFromString(s string) (Priority, error) {
	switch s {
	case "low":
		return PriorityLow, nil
	case "medium":
		return PriorityMedium, nil
	case "high":
		return PriorityHigh, nil
	case "critical":
		return PriorityCritical, nil
	default:
		return 0, ErrInvalidPriority
	}
}
