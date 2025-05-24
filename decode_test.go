package fixedlength

import (
	"errors"
	"reflect"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	type nested struct {
		A string `range:"0,1"`
		B int    `range:"2,3"`
		C float64
	}

	type testStruct struct {
		Nested nested `range:"0,4"`
		Empty  string
	}

	data := []byte("AB2D")
	var v testStruct
	err := Unmarshal(data, &v)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if v.Nested.A != "A" {
		t.Errorf("Expected v.Nested.A to be 'A', got '%s'", v.Nested.A)
	}

	if v.Nested.B != 2 {
		t.Errorf("Expected v.Nested.B to be 2, got %d", v.Nested.B)
	}

	if v.Nested.C != 0 {
		t.Errorf("Expected v.Nested.C to be 0, got %f", v.Nested.C)
	}

	if v.Empty != "" {
		t.Errorf("Expected v.Empty to be empty, got '%s'", v.Empty)
	}
}

func TestUnmarshalError(t *testing.T) {
	data := []byte("AB2D")

	t.Run("non-pointer int", func(t *testing.T) {
		var v int
		err := Unmarshal(data, v)
		if err == nil {
			t.Fatalf("Expected Unmarshal to fail")
		}

		if err.Error() != "range: Unmarshal(non-pointer int)" {
			t.Errorf("Expected error to be 'range: Unmarshal(non-pointer int)', got '%s'", err.Error())
		}
	})

	t.Run("nil *int", func(t *testing.T) {
		var v2 *int
		err := Unmarshal(data, v2)
		if err == nil {
			t.Fatalf("Expected Unmarshal to fail")
		}

		if err.Error() != "range: Unmarshal(nil *int)" {
			t.Errorf("Expected error to be 'range: Unmarshal(nil *int)', got '%s'", err.Error())
		}
	})

	t.Run("non-pointer struct {}", func(t *testing.T) {
		var v3 struct{}
		err := Unmarshal(data, v3)
		if err == nil {
			t.Fatalf("Expected Unmarshal to fail")
		}

		if err.Error() != "range: Unmarshal(non-pointer struct {})" {
			t.Errorf("Expected error to be 'range: Unmarshal(non-pointer struct {})', got '%s'", err.Error())
		}
	})

	t.Run("nil", func(t *testing.T) {
		err := Unmarshal(data, nil)
		if err == nil {
			t.Fatalf("Expected Unmarshal to fail")
		}

		if err.Error() != "range: Unmarshal(nil)" {
			t.Errorf("Expected error to be 'range: Unmarshal(non-pointer struct {})', got '%s'", err.Error())
		}
	})
}

// Define a type that implements the Unmarshaler interface
type CustomTime struct{}

var _ Unmarshaler = (*CustomTime)(nil)

func (ct *CustomTime) Unmarshal(data []byte) error {
	return nil
}

// Define a type that does NOT implement the Unmarshaler interface
type NotUnmarshaler struct{}

// Test suite for implementsUnmarshaler function
func TestImplementsUnmarshaler(t *testing.T) {
	// Case 1: Type that implements Unmarshaler
	t.Run("implements Unmarshaler", func(t *testing.T) {
		// Create a reflect.Value for a type that implements Unmarshaler
		ct := &CustomTime{}
		val := reflect.ValueOf(ct)

		if !implementsUnmarshaler(val) {
			t.Errorf("expected true, got false for CustomTime")
		}
	})

	// Case 2: Type that does NOT implement Unmarshaler
	t.Run("does not implement Unmarshaler", func(t *testing.T) {
		// Create a reflect.Value for a type that does not implement Unmarshaler
		notUnmarshaler := NotUnmarshaler{}
		val := reflect.ValueOf(notUnmarshaler)

		if implementsUnmarshaler(val) {
			t.Errorf("expected false, got true for NotUnmarshaler")
		}
	})

	// Case 3: Non-addressable value
	t.Run("non-addressable value", func(t *testing.T) {
		// Use a literal value, which is non-addressable
		val := reflect.ValueOf(42) // An int, which is non-addressable in this context

		if implementsUnmarshaler(val) {
			t.Errorf("expected false, got true for non-addressable int")
		}
	})

	// Case 4: Nil pointer that implements Unmarshaler
	t.Run("nil pointer implements Unmarshaler", func(t *testing.T) {
		// Use a nil pointer to CustomTime (which implements Unmarshaler)
		var ct *CustomTime
		val := reflect.ValueOf(ct)

		if !implementsUnmarshaler(val) {
			t.Errorf("expected true, got false for nil pointer to CustomTime")
		}
	})

	// Case 5: Nil pointer that does NOT implement Unmarshaler
	t.Run("nil pointer does not implement Unmarshaler", func(t *testing.T) {
		// Use a nil pointer to NotUnmarshaler (which does not implement Unmarshaler)
		var notUnmarshaler *NotUnmarshaler
		val := reflect.ValueOf(notUnmarshaler)

		if implementsUnmarshaler(val) {
			t.Errorf("expected false, got true for nil pointer to NotUnmarshaler")
		}
	})
}

