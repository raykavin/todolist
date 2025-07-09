package valueobject

import (
	"errors"
	"testing"
)

func TestNewTaxID(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  error
		wantType TaxIDType
		wantNum  string
		wantFmt  string
	}{
		// Valid CPF cases
		{
			name:     "valid CPF with formatting",
			input:    "123.456.789-09",
			wantErr:  nil,
			wantType: CPF,
			wantNum:  "12345678909",
			wantFmt:  "123.456.789-09",
		},
		{
			name:     "valid CPF without formatting",
			input:    "12345678909",
			wantErr:  nil,
			wantType: CPF,
			wantNum:  "12345678909",
			wantFmt:  "123.456.789-09",
		},
		{
			name:     "valid CPF with spaces",
			input:    " 123.456.789-09 ",
			wantErr:  nil,
			wantType: CPF,
			wantNum:  "12345678909",
			wantFmt:  "123.456.789-09",
		},
		// Valid CNPJ cases
		{
			name:     "valid CNPJ with formatting",
			input:    "11.222.333/0001-81",
			wantErr:  nil,
			wantType: CNPJ,
			wantNum:  "11222333000181",
			wantFmt:  "11.222.333/0001-81",
		},
		{
			name:     "valid CNPJ without formatting",
			input:    "11222333000181",
			wantErr:  nil,
			wantType: CNPJ,
			wantNum:  "11222333000181",
			wantFmt:  "11.222.333/0001-81",
		},
		// Invalid cases
		{
			name:    "empty document",
			input:   "",
			wantErr: ErrEmptyTaxID,
		},
		{
			name:    "only spaces",
			input:   "   ",
			wantErr: ErrEmptyTaxID,
		},
		{
			name:    "invalid CPF - all same digits",
			input:   "111.111.111-11",
			wantErr: ErrInvalidTaxIDNumber,
		},
		{
			name:    "invalid CPF - wrong check digit",
			input:   "123.456.789-00",
			wantErr: ErrInvalidTaxIDNumber,
		},
		{
			name:    "invalid CNPJ - all same digits",
			input:   "11.111.111/1111-11",
			wantErr: ErrInvalidTaxIDNumber,
		},
		{
			name:    "invalid CNPJ - wrong check digit",
			input:   "11.222.333/0001-80",
			wantErr: ErrInvalidTaxIDNumber,
		},
		{
			name:    "invalid format - too short",
			input:   "123456",
			wantErr: ErrInvalidTaxIDNumber,
		},
		{
			name:    "invalid format - too long",
			input:   "123456789012345",
			wantErr: ErrInvalidTaxIDNumber,
		},
		{
			name:    "invalid format - letters",
			input:   "ABC.DEF.GHI-JK",
			wantErr: ErrInvalidTaxIDNumber,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := NewTaxID(tt.input)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewTaxID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return // Skip further checks if error is expected
			}

			if doc.Type() != tt.wantType {
				t.Errorf("Type() = %v, want %v", doc.Type(), tt.wantType)
			}

			if doc.Number() != tt.wantNum {
				t.Errorf("Number() = %v, want %v", doc.Number(), tt.wantNum)
			}

			if doc.Formatted() != tt.wantFmt {
				t.Errorf("Formatted() = %v, want %v", doc.Formatted(), tt.wantFmt)
			}

			if doc.String() != tt.wantFmt {
				t.Errorf("String() = %v, want %v", doc.String(), tt.wantFmt)
			}
		})
	}
}

func TestTaxID_Equals(t *testing.T) {
	doc1, _ := NewTaxID("123.456.789-09")
	doc2, _ := NewTaxID("12345678909") // Same CPF, different format
	doc3, _ := NewTaxID("987.654.321-00")

	if !doc1.Equals(doc2) {
		t.Error("Expected doc1 and doc2 to be equal")
	}

	if doc1.Equals(doc3) {
		t.Error("Expected doc1 and doc3 to be different")
	}
}

// TestRealTaxIDs tests with some real valid document numbers
func TestRealTaxIDs(t *testing.T) {
	validDocs := []struct {
		input    string
		docType  TaxIDType
		expected string
	}{
		// Real valid CPFs (test numbers)
		{"191.000.000-00", CPF, "191.000.000-00"},
		{"000.000.001-91", CPF, "000.000.001-91"},

		// Real valid CNPJs (test numbers)
		{"11.444.777/0001-61", CNPJ, "11.444.777/0001-61"},
		{"82.373.077/0001-71", CNPJ, "82.373.077/0001-71"},
	}

	for _, td := range validDocs {
		t.Run(td.input, func(t *testing.T) {
			doc, err := NewTaxID(td.input)
			if err != nil {
				t.Fatalf("Expected valid document, got error: %v", err)
			}

			if doc.Type() != td.docType {
				t.Errorf("Expected type %s, got %s", td.docType, doc.Type())
			}

			if doc.Formatted() != td.expected {
				t.Errorf("Expected formatted %s, got %s", td.expected, doc.Formatted())
			}
		})
	}
}

// BenchmarkNewTaxID benchmarks document creation
func BenchmarkNewTaxID(b *testing.B) {
	inputs := []string{
		"123.456.789-09",
		"12345678909",
		"11.222.333/0001-81",
		"11222333000181",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = NewTaxID(inputs[i%len(inputs)])
	}
}
