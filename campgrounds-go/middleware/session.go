package middleware

import (
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func SessionMiddleware() gin.HandlerFunc {
	secret := os.Getenv("SESSION_SECRET")
	if secret == "" {
		secret = "your-secret-key-change-in-production"
	}

	store := cookie.NewStore([]byte(secret))
	store.Options(sessions.Options{
		MaxAge:   86400 * 7, // 7 days
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
	})

	return sessions.Sessions("campgrounds-session", store)
}

func FlashMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		
		// Get flash messages
		success := session.Get("success")
		error := session.Get("error")
		
		// Clear flash messages
		session.Delete("success")
		session.Delete("error")
		session.Save()
		
		// Set in context for templates
		if success != nil {
			c.Set("success", success)
		}
		if error != nil {
			c.Set("error", error)
		}
		
		c.Next()
	}
}