func TestSetFieldValue(t *testing.T) {
	tests := []struct {
		name      string
		kind      reflect.Kind
		value     string
		wantValue interface{}
		wantErr   error
	}{
		// Int tests
		{
			name:      "valid int",
			kind:      reflect.Int,
			value:     "42",
			wantValue: int64(42),
			wantErr:   nil,
		},
		{
			name:      "invalid int",
			kind:      reflect.Int,
			value:     "invalid",
			wantValue: int64(0),
			wantErr:   ErrInvalidIntValue,
		},

		// Float tests
		{
			name:      "valid float",
			kind:      reflect.Float64,
			value:     "42.5",
			wantValue: float64(42.5),
			wantErr:   nil,
		},
		{
			name:      "invalid float",
			kind:      reflect.Float64,
			value:     "invalid",
			wantValue: float64(0),
			wantErr:   ErrInvalidFloatValue,
		},

		// String tests
		{
			name:      "valid string",
			kind:      reflect.String,
			value:     "hello",
			wantValue: "hello",
			wantErr:   nil,
		},

		// Bool tests
		{
			name:      "valid bool true",
			kind:      reflect.Bool,
			value:     "true",
			wantValue: true,
			wantErr:   nil,
		},
		{
			name:      "valid bool false",
			kind:      reflect.Bool,
			value:     "false",
			wantValue: false,
			wantErr:   nil,
		},
		{
			name:      "invalid bool",
			kind:      reflect.Bool,
			value:     "invalid",
			wantValue: false,
			wantErr:   ErrInvalidBooleanValue,
		},

		// Unsupported kind
		{
			name:      "unsupported kind",
			kind:      reflect.Slice,
			value:     "somevalue",
			wantValue: nil,
			wantErr:   ErrUnsupportedKind,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a zero value for the given type
			var field reflect.Value
			switch tt.kind {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				field = reflect.New(reflect.TypeOf(int64(0))).Elem()
			case reflect.Float32, reflect.Float64:
				field = reflect.New(reflect.TypeOf(float64(0))).Elem()
			case reflect.String:
				field = reflect.New(reflect.TypeOf("")).Elem()
			case reflect.Bool:
				field = reflect.New(reflect.TypeOf(false)).Elem()
			default:
				field = reflect.New(reflect.TypeOf([]byte{})).Elem() // Unsupported kind
			}

			// Call setFieldValue and check for errors
			// Create an empty tag for testing purposes
			emptyTag := tag{}
			err := setFieldValue(field, tt.value, emptyTag)

			// Check for expected error
			if err != nil && tt.wantErr == nil {
				t.Errorf("expected no error, got %v", err)
			}
			if err == nil && tt.wantErr != nil {
				t.Errorf("expected error %v, got none", tt.wantErr)
			}
			if err != nil && tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("expected error %v, got %v", tt.wantErr, err)
			}

			// Check the value set in the field (only for supported kinds)
			if err == nil {
				switch tt.kind {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					if field.Int() != tt.wantValue {
						t.Errorf("expected int value %v, got %v", tt.wantValue, field.Int())
					}
				case reflect.Float32, reflect.Float64:
					if field.Float() != tt.wantValue {
						t.Errorf("expected float value %v, got %v", tt.wantValue, field.Float())
					}
				case reflect.String:
					if field.String() != tt.wantValue {
						t.Errorf("expected string value %v, got %v", tt.wantValue, field.String())
					}
				case reflect.Bool:
					if field.Bool() != tt.wantValue {
						t.Errorf("expected bool value %v, got %v", tt.wantValue, field.Bool())
					}
				}
			}
		})
	}
}

func TestUnmarshalWithDecimals(t *testing.T) {
	type testStruct struct {
		Amount float64 `range:"0,6" decimals:"2"`
	}

	tests := []struct {
		name        string
		data        []byte
		expected    float64
		expectedErr error
	}{
		{
			name:        "positive float with decimals",
			data:        []byte("012345"),
			expected:    123.45,
			expectedErr: nil,
		},
		{
			name:        "zero padded with decimals",
			data:        []byte("000199"),
			expected:    1.99,
			expectedErr: nil,
		},
		{
			name:        "negative number with decimals",
			data:        []byte("01234K"), // K represents negative 2 in EBCDIC
			expected:    -123.42,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var v testStruct
			err := Unmarshal(tt.data, &v)

			if (err != nil) != (tt.expectedErr != nil) {
				t.Errorf("Unmarshal() error = %v, expected error %v", err, tt.expectedErr)
				return
			}

			if err != nil && tt.expectedErr != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("Unmarshal() error = %v, expected error %v", err, tt.expectedErr)
				return
			}

			if v.Amount != tt.expected {
				t.Errorf("Expected Amount to be %v, got %v", tt.expected, v.Amount)
			}
		})
	}
}
