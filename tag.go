package fixedlength

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var (
	ErrTagInefectualRange    = errors.New("inefectual range")
	ErrTagEmpty              = errors.New("tag is empty")
	ErrTagInvalidRangeValues = errors.New("invalid range values")
	ErrTagInvalidUpperBound  = errors.New("invalid upper bound")
)

type tag struct {
	fromPos  int
	toPos    int
	flags    flags
	align    AlignmentType
	decimals int
}

func (t tag) Len() int {
	return t.toPos - t.fromPos
}

func (t tag) Validate(maxPos int) error {
	if t.toPos > maxPos {
		return fmt.Errorf("to pos is higher that length (%d): %s", maxPos, t)
	}

	if t.fromPos < 0 || t.toPos <= t.fromPos {
		return fmt.Errorf("invalid range values: %s", t)
	}

	return nil
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

func parseFieldTag(t reflect.StructTag) (tag, error) {
	res := tag{}

	flagsTag := t.Get("flags")
	flags, err := parseFlagsTag(flagsTag)
	if err != nil {
		return res, err
	}
	res.flags = flags

	alignTag := t.Get("align")
	align, err := parseAlignTag(alignTag)
	if err != nil {
		return res, err
	}
	res.align = align

	decimalsTag := t.Get("decimals")
	decimals, err := parseDecimalsTag(decimalsTag)
	if err != nil {
		return res, err
	}
	res.decimals = decimals

	rangeTag := t.Get("range")
	start, end, err := parseRangeTag(rangeTag)
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

func parseDecimalsTag(tag string) (int, error) {
	if tag == "" {
		return -1, nil // -1 means no decimals
	}

	decimals, err := strconv.ParseInt(tag, 10, 0)
	if err != nil {
		return 0, fmt.Errorf("invalid decimals tag: %s", tag)
	}

	return int(decimals), nil
}

func parseAlignTag(tag string) (AlignmentType, error) {
	if tag == "" {
		return AlignmentTypeNone, nil
	}

	switch tag {
	case "left":
		return AlignmentTypeLeft, nil
	case "right":
		return AlignmentTypeRight, nil
	}

	return AlignmentTypeNone, fmt.Errorf("invalid align type: %s", tag)
}

// parseRangeTag splits a struct field's json tag into its name and
// comma-separated options.
func parseRangeTag(tag string) (int, int, error) {
	if tag == "" {
		return 0, 0, ErrTagEmpty
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

	return x, y, nil
}
