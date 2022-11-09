package trim

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringSlice(t *testing.T) {
	var tests = []struct {
		arg  []string
		want []string
	}{
		{arg: []string{"1", "2", "3"}, want: []string{"1", "2", "3"}},
		{arg: []string{"1", "2", ""}, want: []string{"1", "2"}},
		{arg: []string{"1", "", ""}, want: []string{"1"}},
		{arg: []string{"", "", ""}, want: nil},
		{arg: []string{"   ", "", ""}, want: []string{"   "}},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, StringSlice(tt.arg))
	}
}

func TestNewLine(t *testing.T) {
	var tests = []struct {
		arg  string
		want string
	}{
		{arg: "123", want: "123"},
		{arg: "12\n3", want: "123"},
		{arg: "12\r3", want: "123"},
		{arg: "12\r\n3", want: "123"},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, NewLine(tt.arg))
	}
}
