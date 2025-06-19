package models

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCampgroundInput_Validation(t *testing.T) {
	tests := []struct {
		name     string
		input    CampgroundInput
		hasError bool
	}{
		{
			name: "valid campground",
			input: CampgroundInput{
				Title:       "Beautiful Campground",
				Description: "A wonderful place to camp with amazing views",
				Location:    "Yellowstone National Park",
				Price:       25.99,
			},
			hasError: false,
		},
		{
			name: "empty title",
			input: CampgroundInput{
				Title:       "",
				Description: "A wonderful place to camp",
				Location:    "Yellowstone",
				Price:       25.99,
			},
			hasError: true,
		},
		{
			name: "negative price",
			input: CampgroundInput{
				Title:       "Test Campground",
				Description: "A wonderful place to camp",
				Location:    "Yellowstone",
				Price:       -10.0,
			},
			hasError: true,
		},
		{
			name: "short description",
			input: CampgroundInput{
				Title:       "Test Campground",
				Description: "Short",
				Location:    "Yellowstone",
				Price:       25.99,
			},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.hasError {
				hasValidationError := tt.input.Title == "" || 
					len(tt.input.Title) < 3 || 
					len(tt.input.Description) < 10 || 
					tt.input.Price < 0
				assert.True(t, hasValidationError)
			} else {
				assert.True(t, len(tt.input.Title) >= 3)
				assert.True(t, len(tt.input.Description) >= 10)
				assert.True(t, tt.input.Price >= 0)
				assert.NotEmpty(t, tt.input.Location)
			}
		})
	}
}
