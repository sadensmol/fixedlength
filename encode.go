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

	//todo: check for overlapping ranges

	sb := strings.Builder{}

	// use runes to handle utf-8
	lastPos := 0
	for _, tagWitPos := range tagsWithPos {
		tagLen := tagWitPos.tag.toPos - tagWitPos.tag.fromPos
		field := structVal.Field(tagWitPos.fieldNum)
		gap := tagWitPos.tag.fromPos - lastPos

		m, err := MarshalField(field, tagWitPos.tag)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field %s : %w", structVal.Type().Field(tagWitPos.fieldNum).Name, err)
		}

		str := string(m)
		strLen := utf8.RuneCountInString(str)
		// check if field is too long
		if strLen > tagLen {
			return nil, fmt.Errorf("field %s is too long, required: %d but %d", structVal.Type().Field(tagWitPos.fieldNum).Name, tagLen, strLen)
		}

		sb.WriteString(formatStringWithPadding(str, tagLen+gap, stringPaddingTypeLeft))
		lastPos = tagWitPos.tag.toPos
	}
	return []byte(sb.String()), nil
}

func MarshalField(field reflect.Value, t tag) ([]byte, error) {
	var str string
	switch field.Kind() {
	case reflect.Int:
		str = fmt.Sprintf("%d", field.Int())
	case reflect.Float64:
		str = fmt.Sprintf("%.*g", t.Len(), field.Float())
	case reflect.String:
		str = field.String()
	case reflect.Struct:
		if implementsMarshaler(field) {
			m := field.Interface().(Marshaler)
			b, err := m.Marshal()
			if err != nil {
				return nil, err
			}
			return b, nil
		}
	}

	return []byte(str), nil
}
