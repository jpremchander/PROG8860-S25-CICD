package utils

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

// Helper function for testing
func ValidateEmail(email string) bool {
	return len(email) > 0 && 
		   len(email) <= 254 && 
		   emailRegex(email)
}

func emailRegex(email string) bool {
	// Simple email validation for testing
	atCount := 0
	for _, char := range email {
		if char == '@' {
			atCount++
		}
	}
	return atCount == 1 && len(email) > 3
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected bool
	}{
		{
			name:     "valid email",
			email:    "test@example.com",
			expected: true,
		},
		{
			name:     "invalid email - no @",
			email:    "testexample.com",
			expected: false,
		},
		{
			name:     "invalid email - multiple @",
			email:    "test@@example.com",
			expected: false,
		},
		{
			name:     "invalid email - too short",
			email:    "a@b",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateEmail(tt.email)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Helper function for price validation
func ValidatePrice(price float64) bool {
	return price >= 0 && price <= 10000
}

func TestValidatePrice(t *testing.T) {
	tests := []struct {
		name     string
		price    float64
		expected bool
	}{
		{
			name:     "valid price",
			price:    25.99,
			expected: true,
		},
		{
			name:     "zero price",
			price:    0,
			expected: true,
		},
		{
			name:     "negative price",
			price:    -10,
			expected: false,
		},
		{
			name:     "too high price",
			price:    15000,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidatePrice(tt.price)
			assert.Equal(t, tt.expected, result)
		})
	}
}
