package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GenerateJWTSecret generates a cryptographically secure JWT secret
func GenerateJWTSecret() string {
	bytes := make([]byte, 32) // 256 bits
	if _, err := rand.Read(bytes); err != nil {
		log.Printf("⚠️ Warning: Could not generate secure random bytes, using fallback")
		// Fallback to a reasonably secure default
		return fmt.Sprintf("yelpcamp_jwt_secret_%d_%d", os.Getpid(), time.Now().Unix())
	}
	return hex.EncodeToString(bytes)
}

// EnsureJWTSecret ensures a JWT secret exists, generates one if needed
func EnsureJWTSecret() string {
	// Check if JWT_SECRET is already set and valid
	if secret := os.Getenv("JWT_SECRET"); secret != "" && len(secret) >= 32 {
		log.Printf("✅ Using existing JWT secret (%d characters)", len(secret))
		return secret
	}

	// Generate new secret
	newSecret := GenerateJWTSecret()
	
	// Set in environment for current session
	os.Setenv("JWT_SECRET", newSecret)
	
	// Try to update .env file for persistence
	if err := updateEnvFile(newSecret); err != nil {
		log.Printf("⚠️ Warning: Could not update .env file: %v", err)
		log.Printf("🔑 Generated JWT Secret: %s", newSecret)
		log.Printf("💡 Consider adding this to your .env file manually")
	} else {
		log.Printf("✅ Generated and saved new JWT secret to .env file")
	}
	
	return newSecret
}

// GenerateJWTToken generates a JWT token for a user - INTEGRATED WITH AUTO SECRET
func GenerateJWTToken(userID primitive.ObjectID) (string, error) {
	// Ensure JWT secret is available
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" || len(jwtSecret) < 32 {
		log.Printf("⚠️ JWT secret not properly configured, ensuring it's set...")
		jwtSecret = EnsureJWTSecret()
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

// ValidateJWTToken validates a JWT token - INTEGRATED WITH AUTO SECRET
func ValidateJWTToken(tokenString string) (jwt.MapClaims, error) {
	// Ensure JWT secret is available
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" || len(jwtSecret) < 32 {
		log.Printf("⚠️ JWT secret not properly configured during validation")
		jwtSecret = EnsureJWTSecret()
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
		return nil, err
	}

	if !token.Valid {
		log.Printf("❌ JWT token is invalid")
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Printf("❌ JWT token claims are invalid")
		return nil, fmt.Errorf("invalid token claims")
	}

	// Validate expiration
	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			log.Printf("❌ JWT token has expired")
			return nil, fmt.Errorf("token has expired")
		}
	}

	log.Printf("✅ JWT token validated successfully")
	return claims, nil
}

// updateEnvFile updates or creates .env file with new JWT secret
func updateEnvFile(secret string) error {
	envPath := ".env"
	
	// Read existing .env file if it exists
	var lines []string
	if content, err := os.ReadFile(envPath); err == nil {
		lines = strings.Split(string(content), "\n")
	} else {
		// Create basic .env structure if file doesn't exist
		lines = []string{
			"# Auto-generated environment configuration",
			"# Database Configuration",
			"MONGO_HOST=localhost",
			"MONGO_PORT=27017",
			"MONGO_DATABASE=yelpcamp_dev",
			"",
			"# Server Configuration",
			"PORT=3000",
			"GIN_MODE=debug",
			"",
			"# Upload Configuration",
			"UPLOAD_PATH=./static/uploads",
			"MAX_UPLOAD_SIZE=5242880",
			"",
			"# JWT Configuration (Auto-generated)",
		}
	}
	
	// Find and update JWT_SECRET line, or add it
	jwtSecretUpdated := false
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "JWT_SECRET=") {
			lines[i] = fmt.Sprintf("JWT_SECRET=%s", secret)
			jwtSecretUpdated = true
			break
		}
	}
	
	// If JWT_SECRET wasn't found, add it
	if !jwtSecretUpdated {
		lines = append(lines, fmt.Sprintf("JWT_SECRET=%s", secret))
	}
	
	// Write back to file
	content := strings.Join(lines, "\n")
	
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(envPath), 0755); err != nil {
		return err
	}
	
	return os.WriteFile(envPath, []byte(content), 0644)
}

// GetEnvironmentType determines if we're in development or production
func GetEnvironmentType() string {
	if ginMode := os.Getenv("GIN_MODE"); ginMode == "release" {
		return "production"
	}
	if env := os.Getenv("ENVIRONMENT"); env != "" {
		return env
	}
	return "development"
}

// GetJWTSecretInfo returns information about the current JWT secret
func GetJWTSecretInfo() map[string]interface{} {
	secret := os.Getenv("JWT_SECRET")
	return map[string]interface{}{
		"configured":    secret != "",
		"length":        len(secret),
		"secure":        len(secret) >= 32,
		"environment":   GetEnvironmentType(),
		"auto_generated": true,
	}
}
