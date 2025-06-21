package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"yelpcamp-go/config"
	"yelpcamp-go/controllers"
	"yelpcamp-go/models"
	"yelpcamp-go/routes"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Set JWT secret if not provided
	if os.Getenv("JWT_SECRET") == "" {
		os.Setenv("JWT_SECRET", "your-super-secret-jwt-key-change-this-in-production")
		log.Println("🔑 JWT_SECRET set to default value")
	}

	// Initialize database
	log.Println("🔌 Connecting to database...")
	config.ConnectDB()

	// Seed database
	log.Println("🌱 Seeding database...")
	if err := models.SeedDatabase(); err != nil {
		log.Printf("⚠️ Warning: Could not seed database: %v", err)
	}

	// Initialize Gin
	r := gin.Default()

	// Load HTML templates
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	// Initialize controllers
	authController := controllers.NewAuthController()
	campgroundController := controllers.NewCampgroundController()
	reviewController := controllers.NewReviewController()

	// Setup routes
	routes.SetupWebRoutes(r, authController, campgroundController, reviewController)
	routes.SetupAPIRoutes(r)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "YelpCamp Go is running!",
			"version": "1.0.0",
		})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("🌐 Access your app at: http://localhost:%s", port)
	
	if err := r.Run(":" + port); err != nil {
		log.Fatal("❌ Failed to start server:", err)
	}
}
