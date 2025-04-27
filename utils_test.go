package fixedlength

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatStringWithPadding(t *testing.T) {
	tests := []struct {
		name        string
		str         string
		paddingType AlignmentType
		length      int
		wantStr     string
		wantErr     error
	}{
		{
			name:        "right alignment",
			str:         "test_right",
			paddingType: AlignmentTypeRight,
			length:      13,
			wantStr:     "   test_right",
		},
		{
			name:        "left alignment",
			str:         "test_left",
			paddingType: AlignmentTypeLeft,
			length:      12,
			wantStr:     "test_left   ",
		},
		{
			name:        "center alignment",
			str:         "test_center",
			paddingType: AlignmentTypeCenter,
			length:      17,
			wantStr:     "   test_center   ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := FormatStringWithAlignment(tt.str, tt.length, tt.paddingType)
			if err != nil || tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			}
			assert.Equal(t, tt.wantStr, res)
		})
	}
}

func TestFormatFloatWithPadding(t *testing.T) {
	tests := []struct {
		name          string
		num           float64
		numDecimals   int
		leadingZeroes bool
		paddingType   AlignmentType
		length        int

		wantStr string
		wantErr error
	}{
		{
			name:          "left alignment float 2 decimals no leading zeroes",
			num:           12345.12,
			numDecimals:   2,
			leadingZeroes: false,
			paddingType:   AlignmentTypeLeft,
			length:        15,
			wantStr:       "12345.12       ",
		},
		{
			name:          "right alignment float 2 decimals no leading zeroes",
			num:           12345.12,
			numDecimals:   2,
			leadingZeroes: false,
			paddingType:   AlignmentTypeRight,
			length:        15,
			wantStr:       "       12345.12",
		},

		{
			name:          "left alignment float 2 decimals leading zeroes",
			num:           12345.12,
			numDecimals:   2,
			leadingZeroes: true,
			paddingType:   AlignmentTypeLeft,
			length:        15,
			wantStr:       "000000012345.12",
		},
		{
			name:          "right alignment float no decimals leading zeroes",
			num:           12345.12,
			numDecimals:   0,
			leadingZeroes: true,
			paddingType:   AlignmentTypeRight,
			length:        15,
			wantStr:       "000000000012345",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := FormatFloatWithAlignment(tt.num, tt.numDecimals, tt.length, tt.leadingZeroes, tt.paddingType)
			if err != nil || tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			}
			assert.Equal(t, tt.wantStr, res)

		})
	}
}

func TestFormatIntWithPadding(t *testing.T) {
	tests := []struct {
		name          string
		num           int64
		leadingZeroes bool
		paddingType   AlignmentType
		length        int

		wantStr string
		wantErr error
	}{
		{
			name:          "left alignment int with no leading zeroes",
			num:           12345,
			leadingZeroes: false,
			paddingType:   AlignmentTypeLeft,
			length:        15,
			wantStr:       "12345          ",
		},
		{
			name:          "right alignment int with no leading zeroes",
			num:           12345,
			leadingZeroes: false,
			paddingType:   AlignmentTypeRight,
			length:        15,
			wantStr:       "          12345",
		},
		{
			name:          "left alignment int with leading zeroes",
			num:           12345,
			leadingZeroes: true,
			paddingType:   AlignmentTypeLeft,
			length:        15,
			wantStr:       "000000000012345",
		},
		{
			name:          "right alignment int with no leading zeroes",
			num:           12345,
			leadingZeroes: true,
			paddingType:   AlignmentTypeRight,
			length:        15,
			wantStr:       "000000000012345",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := FormatIntWithAlignment(tt.num, tt.length, tt.leadingZeroes, tt.paddingType)
			if err != nil || tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			}
			assert.Equal(t, tt.wantStr, res)

		})
	}
}
