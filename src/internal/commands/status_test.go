package commands

import "testing"

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		name     string
		given    int64
		expected string
	}{
		{"zero bytes", 0, "0 B"},
		{"500 bytes", 500, "500 B"},
		{"1 KB", 1024, "1.0 KB"},
		{"1.5 KB", 1536, "1.5 KB"},
		{"1 MB", 1048576, "1.0 MB"},
		{"1 GB", 1073741824, "1.0 GB"},
		{"1.5 GB", 1610612736, "1.5 GB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given: byte size
			bytes := tt.given

			// When: formatting bytes
			result := formatBytes(bytes)

			// Then: it should return correct string
			if result != tt.expected {
				t.Errorf("formatBytes(%d) = %q, want %q", bytes, result, tt.expected)
			}
		})
	}
}
