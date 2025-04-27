package fixedlength

import (
	"sync"
)

type Config struct {
	AlignmentType            AlignmentType
	NumbersDecimalPlaces     int
	NumbersWithLeadingZeroes bool
}

var once sync.Once
var instance Config

func GetConfig() *Config {
	once.Do(func() {
		instance = Config{
			AlignmentType:            AlignmentTypeLeft,
			NumbersDecimalPlaces:     2,
			NumbersWithLeadingZeroes: true,
		}
	})

	return &instance
}
