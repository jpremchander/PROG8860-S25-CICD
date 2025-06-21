package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// EnsureJWTSecret ensures a secure JWT secret is configured
func EnsureJWTSecret() error {
	jwtSecret := os.Getenv("JWT_SECRET")
	
	// Check if JWT secret exists and is secure
	if jwtSecret == "" || len(jwtSecret) < 32 || strings.Contains(jwtSecret, "your_super_secret") {
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
		log.Printf("✅ Using existing JWT secret (%d characters)", len(jwtSecret))
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
func GenerateJWTToken(userID primitive.ObjectID) (string, error) {
	// Ensure JWT secret is configured
	if err := EnsureJWTSecret(); err != nil {
		return "", err
	}
	
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", fmt.Errorf("JWT_SECRET not configured")
	}

	claims := jwt.MapClaims{
		"user_id": userID.Hex(),
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
		"iat":     time.Now().Unix(),
		"iss":     "yelpcamp-go",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jwtSecret))
	
	if err != nil {
		log.Printf("❌ Error signing JWT token: %v", err)
		return "", err
	}
	
	log.Printf("✅ JWT token generated successfully for user: %s", userID.Hex())
	return signedToken, nil
}

// ValidateJWTToken validates a JWT token and returns the user ID
func ValidateJWTToken(tokenString string) (string, error) {
	// Ensure JWT secret is configured
	if err := EnsureJWTSecret(); err != nil {
		return "", err
	}
	
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", fmt.Errorf("JWT_SECRET not configured")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		log.Printf("❌ JWT token parsing error: %v", err)
		return "", err
	}

	if !token.Valid {
		log.Printf("❌ JWT token is invalid")
		return "", fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Printf("❌ JWT token claims are invalid")
		return "", fmt.Errorf("invalid token claims")
	}

	// Validate expiration
	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			log.Printf("❌ JWT token has expired")
			return "", fmt.Errorf("token has expired")
		}
	}

	// Extract user ID
	if userID, ok := claims["user_id"].(string); ok {
		log.Printf("✅ JWT token validated successfully for user: %s", userID)
		return userID, nil
	}

	return "", fmt.Errorf("user_id not found in token claims")
}
