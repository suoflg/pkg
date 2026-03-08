package util

import (
	"testing"
)

func TestStringToBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []byte
	}{
		{
			name:     "empty string",
			input:    "",
			expected: []byte{},
		},
		{
			name:     "ascii string",
			input:    "hello",
			expected: []byte{'h', 'e', 'l', 'l', 'o'},
		},
		{
			name:     "utf8 string",
			input:    "你好世界",
			expected: []byte{0xe4, 0xbd, 0xa0, 0xe5, 0xa5, 0xbd, 0xe4, 0xb8, 0x96, 0xe7, 0x95, 0x8c},
		},
		{
			name:     "mixed string",
			input:    "Hello 世界 123",
			expected: []byte{'H', 'e', 'l', 'l', 'o', ' ', 0xe4, 0xb8, 0x96, 0xe7, 0x95, 0x8c, ' ', '1', '2', '3'},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StringToBytes(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("StringToBytes() length = %v, want %v", len(result), len(tt.expected))
			}
			for i, b := range result {
				if b != tt.expected[i] {
					t.Errorf("StringToBytes()[%d] = %v, want %v", i, b, tt.expected[i])
				}
			}
		})
	}
}

func TestBytesToString(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "empty bytes",
			input:    []byte{},
			expected: "",
		},
		{
			name:     "ascii bytes",
			input:    []byte{'h', 'e', 'l', 'l', 'o'},
			expected: "hello",
		},
		{
			name:     "utf8 bytes",
			input:    []byte{0xe4, 0xbd, 0xa0, 0xe5, 0xa5, 0xbd, 0xe4, 0xb8, 0x96, 0xe7, 0x95, 0x8c},
			expected: "你好世界",
		},
		{
			name:     "mixed bytes",
			input:    []byte{'H', 'e', 'l', 'l', 'o', ' ', 0xe4, 0xb8, 0x96, 0xe7, 0x95, 0x8c, ' ', '1', '2', '3'},
			expected: "Hello 世界 123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BytesToString(tt.input)
			if result != tt.expected {
				t.Errorf("BytesToString() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestStringToBytesAndBack(t *testing.T) {
	original := "Hello 世界 123"
	bytes := StringToBytes(original)
	result := BytesToString(bytes)
	if result != original {
		t.Errorf("StringToBytes -> BytesToString = %v, want %v", result, original)
	}
}

func TestBytesToStringAndBack(t *testing.T) {
	original := []byte{'H', 'e', 'l', 'l', 'o', ' ', 0xe4, 0xb8, 0x96, 0xe7, 0x95, 0x8c}
	str := BytesToString(original)
	result := StringToBytes(str)
	for i, b := range result {
		if b != original[i] {
			t.Errorf("BytesToString -> StringToBytes[%d] = %v, want %v", i, b, original[i])
		}
	}
}
