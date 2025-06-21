package middleware

import (
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
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
			c.Abort()
			return
		}

		claims, err := utils.ValidateJWTToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if userIDStr, ok := claims["user_id"].(string); ok {
			if userID, err := primitive.ObjectIDFromHex(userIDStr); err == nil {
				c.Set("user_id", userID)
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
				c.Abort()
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func WebAuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		if err != nil || token == "" {
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		claims, err := utils.ValidateJWTToken(token)
		if err != nil {
			c.SetCookie("token", "", -1, "/", "", false, true)
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		if userIDStr, ok := claims["user_id"].(string); ok {
			if userID, err := primitive.ObjectIDFromHex(userIDStr); err == nil {
				c.Set("user_id", userID)
				c.Set("authenticated", true)
			} else {
				c.SetCookie("token", "", -1, "/", "", false, true)
				c.Redirect(http.StatusSeeOther, "/login")
				c.Abort()
				return
			}
		} else {
			c.SetCookie("token", "", -1, "/", "", false, true)
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		c.Next()
	}
}

func WebAuthOptional() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		if err != nil || token == "" {
			c.Set("authenticated", false)
			c.Next()
			return
		}

		claims, err := utils.ValidateJWTToken(token)
		if err != nil {
			c.SetCookie("token", "", -1, "/", "", false, true)
			c.Set("authenticated", false)
			c.Next()
			return
		}

		if userIDStr, ok := claims["user_id"].(string); ok {
			if userID, err := primitive.ObjectIDFromHex(userIDStr); err == nil {
				c.Set("user_id", userID)
				c.Set("authenticated", true)
			} else {
				c.SetCookie("token", "", -1, "/", "", false, true)
				c.Set("authenticated", false)
			}
		} else {
			c.SetCookie("token", "", -1, "/", "", false, true)
			c.Set("authenticated", false)
		}

		c.Next()
	}
}
