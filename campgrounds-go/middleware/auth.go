package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		
		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			jwtSecret = "your_super_secret_jwt_key_here"
		}
		
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		userIDStr := claims["user_id"].(string)
		userID, err := primitive.ObjectIDFromHex(userIDStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}

func WebAuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for JWT token in cookie
		tokenString, err := c.Cookie("token")
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			jwtSecret = "your_super_secret_jwt_key_here"
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			// Clear invalid cookie
			c.SetCookie("token", "", -1, "/", "", false, true)
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.SetCookie("token", "", -1, "/", "", false, true)
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			c.SetCookie("token", "", -1, "/", "", false, true)
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		userID, err := primitive.ObjectIDFromHex(userIDStr)
		if err != nil {
			c.SetCookie("token", "", -1, "/", "", false, true)
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		c.Set("user_id", userID)
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

		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			jwtSecret = "your_super_secret_jwt_key_here"
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.Set("authenticated", false)
			c.Next()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.Set("authenticated", false)
			c.Next()
			return
		}

		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			c.Set("authenticated", false)
			c.Next()
			return
		}

		userID, err := primitive.ObjectIDFromHex(userIDStr)
		if err != nil {
			c.Set("authenticated", false)
			c.Next()
			return
		}

		c.Set("user_id", userID)
		c.Set("authenticated", true)
		c.Next()
	}
}
