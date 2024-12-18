package fixedlength

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type intS struct {
	string
}

var _ Marshaler = &intS{}

func (i *intS) Marshal() ([]byte, error) {
	switch i.string {
	case "yes":
		return []byte("42"), nil
	}
	return nil, errors.New("invalid value")
}

func TestMarshal(t *testing.T) {

	t.Run("simple struct", func(t *testing.T) {
		type s1 struct {
			Field1 string  `range:"0,5" flags:"optional"`
			Field2 int     `range:"5,10"`
			Field3 float64 `range:"10,20"`
		}

		res, err := Marshal(&s1{Field1: "hello", Field2: 42, Field3: 3.14})
		require.NoError(t, err)

		require.Equal(t, "hello42   3.14      ", string(res))
	})
	t.Run("struct with internal struct", func(t *testing.T) {

		type s1 struct {
			Field1    string `range:"0,5" flags:"optional"`
			IntStruct intS   `range:"5,10"`
		}

		res, err := Marshal(&s1{Field1: "hello", IntStruct: intS{"yes"}})
		require.NoError(t, err)

		require.Equal(t, "hello42   ", string(res))
	})
}
