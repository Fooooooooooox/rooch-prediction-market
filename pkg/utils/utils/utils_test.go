package utils

import (
	"fmt"
	"testing"
	"time"
)

// TestStruct is a sample DTO for testing
type TestStruct struct {
	ID        int
	Name      string
	IsActive  bool
	Price     float64
	CreatedAt time.Time
	Pointer   *string
}

func TestToString(t *testing.T) {
	// Setup test data
	sampleText := "pointer value"
	now := time.Now()

	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name: "basic struct with various types",
			input: TestStruct{
				ID:        123,
				Name:      "test item",
				IsActive:  true,
				Price:     99.99,
				CreatedAt: now,
				Pointer:   &sampleText,
			},
			expected: fmt.Sprintf("TestStruct{ID: 123, Name: test item, IsActive: true, Price: 99.99, CreatedAt: %v, Pointer: %s}", now, sampleText),
		},
		{
			name: "struct with nil pointer",
			input: TestStruct{
				ID:        456,
				Name:      "another test",
				IsActive:  false,
				Price:     0,
				CreatedAt: now,
				Pointer:   nil,
			},
			expected: fmt.Sprintf("TestStruct{ID: 456, Name: another test, IsActive: false, Price: 0, CreatedAt: %v, Pointer: <nil>}", now),
		},
		{
			name:     "nil input",
			input:    nil,
			expected: "nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToString(tt.input)
			if result != tt.expected {
				t.Errorf("ToString() = %v, want %v", result, tt.expected)
			}
		})
	}
}
