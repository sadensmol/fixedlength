package fixedlength

import (
	"strconv"
	"strings"
)

var ebcdicToASCIINegativeMap = map[rune]int{
	'J': 1,
	'K': 2,
	'L': 3,
	'M': 4,
	'N': 5,
	'O': 6,
	'P': 7,
	'Q': 8,
	'R': 9,
	'ü': 0,
}

// ConvertEBCDICToAsciiNumber converts an EBCDIC number to an ASCII number.
// input EBCDIC value doesn't contain any decimal point.
// It handles negative numbers represented by specific characters.
// The function also allows specifying the number of decimal places, this should apply after the negative conversion if decimal
// places number > 0.
// It returns the converted string and an error if any issues occur during conversion.
func ConvertEBCDICToAsciiNumber(value string, decimalPlaces int) (string, error) {
	if len(value) == 0 {
		value = "0"
	}

	// Check for negative number representation
	runes := []rune(value)
	lastRune := runes[len(runes)-1]
	negative := false

	// Handle negative number if last character is in the map
	if digit, ok := ebcdicToASCIINegativeMap[lastRune]; ok {
		negative = true
		// Replace the last character with its numeric representation
		if len(runes) == 1 {
			value = strconv.Itoa(digit)
		} else {
			value = string(runes[:len(runes)-1]) + strconv.Itoa(digit)
		}
	}

	value = strings.TrimLeft(value, "0")
	if value == "" {
		value = "0"
	}
	minLength := decimalPlaces + 1
	if minLength < 1 {
		minLength = 1
	}
	for len(value) < minLength {
		value = "0" + value
	}

	if decimalPlaces > 0 {
		insertAt := len(value) - decimalPlaces
		value = value[:insertAt] + "." + value[insertAt:]
	}

	if negative {
		value = "-" + value
	}

	return value, nil
}
