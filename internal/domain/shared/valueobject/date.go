package valueobject

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

var (
	ErrInvalidDate = errors.New("invalid date")
	ErrEmptyDate   = errors.New("date cannot be empty")
)

// Date represents a date without time components
type Date struct {
	value time.Time
}

// NewDate creates a new Date from various input types
func NewDate(input any) (Date, error) {
	switch v := input.(type) {
	case time.Time:
		return Date{value: truncateTime(v)}, nil
	case string:
		if v == "" {
			return Date{}, ErrEmptyDate
		}
		return parseDate(v)
	case *time.Time:
		if v == nil {
			return Date{}, ErrEmptyDate
		}
		return Date{value: truncateTime(*v)}, nil
	default:
		return Date{}, fmt.Errorf("%w: unsupported type %T", ErrInvalidDate, input)
	}
}

// NewDateFromComponents creates a Date from year, month, and day
func NewDateFromComponents(year, month, day int) Date {
	return Date{value: time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)}
}

// Today returns the current date
func Today() Date {
	return Date{value: truncateTime(time.Now())}
}

func (d Date) Year() int {
	return d.value.Year()
}

func (d Date) Month() time.Month {
	return d.value.Month()
}

func (d Date) Day() int {
	return d.value.Day()
}

func (d Date) Weekday() time.Weekday {
	return d.value.Weekday()
}

func (d Date) Time() time.Time {
	return d.value
}

func (d Date) String() string {
	return d.value.Format("2006-01-02")
}

func (d Date) Format(layout string) string {
	return d.value.Format(layout)
}

func (d Date) Equals(other Date) bool {
	return d.value.Equal(other.value)
}

func (d Date) Before(other Date) bool {
	return d.value.Before(other.value)
}

func (d Date) After(other Date) bool {
	return d.value.After(other.value)
}

func (d Date) AddDays(days int) Date {
	return Date{value: d.value.AddDate(0, 0, days)}
}

func (d Date) AddMonths(months int) Date {
	return Date{value: d.value.AddDate(0, months, 0)}
}

func (d Date) AddYears(years int) Date {
	return Date{value: d.value.AddDate(years, 0, 0)}
}

func (d Date) DaysBetween(other Date) int {
	duration := other.value.Sub(d.value)
	return int(duration.Hours() / 24)
}

func (d Date) IsZero() bool {
	return d.value.IsZero()
}

// IsWorkday returns true if the date is Monday through Friday
func (d Date) IsWorkday() bool {
	wd := d.value.Weekday()
	return wd >= time.Monday && wd <= time.Friday
}

// MarshalJSON implements json.Marshaler
func (d Date) MarshalJSON() ([]byte, error) {
	if d.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(d.String())
}

// UnmarshalJSON implements json.Unmarshaler
func (d *Date) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	if str == "" {
		return nil
	}

	parsed, err := parseDate(str)
	if err != nil {
		return err
	}

	*d = parsed
	return nil
}

// Value implements driver.Valuer for database operations
func (d Date) Value() (driver.Value, error) {
	if d.IsZero() {
		return nil, nil
	}
	return d.value, nil
}

// Scan implements sql.Scanner for database operations
func (d *Date) Scan(value any) error {
	if value == nil {
		*d = Date{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		*d = Date{value: truncateTime(v)}
		return nil
	case string:
		parsed, err := parseDate(v)
		if err != nil {
			return err
		}
		*d = parsed
		return nil
	default:
		return fmt.Errorf("cannot scan %T into Date", value)
	}
}

func truncateTime(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func parseDate(str string) (Date, error) {
	layouts := []string{
		"2006-01-02",
		"02/01/2006",
		"01/02/2006",
		"2006/01/02",
		"02-01-2006",
		"01-02-2006",
		time.RFC3339,
		time.RFC3339Nano,
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, str); err == nil {
			return Date{value: truncateTime(t)}, nil
		}
	}

	// Try parsing with Go's smart date parser
	if t, err := time.Parse(time.DateOnly, str); err == nil {
		return Date{value: truncateTime(t)}, nil
	}

	return Date{}, fmt.Errorf("%w: unable to parse %q", ErrInvalidDate, str)
}
