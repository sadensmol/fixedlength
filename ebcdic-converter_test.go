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
