package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var jwtSecret []byte

func init() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-super-secret-jwt-key-change-this-in-production"
	}
	jwtSecret = []byte(secret)
}

// EnsureJWTSecret ensures a secure JWT secret is configured
func EnsureJWTSecret() error {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" || len(secret) < 32 || strings.Contains(secret, "your_super_secret") {
		log.Printf("⚠️ JWT secret not properly configured, ensuring it's set...")
		
		// Generate a new secure secret
		newSecret, err := generateSecureSecret()
		if err != nil {
			return fmt.Errorf("failed to generate JWT secret: %v", err)
		}
		
		// Set environment variable
		os.Setenv("JWT_SECRET", newSecret)
		
		// Try to save to .env file
		if err := saveToEnvFile("JWT_SECRET", newSecret); err != nil {
			log.Printf("⚠️ Could not save JWT secret to .env file: %v", err)
		} else {
			log.Printf("✅ Generated and saved new JWT secret to .env file")
		}
		
		environment := os.Getenv("ENVIRONMENT")
		if environment == "" {
			environment = "development"
		}
		
		log.Printf("🔐 JWT Secret configured (%d characters) for %s environment", len(newSecret), environment)
	} else {
		log.Printf("✅ Using existing JWT secret (%d characters)", len(secret))
	}
	
	return nil
}

// generateSecureSecret generates a cryptographically secure random string
func generateSecureSecret() (string, error) {
	bytes := make([]byte, 32) // 256 bits
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// saveToEnvFile saves a key-value pair to .env file
func saveToEnvFile(key, value string) error {
	envFile := ".env"
	
	// Read existing content
	content := ""
	if data, err := os.ReadFile(envFile); err == nil {
		content = string(data)
	}
	
	// Check if key already exists
	lines := strings.Split(content, "\n")
	found := false
	for i, line := range lines {
		if strings.HasPrefix(line, key+"=") {
			lines[i] = fmt.Sprintf("%s=%s", key, value)
			found = true
			break
		}
	}
	
	// If not found, append
	if !found {
		if content != "" && !strings.HasSuffix(content, "\n") {
			content += "\n"
		}
		content += fmt.Sprintf("%s=%s\n", key, value)
	} else {
		content = strings.Join(lines, "\n")
	}
	
	// Write back to file
	return os.WriteFile(envFile, []byte(content), 0644)
}

// GenerateJWTToken generates a JWT token for a user
func GenerateJWTToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateJWTToken validates a JWT token and returns the user ID
func ValidateJWTToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
