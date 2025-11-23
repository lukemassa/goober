package goober

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHexEncode(t *testing.T) {
	cases := []struct {
		input          byte
		expectedOutput byte
	}{
		{
			input:          4,
			expectedOutput: '4',
		},
		{
			input:          10,
			expectedOutput: 'a',
		},
		{
			input:          15,
			expectedOutput: 'f',
		},
	}
	for _, tt := range cases {
		t.Run(fmt.Sprintf("%v", tt.input), func(t *testing.T) {
			actualOutput := hexEncode(tt.input)
			assert.Equal(t, tt.expectedOutput, actualOutput)
		})
	}
}

func TestGenerateOneLine(t *testing.T) {
	cases := []struct {
		input          string
		expectedOutput string
	}{
		{
			input:          "AAAAAAAAAAAAAAAA",
			expectedOutput: "41 41 41 41 41 41 41 41  41 41 41 41 41 41 41 41  |AAAAAAAAAAAAAAAA|\n",
		},
		{
			input:          "\xbd\xb2=\xbc \u2318",
			expectedOutput: "bd b2 3d bc 20 e2 8c 98                           |..=. ...|\n",
		},
	}
	for _, tt := range cases {
		t.Run(tt.input, func(t *testing.T) {
			actualOutput := string(generateOneLine([]byte(tt.input)))
			assert.Equal(t, tt.expectedOutput, actualOutput)
		})
	}
}
