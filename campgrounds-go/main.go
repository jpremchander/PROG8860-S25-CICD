package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"yelpcamp-go/config"
	"yelpcamp-go/middleware"
	"yelpcamp-go/models"
	"yelpcamp-go/routes"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Connect to MongoDB
	config.ConnectMongoDB()
	defer config.DisconnectMongoDB()

	// Create indexes
	models.AutoMigrate()

	// Seed sample data in development mode
	if os.Getenv("GIN_MODE") != "release" {
		models.SeedData()
	}

	// Initialize Gin router
	r := gin.Default()

	// Load HTML templates
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	// Basic middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())

	// Setup routes
	routes.SetupRoutes(r)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("🚀 YelpCamp Go Server starting on port %s", port)
	log.Printf("🌐 Web interface: http://localhost:%s", port)
	log.Printf("📡 API endpoints: http://localhost:%s/api", port)
	log.Printf("🗄️  MongoDB: %s:%s", os.Getenv("MONGO_HOST"), os.Getenv("MONGO_PORT"))
	
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
