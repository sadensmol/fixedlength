package fixedlength

import (
	"fmt"
)

type AlignmentType int

var (
	AlignmentTypeRight  AlignmentType = 1
	AlignmentTypeLeft   AlignmentType = 2
	AlignmentTypeCenter AlignmentType = 3
)

func FormatStringWithAlignment(str string, length int, alignmentType AlignmentType) (string, error) {
	return formatWithFiller(str, length, alignmentType, ' ')
}

// FormatFloatWithAlignment formats float64 with specified decimal places
func FormatFloatWithAlignment(num float64, decimals int, leadingZeroes bool, length int, alignmentType AlignmentType) (string, error) {
	str := fmt.Sprintf("%.*f", decimals, num)
	filler := byte(' ')
	if leadingZeroes {
		filler = '0'
		alignmentType = AlignmentTypeRight // we always add zeroes in the beginning
	}
	return formatWithFiller(str, length, alignmentType, filler)
}

// FormatIntWithAlignment formats integer
func FormatIntWithAlignment(num int, length int, leadingZeroes bool, alignmentType AlignmentType) (string, error) {
	str := fmt.Sprintf("%d", num)
	filler := byte(' ')
	if leadingZeroes {
		filler = '0'
		alignmentType = AlignmentTypeRight // we always add zeroes in the beginning
	}
	return formatWithFiller(str, length, alignmentType, filler)
}

// formatWithFiller handles string formatting with specified filler character
func formatWithFiller(str string, length int, alignmentType AlignmentType, filler byte) (string, error) {
	if len(str) > length {
		return "", fmt.Errorf("string length exceeds target length")
	}

	buf := make([]byte, length)
	strLen := len(str)
	diff := length - strLen

	switch alignmentType {
	case AlignmentTypeLeft:
		// Copy string first
		copy(buf[0:], str)
		// Fill remaining with filler
		for i := strLen; i < length; i++ {
			buf[i] = filler
		}

	case AlignmentTypeRight:
		// Fill padding first
		for i := 0; i < diff; i++ {
			buf[i] = filler
		}
		// Copy string after padding
		copy(buf[diff:], str)

	case AlignmentTypeCenter:
		leftPad := diff / 2

		// Fill left padding
		for i := 0; i < leftPad; i++ {
			buf[i] = filler
		}
		// Copy string in middle
		copy(buf[leftPad:], str)
		// Fill right padding
		for i := leftPad + strLen; i < length; i++ {
			buf[i] = filler
		}

	default:
		return "", fmt.Errorf("unsupported alignment type")
	}

	return string(buf), nil
}
