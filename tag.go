package fixedlength

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var (
	ErrTagInefectualRange    = errors.New("fixedlength: inefectual range")
	ErrTagEmpty              = errors.New("fixedlength: tag is empty")
	ErrTagInvalidRangeValues = errors.New("fixedlength: invalid range values")
	ErrTagInvalidUpperBound  = errors.New("fixedlength: invalid upper bound")
)

type tag struct {
	fromPos int
	toPos   int
	flags   flags
}

type flags struct {
	optional bool
}

func (t tag) String() string {
	return fmt.Sprintf("range:%d,%d flags:%v", t.fromPos, t.toPos, t.flags)
}

func (f flags) String() string {
	return fmt.Sprintf("optional:%t", f.optional)
}

func parseFieldTag(t reflect.StructTag, upperBound int) (tag, error) {
	res := tag{}

	flagsTag := t.Get("flags")
	flags, err := parseFlagsTag(flagsTag)
	if err != nil {
		return res, err
	}

	res.flags = flags

	rangeTag := t.Get("range")
	start, end, err := parseRangeTag(rangeTag, upperBound)
	if err != nil {
		return res, err
	}
	res.fromPos = start
	res.toPos = end

	return res, nil
}

func parseFlagsTag(tag string) (flags, error) {
	f := flags{}
	if tag == "" {
		return f, nil
	}

	parts := strings.Split(tag, ",")
	for _, part := range parts {
		switch part {
		case "optional":
			f.optional = true
		}
	}

	return f, nil
}

// parseRangeTag splits a struct field's json tag into its name and
// comma-separated options.
func parseRangeTag(tag string, upperBound int) (int, int, error) {
	if tag == "" {
		return 0, 0, ErrTagEmpty
	}

	if upperBound == 0 {
		return 0, 0, fmt.Errorf("%w: %d", ErrTagInvalidUpperBound, upperBound)
	}

	parts := strings.Split(tag, ",")
	x, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, errors.Join(ErrTagInvalidRangeValues, err)
	}

	y, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, errors.Join(ErrTagInvalidRangeValues, err)
	}

	start := max(x, 0)
	end := min(y, upperBound)

	// -1 is used to indicate that the end of the string should be used
	if end == -1 {
		end = upperBound
	}

	if start == end {
		return 0, 0, fmt.Errorf("%w: %s", ErrTagInefectualRange, tag)
	}

	if start > end {
		return 0, 0, fmt.Errorf("%w: x > y (%d) from %s", ErrTagInefectualRange, upperBound, tag)
	}

	return start, end, nil
}
