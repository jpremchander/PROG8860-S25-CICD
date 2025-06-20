package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize Gin router
	r := gin.Default()

	// Basic middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "ok",
			"message":   "YelpCamp Go API is running",
			"timestamp": "2024-06-19",
		})
	})

	// Basic API routes
	api := r.Group("/api")
	{
		// Campgrounds routes
		api.GET("/campgrounds", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Campgrounds endpoint working",
				"data":    []gin.H{},
			})
		})

		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Register endpoint working"})
			})
			
			auth.POST("/login", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Login endpoint working"})
			})
		}
	}

	// Root endpoint
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to YelpCamp Go API",
			"version": "1.0.0",
			"endpoints": gin.H{
				"health":      "/health",
				"campgrounds": "/api/campgrounds",
				"auth":        "/api/auth",
			},
		})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("🌐 Health check: http://localhost:%s/health", port)
	log.Printf("📡 API endpoints: http://localhost:%s/api", port)
	
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
