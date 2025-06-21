package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"yelpcamp-go/utils"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		
		// Use the integrated JWT validation
		userID, err := utils.ValidateJWTToken(tokenString)
		if err != nil {
			log.Printf("❌ API Auth failed: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		userObjectID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
			c.Abort()
			return
		}

		c.Set("user_id", userObjectID)
		c.Next()
	}
}

func WebAuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for JWT token in cookie
		tokenString, err := c.Cookie("token")
		if err != nil {
			log.Printf("🔐 No token cookie found, redirecting to login")
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		// Use the integrated JWT validation
		userID, err := utils.ValidateJWTToken(tokenString)
		if err != nil {
			log.Printf("🔐 Invalid token, clearing cookie and redirecting to login: %v", err)
			// Clear invalid cookie
			c.SetCookie("token", "", -1, "/", "", false, true)
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		userObjectID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			log.Printf("🔐 Invalid user ID format: %s", userID)
			c.SetCookie("token", "", -1, "/", "", false, true)
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		log.Printf("✅ Web Authentication successful for user ID: %s", userObjectID.Hex())
		c.Set("user_id", userObjectID)
		c.Set("authenticated", true)
		c.Next()
	}
}

// Optional auth - doesn't redirect if not authenticated
func WebAuthOptional() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("token")
		if err != nil {
			c.Set("authenticated", false)
			c.Next()
			return
		}

		// Use the integrated JWT validation
		userID, err := utils.ValidateJWTToken(tokenString)
		if err != nil {
			c.Set("authenticated", false)
			c.Next()
			return
		}

		userObjectID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.Set("authenticated", false)
			c.Next()
			return
		}

		c.Set("user_id", userObjectID)
		c.Set("authenticated", true)
		c.Next()
	}
}
