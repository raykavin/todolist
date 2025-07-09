package valueobject

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNewDate(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		wantErr bool
		check   func(Date) bool
	}{
		{
			name:    "from time.Time",
			input:   time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
			wantErr: false,
			check: func(d Date) bool {
				return d.Year() == 2024 && d.Month() == 1 && d.Day() == 15
			},
		},
		{
			name:    "from ISO string",
			input:   "2024-01-15",
			wantErr: false,
			check: func(d Date) bool {
				return d.String() == "2024-01-15"
			},
		},
		{
			name:    "from DD/MM/YYYY",
			input:   "15/01/2024",
			wantErr: false,
			check: func(d Date) bool {
				return d.Year() == 2024 && d.Month() == 1 && d.Day() == 15
			},
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "invalid string",
			input:   "not-a-date",
			wantErr: true,
		},
		{
			name:    "nil pointer",
			input:   (*time.Time)(nil),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			date, err := NewDate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.check != nil && !tt.check(date) {
				t.Errorf("Date check failed for input %v", tt.input)
			}
		})
	}
}

func TestNewDateFromComponents(t *testing.T) {
	date := NewDateFromComponents(2024, 12, 25)

	if date.Year() != 2024 {
		t.Errorf("Year() = %v, want 2024", date.Year())
	}
	if date.Month() != 12 {
		t.Errorf("Month() = %v, want 12", date.Month())
	}
	if date.Day() != 25 {
		t.Errorf("Day() = %v, want 25", date.Day())
	}
}

func TestDate_Comparisons(t *testing.T) {
	date1 := NewDateFromComponents(2024, 1, 15)
	date2 := NewDateFromComponents(2024, 1, 15)
	date3 := NewDateFromComponents(2024, 1, 16)

	if !date1.Equals(date2) {
		t.Error("Expected date1 and date2 to be equal")
	}

	if date1.Equals(date3) {
		t.Error("Expected date1 and date3 to be different")
	}

	if !date1.Before(date3) {
		t.Error("Expected date1 to be before date3")
	}

	if !date3.After(date1) {
		t.Error("Expected date3 to be after date1")
	}
}

func TestDate_Arithmetic(t *testing.T) {
	base := NewDateFromComponents(2024, 1, 15)

	tests := []struct {
		name     string
		op       func() Date
		expected string
	}{
		{
			name:     "add 10 days",
			op:       func() Date { return base.AddDays(10) },
			expected: "2024-01-25",
		},
		{
			name:     "subtract 5 days",
			op:       func() Date { return base.AddDays(-5) },
			expected: "2024-01-10",
		},
		{
			name:     "add 2 months",
			op:       func() Date { return base.AddMonths(2) },
			expected: "2024-03-15",
		},
		{
			name:     "add 1 year",
			op:       func() Date { return base.AddYears(1) },
			expected: "2025-01-15",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.op()
			if result.String() != tt.expected {
				t.Errorf("got %v, want %v", result.String(), tt.expected)
			}
		})
	}
}

func TestDate_DaysBetween(t *testing.T) {
	date1 := NewDateFromComponents(2024, 1, 1)
	date2 := NewDateFromComponents(2024, 1, 11)

	days := date1.DaysBetween(date2)
	if days != 10 {
		t.Errorf("DaysBetween() = %v, want 10", days)
	}

	days = date2.DaysBetween(date1)
	if days != -10 {
		t.Errorf("DaysBetween() = %v, want -10", days)
	}
}

func TestDate_IsWorkday(t *testing.T) {
	tests := []struct {
		date      Date
		isWorkday bool
	}{
		{NewDateFromComponents(2024, 1, 15), true},  // Monday
		{NewDateFromComponents(2024, 1, 16), true},  // Tuesday
		{NewDateFromComponents(2024, 1, 17), true},  // Wednesday
		{NewDateFromComponents(2024, 1, 18), true},  // Thursday
		{NewDateFromComponents(2024, 1, 19), true},  // Friday
		{NewDateFromComponents(2024, 1, 20), false}, // Saturday
		{NewDateFromComponents(2024, 1, 21), false}, // Sunday
	}

	for _, tt := range tests {
		t.Run(tt.date.Weekday().String(), func(t *testing.T) {
			if tt.date.IsWorkday() != tt.isWorkday {
				t.Errorf("IsWorkday() = %v, want %v", tt.date.IsWorkday(), tt.isWorkday)
			}
		})
	}
}

func TestDate_JSON(t *testing.T) {
	date := NewDateFromComponents(2024, 1, 15)

	// Test marshaling
	data, err := json.Marshal(date)
	if err != nil {
		t.Fatalf("Failed to marshal date: %v", err)
	}

	expected := `"2024-01-15"`
	if string(data) != expected {
		t.Errorf("Marshaled = %s, want %s", string(data), expected)
	}

	// Test unmarshaling
	var unmarshaled Date
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal date: %v", err)
	}

	if !date.Equals(unmarshaled) {
		t.Errorf("Unmarshaled date doesn't match original")
	}

	// Test zero value
	var zero Date
	data, err = json.Marshal(zero)
	if err != nil {
		t.Fatalf("Failed to marshal zero date: %v", err)
	}
	if string(data) != "null" {
		t.Errorf("Zero date marshaled = %s, want null", string(data))
	}
}

func TestDate_Format(t *testing.T) {
	date := NewDateFromComponents(2024, 1, 15)

	tests := []struct {
		layout   string
		expected string
	}{
		{"2006-01-02", "2024-01-15"},
		{"02/01/2006", "15/01/2024"},
		{"Jan 2, 2006", "Jan 15, 2024"},
		{"Monday, January 2, 2006", "Monday, January 15, 2024"},
	}

	for _, tt := range tests {
		t.Run(tt.layout, func(t *testing.T) {
			result := date.Format(tt.layout)
			if result != tt.expected {
				t.Errorf("Format() = %v, want %v", result, tt.expected)
			}
		})
	}
}
