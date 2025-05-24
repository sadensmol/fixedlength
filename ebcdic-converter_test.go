package fixedlength

import (
	"testing"
)

func TestEbcdicToAsciiNumber(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		decimalPlaces int
		expected      string
	}{
		{
			name:     "Empty string",
			input:    "",
			expected: "0",
		},
		{
			name:          "Empty string",
			input:         "",
			decimalPlaces: 2,
			expected:      "0.00",
		},
		{
			name:     "Normal number",
			input:    "12345",
			expected: "12345",
		},
		{
			name:          "Normal number with 2 decimal places",
			input:         "12345",
			decimalPlaces: 2,
			expected:      "123.45",
		},
		{
			name:     "Negative number with J",
			input:    "0000J",
			expected: "-1",
		},
		{
			name:          "Negative number with J and 2 decimal places",
			input:         "0000J",
			decimalPlaces: 2,
			expected:      "-0.01",
		},
		{
			name:     "Negative number with K",
			input:    "0001K",
			expected: "-12",
		},
		{
			name:          "Negative number with K and 2 decimal places",
			input:         "0001K",
			decimalPlaces: 2,
			expected:      "-0.12",
		},
		{
			name:     "Negative number with ü",
			input:    "0020ü",
			expected: "-200",
		},
		{
			name:          "Negative number with ü and 2 decimal places",
			input:         "0020ü",
			decimalPlaces: 2,
			expected:      "-2.00",
		},
		{
			name:     "Leading zeros with negative",
			input:    "00001J",
			expected: "-11",
		},
		{
			name:     "Only EBCDIC character",
			input:    "J",
			expected: "-1",
		},
		{
			name:     "Zero with negative",
			input:    "0000ü",
			expected: "-0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result string
			var err error

			result, err = ConvertEBCDICToAsciiNumber(tt.input, tt.decimalPlaces)

			// Check result if no error
			if err == nil && result != tt.expected {
				t.Errorf("EbcdicToAsciiNumber() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestAsciiToEbcdicNumber(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		decimalPlaces int
		expected      string
	}{
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Normal number",
			input:    "12345",
			expected: "12345",
		},
		{
			name:          "Normal number with 2 decimal places",
			input:         "123.45",
			decimalPlaces: 2,
			expected:      "12345",
		},
		{
			name:          "Normal number with 2 decimal places but no decimal point",
			input:         "123",
			decimalPlaces: 2,
			expected:      "12300",
		},
		{
			name:     "Negative number",
			input:    "-1",
			expected: "J",
		},
		{
			name:          "Negative number with 2 decimal places",
			input:         "-0.01",
			decimalPlaces: 2,
			expected:      "J",
		},
		{
			name:     "Negative number with 2 digits",
			input:    "-12",
			expected: "1K",
		},
		{
			name:          "Negative number with decimal places",
			input:         "-0.12",
			decimalPlaces: 2,
			expected:      "1K",
		},
		{
			name:     "Negative number ending with 0",
			input:    "-10",
			expected: "1ü",
		},
		{
			name:          "Negative number with decimal places ending with 0",
			input:         "-2.00",
			decimalPlaces: 2,
			expected:      "20ü",
		},
		{
			name:     "Zero",
			input:    "0",
			expected: "0",
		},
		{
			name:     "Negative zero",
			input:    "-0",
			expected: "ü",
		},
		{
			name:          "Decimal handling with padding",
			input:         "123.4",
			decimalPlaces: 2,
			expected:      "12340",
		},
		{
			name:          "Decimal handling with truncation",
			input:         "123.456",
			decimalPlaces: 2,
			expected:      "12345",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ConvertAsciiToEBCDICNumber(tt.input, tt.decimalPlaces)

			if err != nil {
				t.Errorf("AsciiToEbcdicNumber() error = %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("AsciiToEbcdicNumber() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestRoundTripConversion tests that converting from EBCDIC to ASCII and back results in the original value
func TestRoundTripConversion(t *testing.T) {
	tests := []struct {
		name          string
		ebcdicValue   string
		decimalPlaces int
	}{
		{
			name:        "Normal number",
			ebcdicValue: "12345",
		},
		{
			name:          "Normal number with decimal places",
			ebcdicValue:   "12345",
			decimalPlaces: 2,
		},
		{
			name:        "Negative number with J",
			ebcdicValue: "123J",
		},
		{
			name:          "Negative number with J and decimal places",
			ebcdicValue:   "123J",
			decimalPlaces: 2,
		},
		{
			name:        "Negative number with ü",
			ebcdicValue: "12ü",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert EBCDIC to ASCII
			asciiValue, err := ConvertEBCDICToAsciiNumber(tt.ebcdicValue, tt.decimalPlaces)
			if err != nil {
				t.Errorf("ConvertEBCDICToAsciiNumber() error = %v", err)
				return
			}

			// Convert ASCII back to EBCDIC
			ebcdicValue, err := ConvertAsciiToEBCDICNumber(asciiValue, tt.decimalPlaces)
			if err != nil {
				t.Errorf("ConvertAsciiToEBCDICNumber() error = %v", err)
				return
			}

			// Check if we got the original EBCDIC value back
			if ebcdicValue != tt.ebcdicValue {
				t.Errorf("Round trip conversion failed: original = %v, got = %v", tt.ebcdicValue, ebcdicValue)
			}
		})
	}
}
