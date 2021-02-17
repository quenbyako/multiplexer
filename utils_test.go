package multiplexer

import (
	"reflect"
	"testing"
)

func Test_chunkSlice(t *testing.T) {
	tests := []struct {
		in         []string
		chunkCount int
		expected   [][]string
	}{
		{
			in:         []string{"a", "b", "c"},
			chunkCount: 2,
			expected: [][]string{
				{"a", "b"},
				{"c"},
			},
		},
		{
			in:         []string{"a", "b", "c", "d", "e"},
			chunkCount: 2,
			expected: [][]string{
				{"a", "b", "c"},
				{"d", "e"},
			},
		},
		{
			in:         []string{"a", "b", "c", "d", "e"},
			chunkCount: 3,
			expected: [][]string{
				{"a", "b"},
				{"c", "d"},
				{"e"},
			},
		},
		{
			in:         []string{"a", "b", "c"},
			chunkCount: 100,
			expected: [][]string{
				{"a"},
				{"b"},
				{"c"},
			},
		},
		{
			in:         []string{"a", "b", "c"},
			chunkCount: 1,
			expected: [][]string{
				{"a", "b", "c"},
			},
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := chunkSlice(tt.in, tt.chunkCount)

			// здесь лучше использовать https://github.com/stretchr/testify, для красивых ассертов, но раз
			// нельзя значит нельзя, что поделать
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("got %v, want %v", got, tt.expected)
			}
		})
	}
}

func Test_isDataBinary(t *testing.T) {
	tests := []struct {
		in       []byte
		expected bool
	}{
		{
			in:       []byte("экий некий текст для теста\n"),
			expected: false,
		},
		{
			in:       []byte(""),
			expected: false,
		},
		{
			in:       []byte{0x3c, 0x30, 0x41, 0x7f},
			expected: true,
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := isDataBinary(tt.in)

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("got %v, want %v", got, tt.expected)
			}
		})
	}
}
