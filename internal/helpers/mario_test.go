package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePyramid(t *testing.T) {
	tests := []struct {
		height   int
		expected string
	}{
		{1, "#\n"},
		{2, " #\n##\n"},
		{3, "  #\n ##\n###\n"},
		{4, "   #\n  ##\n ###\n####\n"},
		{8, "       #\n      ##\n     ###\n    ####\n   #####\n  ######\n #######\n########\n"},
	}

	for _, tc := range tests {
		result := GeneratePyramid(tc.height)
		assert.Equal(t, tc.expected, result, "height=%d", tc.height)
	}
}

func TestGenerateDoublePyramid(t *testing.T) {
	tests := []struct {
		height   int
		expected string
	}{
		{1, "#  #\n"},
		{2, " #  #\n##  ##\n"},
		{3, "  #  #\n ##  ##\n###  ###\n"},
		{4, "   #  #\n  ##  ##\n ###  ###\n####  ####\n"},
	}

	for _, tc := range tests {
		result := GenerateDoublePyramid(tc.height)
		assert.Equal(t, tc.expected, result, "height=%d", tc.height)
	}
}
