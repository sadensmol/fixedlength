package fixedlength

import (
	"sync"
)

type Config struct {
	AlignmentType AlignmentType
}

var once sync.Once
var instance Config

func GetConfig() *Config {
	once.Do(func() {
		instance = Config{
			AlignmentType: AlignmentTypeLeft,
		}
	})

	return &instance
}
