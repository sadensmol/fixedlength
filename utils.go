package fixedlength

import "fmt"

type stringPaddingType int

var (
	stringPaddingTypeLeft   stringPaddingType = 1
	stringPaddingTypeRight  stringPaddingType = 2
	stringPaddingTypeCenter stringPaddingType = 3
)

func formatStringWithPadding(str string, length int, paddingType stringPaddingType) string {
	switch paddingType {
	case stringPaddingTypeLeft:
		return fmt.Sprintf("%*s", length, str)
	case stringPaddingTypeRight:
		return fmt.Sprintf("%-*s", length, str)
	case stringPaddingTypeCenter:
		padding := (length - len(str)) / 2
		return fmt.Sprintf("%*s%s%*s", padding, "", str, length-len(str)-padding, "")
	}

	// unknown padding type
	return str
}
