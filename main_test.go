package main

import "testing"

func TestValidatePort(t *testing.T) {
	tests := []struct {
		port     int
		expected bool
	}{
		{0, false},
		{1, true},
		{4000, true},
		{65535, true},
		{65536, false},
		{-5, false},
	}

	for _, tt := range tests {
		result := validatePort(tt.port)
		if result != tt.expected {
			t.Errorf("validatePort(%d) = %v; want %v", tt.port, result, tt.expected)
		}
	}
}
