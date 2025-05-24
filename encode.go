package fixedlength

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"unicode/utf8"
)

type Marshaler interface {
	Marshal() ([]byte, error)
}

var marshalerType = reflect.TypeOf((*Marshaler)(nil)).Elem()

func implementsMarshaler(val reflect.Value) bool {
	if !val.IsValid() {
		return false
	}

	if val.Type().Implements(marshalerType) {
		return true
	}

	// If the value is addressable, check if the pointer to it implements json.Unmarshaler
	if val.CanAddr() {
		return val.Addr().Type().Implements(marshalerType)
	}

	// Otherwise, it does not implement the interface
	return false
}

// working with gaps:
// field 1 [200,210)
// field 2 [220, 230)
// the gap between them : 10, starting from 210 (included ) and ended 219 (included)
// 209|..........|220
func Marshal(d interface{}) ([]byte, error) {
	rv := reflect.ValueOf(d)
	var structVal reflect.Value

	switch rv.Kind() {
	case reflect.Ptr:
		if rv.IsNil() {
			return nil, fmt.Errorf("cannot marshal nil pointer")
		}
		structVal = rv.Elem()
	case reflect.Struct:
		structVal = rv
	default:
		return nil, fmt.Errorf("invalid marshal value")
	}

	type tagWithFieldNumber struct {
		tag      tag
		fieldNum int
	}
	tagsWithPos := make([]tagWithFieldNumber, 0)

	for i := 0; i < structVal.NumField(); i++ {
		tag, err := parseFieldTag(structVal.Type().Field(i).Tag)
		if err != nil {
			if errors.Is(err, ErrTagEmpty) {
				continue
			}
			if tag.flags.optional {
				continue
			}
			return nil, fmt.Errorf("failed to parse tag %s (%s) : %w", structVal.Type().Field(i).Name, tag, err)
		}

		tagsWithPos = append(tagsWithPos, tagWithFieldNumber{tag, i})
	}

	sort.Slice(tagsWithPos, func(i, j int) bool {
		return tagsWithPos[i].tag.fromPos < tagsWithPos[j].tag.fromPos
	})

	sb := strings.Builder{}
	// use runes to handle utf-8
	lastPos := 2 // we always start at 2 since the first two characters are the type and always filled outside
	for _, tagWitPos := range tagsWithPos {

		field := structVal.Field(tagWitPos.fieldNum)
		str, err := MarshalField(field, tagWitPos.tag)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field %s : %w", structVal.Type().Field(tagWitPos.fieldNum).Name, err)
		}

		strStr := string(str)
		strLen := utf8.RuneCountInString(strStr)
		// check if field is too long
		tagLen := tagWitPos.tag.toPos - tagWitPos.tag.fromPos
		if strLen > tagLen {
			return nil, fmt.Errorf("field %s is too long, required: %d but %d", structVal.Type().Field(tagWitPos.fieldNum).Name, tagLen, strLen)
		}

		gap := tagWitPos.tag.fromPos - lastPos
		if gap < 0 {
			return nil, fmt.Errorf("field %s is overlapping with previous field", structVal.Type().Field(tagWitPos.fieldNum).Name)
		}

		if gap > 0 {
			sb.WriteString(fmt.Sprintf("%*s", gap, ""))
		}

		// write the original string
		sb.WriteString(strStr)

		lastPos = tagWitPos.tag.toPos
	}
	return []byte(sb.String()), nil
}

func MarshalField(field reflect.Value, t tag) ([]byte, error) {
	var str string
	var err error

	align := GetConfig().AlignmentType
	if t.align != AlignmentTypeNone {
		align = t.align
	}

	switch field.Kind() {
	case reflect.String:
		str = field.String()

		str, err = FormatStringWithAlignment(str, t.Len(), align)
		if err != nil {
			return nil, err
		}

	case reflect.Int:
		val := field.Int()

		cVal, err := ConvertAsciiToEBCDICNumber(fmt.Sprintf("%d", val), tag.decimals)
		if err != nil {
			return nil, fmt.Errorf("failed to convert int to EBCDIC: %w", err)
		}

		str, err = FormatStrNumberWithAlignment(cVal, t.Len(), GetConfig().NumbersWithLeadingZeroes, align)
		if err != nil {
			return nil, err
		}
	case reflect.Float64:
		val := field.Float()
		cVal, err := ConvertAsciiToEBCDICNumber(fmt.Sprintf("%f", val), tag.decimals)
		str, err = FormatStrNumberWithAlignment(cVal, t.Len(), GetConfig().NumbersWithLeadingZeroes, align)
		if err != nil {
			return nil, err
		}
	case reflect.Struct:
		if implementsMarshaler(field) {
			m := field.Interface().(Marshaler)
			ba, err := m.Marshal()
			if err != nil {
				return nil, err
			}
			str, err = FormatStringWithAlignment(string(ba), t.Len(), align)
			if err != nil {
				return nil, err
			}
		}
	}

	return []byte(str), nil
}
