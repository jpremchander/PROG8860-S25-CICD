package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("userID")
		
		if userID == nil {
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}
		
		c.Set("userID", userID)
		c.Next()
	}
}

func SetCurrentUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("userID")
		username := session.Get("username")
		
		if userID != nil {
			c.Set("userID", userID)
			c.Set("username", username)
			c.Set("isLoggedIn", true)
		} else {
			c.Set("isLoggedIn", false)
		}
		
		c.Next()
	}
}
