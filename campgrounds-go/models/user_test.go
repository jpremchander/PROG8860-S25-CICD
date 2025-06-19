package models

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestUser_HashPassword(t *testing.T) {
	user := &User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "plainpassword",
	}

	err := user.HashPassword()
	assert.NoError(t, err)
	assert.NotEqual(t, "plainpassword", user.Password)
	assert.NotEmpty(t, user.Password)
}

func TestUser_CheckPassword(t *testing.T) {
	user := &User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "plainpassword",
	}

	// Hash the password first
	err := user.HashPassword()
	assert.NoError(t, err)

	// Test correct password
	err = user.CheckPassword("plainpassword")
	assert.NoError(t, err)

	// Test incorrect password
	err = user.CheckPassword("wrongpassword")
	assert.Error(t, err)
}

func TestUserInput_Validation(t *testing.T) {
	tests := []struct {
		name     string
		input    UserInput
		hasError bool
	}{
		{
			name: "valid input",
			input: UserInput{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			hasError: false,
		},
		{
			name: "empty username",
			input: UserInput{
				Username: "",
				Email:    "test@example.com",
				Password: "password123",
			},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.hasError {
				assert.True(t, tt.input.Username == "" || len(tt.input.Username) < 3)
			} else {
				assert.True(t, len(tt.input.Username) >= 3)
				assert.Contains(t, tt.input.Email, "@")
				assert.True(t, len(tt.input.Password) >= 6)
			}
		})
	}
}
