package fixedlength

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type intS struct {
	string
}

var _ Marshaler = &intS{}

func (i intS) Marshal() ([]byte, error) {
	switch i.string {
	case "yes":
		return []byte("42"), nil
	}
	return nil, errors.New("invalid value")
}

func TestMarshal(t *testing.T) {
	t.Run("simple struct", func(t *testing.T) {
		type s1 struct {
			Field1 string  `range:"0,5" flags:"optional"`
			Field2 int     `range:"5,10"`
			Field3 float64 `range:"10,20"`
		}

		res, err := Marshal(&s1{Field1: "hello", Field2: 42, Field3: 3.14})
		require.NoError(t, err)

		require.Equal(t, "hello000420000003.14", string(res))
	})
	t.Run("struct with internal struct", func(t *testing.T) {

		type s1 struct {
			Field1    string `range:"0,5" flags:"optional"`
			IntStruct intS   `range:"5,10"`
		}

		res, err := Marshal(&s1{Field1: "hello", IntStruct: intS{"yes"}})
		require.NoError(t, err)

		require.Equal(t, "hello42   ", string(res))
	})
}
func TestMarshalIntegers(t *testing.T) {
	type intStruct struct {
		Value int `range:"0,10"`
	}

	type smallIntStruct struct {
		Value int `range:"0,5"`
	}

	type intStructWithRightTag struct {
		Value int `range:"0,10" align:"right"`
	}

	type intStructWithLeftTag struct {
		Value int `range:"0,10" align:"left"`
	}

	tests := []struct {
		name      string
		value     int
		fieldTag  string
		expected  string
		expectErr bool
	}{
		{
			name:     "positive integer",
			value:    42,
			expected: "0000000042",
		},
		{
			name:     "zero integer",
			value:    0,
			expected: "0000000000",
		},
		{
			name:     "negative integer",
			value:    -42,
			expected: "00000004K", // K represents negative 2 in EBCDIC
		},
		{
			name:     "right aligned integer",
			value:    42,
			fieldTag: "right",
			expected: "0000000042", // Fill with zeros on the left side
		},
		{
			name:     "negative right aligned integer",
			value:    -42,
			fieldTag: "right",
			expected: "00000004K", // K represents negative 2 in EBCDIC, zeros on left
		},
		{
			name:     "left aligned integer",
			value:    42,
			fieldTag: "left",
			expected: "042       ",
		},
		{
			name:     "negative left aligned integer",
			value:    -42,
			fieldTag: "left",
			expected: "04K       ", // K represents negative 2 in EBCDIC
		},
		{
			name:     "large integer",
			value:    1234567890,
			expected: "1234567890",
		},
		{
			name:     "larger negative integer",
			value:    -123456789,
			expected: "12345678R", // R represents negative 9 in EBCDIC
		},
		{
			name:      "integer overflow",
			value:     1234567890,
			fieldTag:  "small",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			var res []byte

			if tt.fieldTag == "right" {
				res, err = Marshal(&intStructWithRightTag{Value: tt.value})
			} else if tt.fieldTag == "left" {
				res, err = Marshal(&intStructWithLeftTag{Value: tt.value})
			} else if tt.fieldTag == "small" {
				res, err = Marshal(&smallIntStruct{Value: tt.value})
			} else {
				res, err = Marshal(&intStruct{Value: tt.value})
			}

			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, string(res))
			}
		})
	}
}

func TestMarshalFloats(t *testing.T) {
	type floatStruct struct {
		Value float64 `range:"0,10"`
	}

	type smallFloatStruct struct {
		Value float64 `range:"0,5"`
	}

	type floatStructWithDecimals struct {
		Value float64 `range:"0,10" decimals:"3"`
	}

	type floatStructWithRightAlign struct {
		Value float64 `range:"0,10" align:"right"`
	}

	type floatStructWithLeftAlign struct {
		Value float64 `range:"0,10" align:"left"`
	}

	tests := []struct {
		name      string
		value     float64
		fieldTag  string
		expected  string
		expectErr bool
	}{
		{
			name:     "positive float",
			value:    3.14,
			expected: "0000003.14",
		},
		{
			name:     "float with decimals tag",
			value:    3.14159,
			fieldTag: "decimals",
			expected: "0000003.142",
		},
		{
			name:     "zero float",
			value:    0.0,
			expected: "0000000.00",
		},
		{
			name:     "negative float",
			value:    -42.5,
			expected: "0004250ü", // ü represents negative 0 in EBCDIC
		},
		{
			name:     "right aligned float",
			value:    3.14,
			fieldTag: "right",
			expected: "0000003.14", // Fill with zeros on the left side
		},
		{
			name:     "negative right aligned float",
			value:    -42.5,
			fieldTag: "right",
			expected: "0004250ü", // ü represents negative 0 in EBCDIC, zeros on left
		},
		{
			name:     "left aligned float",
			value:    3.14,
			fieldTag: "left",
			expected: "03.14     ",
		},
		{
			name:     "negative left aligned float",
			value:    -42.5,
			fieldTag: "left",
			expected: "0425ü     ", // ü represents negative 0 in EBCDIC
		},
		{
			name:     "larger negative float",
			value:    -9876.54,
			expected: "000987654Q", // Q represents negative 8 in EBCDIC
		},
		{
			name:      "float overflow",
			value:     1234.56,
			fieldTag:  "small",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			var res []byte

			if tt.fieldTag == "decimals" {
				res, err = Marshal(&floatStructWithDecimals{Value: tt.value})
			} else if tt.fieldTag == "right" {
				res, err = Marshal(&floatStructWithRightAlign{Value: tt.value})
			} else if tt.fieldTag == "left" {
				res, err = Marshal(&floatStructWithLeftAlign{Value: tt.value})
			} else if tt.fieldTag == "small" {
				res, err = Marshal(&smallFloatStruct{Value: tt.value})
			} else {
				res, err = Marshal(&floatStruct{Value: tt.value})
			}

			if tt.expectErr {
				require.Error(t, err)
				if tt.fieldTag == "small" {
					require.Contains(t, err.Error(), "too long")
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, string(res))
			}
		})
	}
}
