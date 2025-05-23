package fixedlength

import (
	"testing"
)

func TestEbcdicToAsciiNumber(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "Empty string",
			input:    "",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Normal number",
			input:    "12345",
			expected: "12345",
			wantErr:  false,
		},
		{
			name:     "Negative number with J",
			input:    "0000J",
			expected: "-1",
			wantErr:  false,
		},
		{
			name:     "Negative number with K",
			input:    "0001K",
			expected: "-12",
			wantErr:  false,
		},
		{
			name:     "Negative number with ü",
			input:    "0020ü",
			expected: "-200",
			wantErr:  false,
		},
		{
			name:     "Leading zeros with negative",
			input:    "00001J",
			expected: "-11",
			wantErr:  false,
		},
		{
			name:     "Only EBCDIC character (still treated as negative)",
			input:    "J",
			expected: "-1",
			wantErr:  false,
		},
		{
			name:     "Zero with negative",
			input:    "0000ü",
			expected: "-0",
			wantErr:  false,
		},
		{
			name:     "Float value",
			input:    "123.45",
			expected: "123.45",
			wantErr:  false,
		},
		{
			name:     "Negative float with J",
			input:    "123.4J",
			expected: "-123.41",
			wantErr:  false,
		},
		{
			name:     "Negative float with K",
			input:    "0.01K",
			expected: "-0.012",
			wantErr:  false,
		},
		{
			name:     "Float with leading zeros and negative",
			input:    "00123.40ü",
			expected: "-123.400",
			wantErr:  false,
		},
		{
			name:     "Float with trailing zeros and negative",
			input:    "123.00J",
			expected: "-123.001",
			wantErr:  false,
		},
		{
			name:     "Zero float with negative",
			input:    "0.00J",
			expected: "-0.001",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := EbcdicToAsciiNumber(tt.input)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("EbcdicToAsciiNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check result if no error
			if err == nil && result != tt.expected {
				t.Errorf("EbcdicToAsciiNumber() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
