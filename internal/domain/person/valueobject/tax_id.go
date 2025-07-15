package valueobject

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var (
	// ErrInvalidTaxIDNumber is returned when the taxID number is invalid
	ErrInvalidTaxIDNumber = errors.New("invalid taxID number")
	// ErrInvalidTaxIDType is returned when the taxID type cannot be determined
	ErrInvalidTaxIDType = errors.New("invalid taxID type")
	// ErrEmptyTaxID is returned when the taxID int64 is empty
	ErrEmptyTaxID = errors.New("taxID cannot be empty")
)

var (
	// cpfRegexp matches CPF format with or without formatting
	// Valid formats: 123.456.789-10 or 12345678910
	cpfRegexp = regexp.MustCompile(`^\d{3}\.?\d{3}\.?\d{3}-?\d{2}$`)

	// cnpjRegexp matches CNPJ format with or without formatting
	// Valid formats: 12.345.678/0001-90 or 12345678000190
	cnpjRegexp = regexp.MustCompile(`^\d{2}\.?\d{3}\.?\d{3}\/?(:?\d{3}[1-9]|\d{2}[1-9]\d|\d[1-9]\d{2}|[1-9]\d{3})-?\d{2}$`)
)

// TaxIDType represents the type of Brazilian taxID
type TaxIDType string

const (
	CPF  TaxIDType = "CPF"  // CPF represents a Brazilian individual taxpayer registry identification
	CNPJ TaxIDType = "CNPJ" // CNPJ represents a Brazilian company taxpayer registry identification
)

// TaxID is a value object representing a Brazilian taxID (CPF or CNPJ).
type TaxID struct {
	number    string    // The normalized taxID number (only digits)
	taxIDType TaxIDType // The type of taxID (CPF or CNPJ)
	formatted string    // The formatted version of the taxID
}

// NewTaxID creates a new TaxID value object.
// It accepts a taxID int64 in various formats:
//   - CPF: "123.456.789-10" or "12345678910"
//   - CNPJ: "12.345.678/0001-90" or "12345678000190"
//
// The function validates the taxID according to Brazilian rules,
// including check digit validation.
func NewTaxID(taxID string) (TaxID, error) {
	taxID = strings.TrimSpace(taxID)
	if taxID == "" {
		return TaxID{}, ErrEmptyTaxID
	}

	doc := TaxID{}

	// Clean and store the normalized number
	doc.number = extractDigits(taxID)

	// Determine taxID type and validate
	if err := doc.determineTypeAndValidate(taxID); err != nil {
		return TaxID{}, err
	}

	// Store formatted version
	doc.formatted = doc.format()

	return doc, nil
}

// Number returns the taxID number without formatting (only digits)
func (d TaxID) Number() string {
	return d.number
}

// Type returns the taxID type (CPF or CNPJ)
func (d TaxID) Type() TaxIDType {
	return d.taxIDType
}

// Formatted returns the taxID in its standard Brazilian format
func (d TaxID) Formatted() string {
	return d.formatted
}

// FormatMasked returns a partially masked version of the TaxID.
// For CPF: XXX.***.***-XX
// For CNPJ: XX.XXX.***/XXXX-XX
func (d TaxID) FormatMasked() string {
	switch d.taxIDType {
	case CPF:
		if len(d.number) == 11 {
			return d.number[:3] + ".***.***-" + d.number[9:]
		}
	case CNPJ:
		if len(d.number) == 14 {
			return d.number[:2] + "." + d.number[2:5] + ".***/" + d.number[8:12] + "-" + d.number[12:]
		}
	}
	return d.number
}

// String implements the Stringer interface, returning the formatted taxID
func (d TaxID) String() string {
	return d.formatted
}

// Equals compares two taxIDs for equality based on their number
func (d TaxID) Equals(other TaxID) bool {
	return d.number == other.number
}

// determineTypeAndValidate determines the taxID type and validates it
func (d *TaxID) determineTypeAndValidate(original string) error {
	// Try CPF first
	if cpfRegexp.MatchString(original) && d.isValidCPF() {
		d.taxIDType = CPF
		return nil
	}

	// Try CNPJ
	if cnpjRegexp.MatchString(original) && d.isValidCNPJ() {
		d.taxIDType = CNPJ
		return nil
	}

	return ErrInvalidTaxIDNumber
}

// isValidCPF validates a CPF number according to Brazilian rules
func (d *TaxID) isValidCPF() bool {
	if len(d.number) != 11 {
		return false
	}

	// Check if all digits are the same (invalid CPFs like 111.111.111-11)
	if hasAllSameDigits(d.number) {
		return false
	}

	// Validate check digits
	return d.validateCheckDigits(9, 10)
}

// isValidCNPJ validates a CNPJ number according to Brazilian rules
func (d *TaxID) isValidCNPJ() bool {
	if len(d.number) != 14 {
		return false
	}

	// Check if all digits are the same
	if hasAllSameDigits(d.number) {
		return false
	}

	// Validate check digits
	return d.validateCheckDigits(12, 5)
}

// validateCheckDigits validates the check digits of a taxID
func (d *TaxID) validateCheckDigits(size, startPos int) bool {
	// Get the base number without check digits
	base := d.number[:size]

	// Calculate first check digit
	digit1 := calculateCheckDigit(base, startPos)
	base = base + digit1

	// Calculate second check digit
	digit2 := calculateCheckDigit(base, startPos+1)

	// Compare with actual check digits
	expected := base + digit2
	return d.number == expected
}

// calculateCheckDigit calculates a check digit for Brazilian taxIDs
func calculateCheckDigit(base string, position int) string {
	sum := 0
	pos := position

	for _, r := range base {
		digit := int(r - '0')
		sum += digit * pos
		pos--

		// For CNPJ calculation, reset position to 9 when it reaches 1
		if pos < 2 {
			pos = 9
		}
	}

	remainder := sum % 11
	if remainder < 2 {
		return "0"
	}

	return strconv.Itoa(11 - remainder)
}

// format returns the taxID in its standard Brazilian format
func (d *TaxID) format() string {
	switch d.taxIDType {
	case CPF:
		// Format: XXX.XXX.XXX-XX
		if len(d.number) == 11 {
			return d.number[:3] + "." + d.number[3:6] + "." +
				d.number[6:9] + "-" + d.number[9:]
		}
	case CNPJ:
		// Format: XX.XXX.XXX/XXXX-XX
		if len(d.number) == 14 {
			return d.number[:2] + "." + d.number[2:5] + "." +
				d.number[5:8] + "/" + d.number[8:12] + "-" + d.number[12:]
		}
	}
	return d.number
}

// extractDigits extracts only digit characters from a string
func extractDigits(s string) string {
	var result strings.Builder
	result.Grow(len(s))

	for _, r := range s {
		if unicode.IsDigit(r) {
			result.WriteRune(r)
		}
	}

	return result.String()
}

// hasAllSameDigits checks if all digits in a string are the same
func hasAllSameDigits(s string) bool {
	if len(s) == 0 {
		return false
	}

	first := s[0]
	for i := 1; i < len(s); i++ {
		if s[i] != first {
			return false
		}
	}

	return true
}
