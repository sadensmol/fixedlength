package fixedlength

import (
	"fmt"
)

type AlignmentType int

var (
	AlignmentTypeNone   AlignmentType = 0
	AlignmentTypeRight  AlignmentType = 1
	AlignmentTypeLeft   AlignmentType = 2
	AlignmentTypeCenter AlignmentType = 3
)

func FormatStringWithAlignment(str string, length int, alignmentType AlignmentType) (string, error) {
	return formatWithFiller(str, length, alignmentType, ' ')
}

func FormatStrNumberWithAlignment(strNum string, length int, leadingZeroes bool, alignmentType AlignmentType) (string, error) {
	str := fmt.Sprintf("%s", strNum)
	filler := byte(' ')
	if leadingZeroes {
		filler = '0'
		alignmentType = AlignmentTypeRight // we always add zeroes in the beginning
	}
	return formatWithFiller(str, length, alignmentType, filler)
}

// formatWithFiller handles string formatting with specified filler character
func formatWithFiller(str string, length int, alignmentType AlignmentType, filler byte) (string, error) {
	runes := []rune(str)
	runeCount := len(runes)
	if runeCount > length {
		return "", fmt.Errorf("string %s length %d exceeds target length %d", str, runeCount, length)
	}

	result := make([]rune, length)
	fillerRune := rune(filler)
	diff := length - runeCount

	switch alignmentType {
	case AlignmentTypeLeft:
		// Copy string first
		copy(result, runes)
		// Fill remaining with filler
		for i := runeCount; i < length; i++ {
			result[i] = fillerRune
		}

	case AlignmentTypeRight:
		// Fill padding first
		for i := 0; i < diff; i++ {
			result[i] = fillerRune
		}
		// Copy string after padding
		copy(result[diff:], runes)

	case AlignmentTypeCenter:
		leftPad := diff / 2

		// Fill left padding
		for i := 0; i < leftPad; i++ {
			result[i] = fillerRune
		}
		// Copy string in middle
		copy(result[leftPad:], runes)
		// Fill right padding
		for i := leftPad + runeCount; i < length; i++ {
			result[i] = fillerRune
		}

	default:
		return "", fmt.Errorf("unsupported alignment type")
	}

	return string(result), nil
}
