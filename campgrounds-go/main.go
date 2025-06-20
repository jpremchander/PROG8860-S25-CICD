package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Gin router
	r := gin.Default()

	// Basic middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Serve static files
	r.Static("/static", "./static")
	r.Static("/public", "./public")

	// Skip template loading for now - comment this out
	// r.LoadHTMLGlob("templates/*")

	// Routes
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "🏕️ Welcome to YelpCamp Go!",
			"status":  "running",
			"version": "1.0.0",
			"endpoints": []string{
				"/health",
				"/api/health",
				"/api/campgrounds",
			},
		})
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "yelpcamp-go",
			"version": "1.0.0",
		})
	})

	// API routes
	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "ok",
				"message": "YelpCamp Go API is running",
				"version": "1.0.0",
			})
		})

		api.GET("/campgrounds", func(c *gin.Context) {
			campgrounds := []gin.H{
				{"id": 1, "name": "Yellowstone National Park", "location": "Wyoming", "price": 25},
				{"id": 2, "name": "Yosemite Valley", "location": "California", "price": 30},
				{"id": 3, "name": "Grand Canyon", "location": "Arizona", "price": 35},
			}

			c.JSON(http.StatusOK, gin.H{
				"success":     true,
				"campgrounds": campgrounds,
				"count":       len(campgrounds),
			})
		})

		api.GET("/campgrounds/:id", func(c *gin.Context) {
			id := c.Param("id")
			campground := gin.H{
				"id":          id,
				"name":        "Sample Campground",
				"location":    "Beautiful Location",
				"price":       25,
				"description": "A wonderful place to camp.",
			}

			c.JSON(http.StatusOK, gin.H{
				"success":    true,
				"campground": campground,
			})
		})
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("🚀 YelpCamp Go Server starting on port %s", port)
	log.Printf("🌐 Frontend: http://localhost:%s", port)
	log.Printf("📡 API: http://localhost:%s/api/health", port)
	log.Printf("🏕️  Campgrounds: http://localhost:%s/api/campgrounds", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatal("❌ Failed to start server:", err)
	}
}
