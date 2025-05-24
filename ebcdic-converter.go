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

// asciiToEBCDICNegativeMap is the reverse of ebcdicToASCIINegativeMap
var asciiToEBCDICNegativeMap = map[int]rune{
	1: 'J',
	2: 'K',
	3: 'L',
	4: 'M',
	5: 'N',
	6: 'O',
	7: 'P',
	8: 'Q',
	9: 'R',
	0: 'ü',
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

// ConvertAsciiToEBCDICNumber converts an ASCII number string to an EBCDIC number string.
// It works in reverse to ConvertEBCDICToAsciiNumber.
// The function handles decimal points and negative numbers, converting them to EBCDIC format.
// For negative numbers, the last digit is replaced with its EBCDIC representation.
// decimalPlaces is used to remove the decimal point and adjust the number format.
// It returns the converted string and an error if any issues occur during conversion.
func ConvertAsciiToEBCDICNumber(value string, decimalPlaces int) (string, error) {
	if value == "" {
		return "", nil
	}
	// Check if the number is negative
	negative := false
	if value[0] == '-' {
		negative = true
		value = value[1:] // Remove the minus sign
	}

	// Handle decimal point
	if decimalPlaces > 0 {
		parts := strings.Split(value, ".")
		if len(parts) > 1 {
			// If decimal part exists, adjust according to specified decimal places
			intPart := parts[0]
			decimalPart := parts[1]

			// Pad decimal part if needed
			for len(decimalPart) < decimalPlaces {
				decimalPart += "0"
			}

			// Truncate decimal part if longer than specified places
			if len(decimalPart) > decimalPlaces {
				decimalPart = decimalPart[:decimalPlaces]
			}

			value = intPart + decimalPart
		} else {
			// No decimal part, add zeros
			for i := 0; i < decimalPlaces; i++ {
				value += "0"
			}
		}
	} else if strings.Contains(value, ".") {
		// If decimalPlaces is 0 but value has a decimal point, remove it
		parts := strings.Split(value, ".")
		value = parts[0]
	}

	// Ensure proper padding for very small decimal values
	if decimalPlaces > 0 && value == "0" && negative {
		// Add leading zeros for cases like -0.01
		for i := 0; i < decimalPlaces; i++ {
			value = "0" + value
		}
	}

	// Ensure we remove leading zeros but not all zeros
	value = strings.TrimLeft(value, "0")
	if value == "" {
		value = "0"
	}

	// If the number is negative, replace the last digit with its EBCDIC representation
	if negative {
		runes := []rune(value)
		lastDigit, err := strconv.Atoi(string(runes[len(runes)-1]))
		if err != nil {
			return "", err
		}

		ebcdicChar, exists := asciiToEBCDICNegativeMap[lastDigit]
		if !exists {
			return "", err
		}

		// Replace last digit with EBCDIC representation
		if len(runes) == 1 {
			value = string(ebcdicChar)
		} else {
			value = string(runes[:len(runes)-1]) + string(ebcdicChar)
		}
	}

	return value, nil
}
