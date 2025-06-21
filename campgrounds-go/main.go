package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Setup Gin router
	r := gin.Default()

	// Basic routes
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to YelpCamp Go!",
			"status":  "running",
		})
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "campgrounds-go",
		})
	})

	r.GET("/campgrounds", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"campgrounds": []map[string]interface{}{
				{
					"id":          1,
					"title":       "Sample Campground",
					"description": "A beautiful campground for testing",
					"price":       25.99,
					"location":    "Test Location",
				},
			},
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
