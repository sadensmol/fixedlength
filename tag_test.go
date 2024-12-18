package fixedlength

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseFieldTag(t *testing.T) {
	st := reflect.StructTag(`range:"1,5" flags:"optional"`)
	tag, err := parseFieldTag(st)
	require.NoError(t, err)
	require.Equal(t, 1, tag.fromPos)
	require.Equal(t, 5, tag.toPos)
	require.True(t, tag.flags.optional)
}

func TestParseTag(t *testing.T) {
	tests := []struct {
		name       string
		tag        string
		upperBound int
		wantStart  int
		wantEnd    int
		wantErr    error
	}{
		{
			name:       "valid tag with positive range",
			tag:        "2,5",
			upperBound: 10,
			wantStart:  2,
			wantEnd:    5,
			wantErr:    nil,
		},
		{
			name:       "invalid upper bound",
			tag:        "2,5",
			upperBound: 0,
			wantStart:  0,
			wantEnd:    0,
			wantErr:    ErrTagInvalidUpperBound,
		},
		{
			name:       "valid tag with end capped by upperBound",
			tag:        "1,20",
			upperBound: 15,
			wantStart:  1,
			wantEnd:    15,
			wantErr:    nil,
		},
		{
			name:       "invalid empty tag",
			tag:        "",
			upperBound: 10,
			wantStart:  0,
			wantEnd:    0,
			wantErr:    ErrTagEmpty,
		},
		{
			name:       "invalid tag with non-numeric start",
			tag:        "a,5",
			upperBound: 10,
			wantStart:  0,
			wantEnd:    0,
			wantErr:    ErrTagInvalidRangeValues,
		},
		{
			name:       "invalid tag with non-numeric end",
			tag:        "2,b",
			upperBound: 10,
			wantStart:  0,
			wantEnd:    0,
			wantErr:    ErrTagInvalidRangeValues,
		},
		{
			name:       "ineffectual range: start equals end",
			tag:        "3,3",
			upperBound: 10,
			wantStart:  0,
			wantEnd:    0,
			wantErr:    ErrTagInefectualRange,
		},
		{
			name:       "ineffectual range: start greater than end",
			tag:        "5,2",
			upperBound: 10,
			wantStart:  0,
			wantEnd:    0,
			wantErr:    ErrTagInefectualRange,
		},
		{
			name:       "valid tag with -1 end (upperBound is respected)",
			tag:        "0,-1",
			upperBound: 10,
			wantStart:  0,
			wantEnd:    10,
			wantErr:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStart, gotEnd, err := parseRangeTag(tt.tag)

			// Check for expected errors
			if err != nil && tt.wantErr == nil {
				t.Errorf("expected no error, got %v", err)
			}
			if err == nil && tt.wantErr != nil {
				t.Errorf("expected error %v, got none", tt.wantErr)
			}
			if err != nil && tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("expected error %v, got %v", tt.wantErr, err)
			}

			// Check the start and end values
			if gotStart != tt.wantStart {
				t.Errorf("expected start %d, got %d", tt.wantStart, gotStart)
			}
			if gotEnd != tt.wantEnd {
				t.Errorf("expected end %d, got %d", tt.wantEnd, gotEnd)
			}
		})
	}
}

func TestTag_Validate(t *testing.T) {

	t.Run("valid tag", func(t *testing.T) {
		tag := tag{fromPos: 1, toPos: 5}
		err := tag.Validate(10)
		require.NoError(t, err)
	})
	t.Run("invalid tag: start position negative", func(t *testing.T) {
		tag := tag{fromPos: -1, toPos: 5}
		err := tag.Validate(10)
		require.EqualError(t, err, "invalid range values: range:-1,5 flags:optional:false")
	})

	t.Run("invalid tag: end position less start position", func(t *testing.T) {
		tag := tag{fromPos: 5, toPos: 3}
		err := tag.Validate(10)
		require.EqualError(t, err, "invalid range values: range:5,3 flags:optional:false")
	})
	t.Run("invalid tag: start and end position the same", func(t *testing.T) {
		tag := tag{fromPos: 5, toPos: 5}
		err := tag.Validate(10)
		require.EqualError(t, err, "invalid range values: range:5,5 flags:optional:false")
	})

}
