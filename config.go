package fixedlength

import (
	"sync"
)

type Config struct {
	AlignmentType            AlignmentType
	NumbersWithLeadingZeroes bool
}

var once sync.Once
var instance Config

func GetConfig() *Config {
	once.Do(func() {
		instance = Config{
			AlignmentType:            AlignmentTypeLeft,
			NumbersWithLeadingZeroes: true,
		}
	})

	return &instance
}
