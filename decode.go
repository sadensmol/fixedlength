package fixedlength

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var (
	ErrInvalidBooleanValue = errors.New("fixedlength: invalid boolean value")
	ErrInvalidIntValue     = errors.New("fixedlength: invalid int value")
	ErrInvalidFloatValue   = errors.New("fixedlength: invalid float value")
	ErrUnsupportedKind     = errors.New("fixedlength: unsupported kind")
)

// setFieldValue sets the value for a struct field using reflection.
func setFieldValue(field reflect.Value, value string, tag tag) error {

	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

		cValue, err := ConvertEBCDICToAsciiNumber(value, tag.decimals)
		if err != nil {
			return err
		}

		intVal, err := strconv.ParseInt(cValue, 10, 64)
		if err != nil {
			return errors.Join(ErrInvalidIntValue, err)
		}
		field.SetInt(intVal)

	case reflect.Float32, reflect.Float64:
		cValue, err := ConvertEBCDICToAsciiNumber(value, tag.decimals)
		if err != nil {
			return err
		}

		floatVal, err := strconv.ParseFloat(cValue, 64)
		if err != nil {
			return errors.Join(ErrInvalidFloatValue, err)
		}
		field.SetFloat(floatVal)

	case reflect.String:
		field.SetString(value)

	case reflect.Bool:
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return errors.Join(ErrInvalidBooleanValue, err)
		}
		field.SetBool(boolVal)

	default:
		if implementsUnmarshaler(field) {
			um := field.Addr().Interface().(Unmarshaler)
			return um.Unmarshal([]byte(value))
		}

		return fmt.Errorf("%w: %s", ErrUnsupportedKind, field.Kind())
	}
	return nil
}

// Unmarshaler is the interface implemented by types
// that can unmarshal themselves.
// Unmarshal must copy the input data if it wishes
// to retain the data after returning.
type Unmarshaler interface {
	Unmarshal([]byte) error
}

var unmarshalerType = reflect.TypeOf((*Unmarshaler)(nil)).Elem()

// implementsUnmarshaler checks if a field implements the Unmarshaler interface
func implementsUnmarshaler(val reflect.Value) bool {
	// If the value is invalid (e.g., a nil value), return false
	if !val.IsValid() {
		return false
	}

	// Check if the value itself implements json.Unmarshaler
	if val.Type().Implements(unmarshalerType) {
		return true
	}

	// If the value is addressable, check if the pointer to it implements json.Unmarshaler
	if val.CanAddr() {
		return val.Addr().Type().Implements(unmarshalerType)
	}

	// Otherwise, it does not implement the interface
	return false
}

// InvalidUnmarshalError describes an invalid argument passed to [Unmarshal].
// (The argument to [Unmarshal] must be a non-nil pointer.)
type InvalidUnmarshalError struct {
	Type reflect.Type
}

func (e InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "range: Unmarshal(nil)"
	}
	if e.Type.Kind() != reflect.Pointer {
		return "range: Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "range: Unmarshal(nil " + e.Type.String() + ")"
}

// Unmarshal parses the given string into the provided struct v.
// v must be a pointer to a struct, and its fields should be tagged with `range:"<start>,<end>"`
// where start and end are the lower and upper bounds of the segment in the string.
// Unmarshal will parse nested structs recursively.
func Unmarshal(data []byte, v any) error {
	// Validate that v is a pointer to a struct
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return InvalidUnmarshalError{reflect.TypeOf(v)}
	}

	// convert to runes since we use utf-8 here
	runes := []rune(string(data))
	// Iterate over struct fields to map segment names to fields
	for i := 0; i < rv.Elem().NumField(); i++ {
		field := rv.Elem().Field(i)

		// Recursively parse the struct
		if field.Kind() == reflect.Struct && !implementsUnmarshaler(field) {
			if err := Unmarshal(data, field.Addr().Interface()); err != nil {
				return err
			}

			continue
		}

		tag, err := parseFieldTag(rv.Elem().Type().Field(i).Tag)
		if err != nil {
			if errors.Is(err, ErrTagEmpty) {
				continue
			}

			if tag.flags.optional {
				continue
			}
			return fmt.Errorf("failed to parse tag %s (%s) : %w", rv.Elem().Type().Field(i).Name, tag, err)
		}

		l := len(runes)
		err = tag.Validate(l)
		if err != nil {
			if tag.flags.optional {
				continue
			}
			return fmt.Errorf("failed to validate tag %s (%s) : %w", rv.Elem().Type().Field(i).Name, tag, err)
		}

		value := strings.TrimSpace(string(runes[tag.fromPos:tag.toPos]))

		if err := setFieldValue(field, value, tag); err != nil {
			if tag.flags.optional {
				continue
			}
			return fmt.Errorf("failed to set field value %s (%s) : %w", rv.Elem().Type().Field(i).Name, tag, err)
		}
	}

	return nil
}
