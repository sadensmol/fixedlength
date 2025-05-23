package fixedlength

import (
	"fmt"
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

func EbcdicToAsciiNumber(ebcdicField string) (string, error) {
	if len(ebcdicField) == 0 {
		return "", fmt.Errorf("empty EBCDIC field")
	}

	// Check only the last character for EBCDIC negative sign
	lastChar := rune(ebcdicField[len(ebcdicField)-1])
	if digit, ok := ebcdicToASCIINegativeMap[lastChar]; ok {
		// It's a negative number in EBCDIC format
		numPart := ebcdicField[:len(ebcdicField)-1]

		// If just a single EBCDIC character
		if numPart == "" {
			return "-" + strconv.Itoa(digit), nil
		}

		// Handle float values (containing a decimal point)
		if strings.Contains(numPart, ".") {
			// Remove leading zeros before decimal point
			parts := strings.Split(numPart, ".")
			intPart := parts[0]

			start := 0
			for start < len(intPart) && intPart[start] == '0' {
				start++
			}

			// All zeros in integer part
			if start == len(intPart) {
				intPart = "0"
			} else {
				intPart = intPart[start:]
			}

			// Reconstruct with decimal part and EBCDIC digit
			if len(parts) > 1 {
				return "-" + intPart + "." + parts[1] + strconv.Itoa(digit), nil
			}
			return "-" + intPart + strconv.Itoa(digit), nil
		}

		// Handle integer values
		// Remove leading zeros
		start := 0
		for start < len(numPart) && numPart[start] == '0' {
			start++
		}

		// All zeros case
		if start == len(numPart) {
			return "-" + strconv.Itoa(digit), nil
		}

		result := numPart[start:] + strconv.Itoa(digit)
		return "-" + result, nil
	}

	// If last character is not an EBCDIC sign, return the original
	return ebcdicField, nil
}
