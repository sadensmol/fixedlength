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
		Value int `range:"2,12"`
	}

	type smallIntStruct struct {
		Value int `range:"2,7"`
	}

	type intStructWithRightTag struct {
		Value int `range:"2,12" align:"right"`
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
			expected: "000000004K", // K represents negative 2 in EBCDIC
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
			expected: "000000004K", // K represents negative 2 in EBCDIC, zeros on left
		},

		{
			name:     "large integer",
			value:    1234567890,
			expected: "1234567890",
		},
		{
			name:     "larger negative integer",
			value:    -123456789,
			expected: "012345678R", // R represents negative 9 in EBCDIC
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
		Value float64 `range:"2,12" decimals:"2"`
	}

	type floatStructWith3Decimals struct {
		Value float64 `range:"2,12" decimals:"3"`
	}

	type smallFloatStruct struct {
		Value float64 `range:"2,7" decimals:"2"`
	}

	type floatStructWithRightAlign struct {
		Value float64 `range:"2,12" decimals:"2" align:"right"`
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
			expected: "0000000314",
		},
		{
			name:     "float with decimals tag",
			value:    3.14159,
			fieldTag: "decimals3",
			expected: "0000003141",
		},
		{
			name:     "zero float",
			value:    0.0,
			expected: "0000000000",
		},
		{
			name:     "negative float",
			value:    -42.5,
			expected: "000000425ü", // ü represents negative 0 in EBCDIC
		},
		{
			name:     "right aligned float",
			value:    3.14,
			fieldTag: "right",
			expected: "0000000314", // Fill with zeros on the left side
		},
		{
			name:     "negative right aligned float",
			value:    -42.5,
			fieldTag: "right",
			expected: "000000425ü", // ü represents negative 0 in EBCDIC, zeros on left
		},

		{
			name:     "larger negative float",
			value:    -9876.58,
			expected: "000098765Q", // Q represents negative 8 in EBCDIC
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

			if tt.fieldTag == "decimals3" {
				res, err = Marshal(&floatStructWith3Decimals{Value: tt.value})
			} else if tt.fieldTag == "right" {
				res, err = Marshal(&floatStructWithRightAlign{Value: tt.value})
			} else if tt.fieldTag == "small" {
				res, err = Marshal(&smallFloatStruct{Value: tt.value})
			} else {
				res, err = Marshal(&floatStruct{Value: tt.value})
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
