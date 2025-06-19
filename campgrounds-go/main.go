package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"campgrounds-app/config"
	"campgrounds-app/controllers"
	"campgrounds-app/middleware"
	"campgrounds-app/models"
	"campgrounds-app/routes"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize database
	config.ConnectDB()
	
	// Auto-migrate models
	models.AutoMigrate()

	// Initialize Gin router
	r := gin.Default()

	// Load HTML templates
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./public")

	// Global middleware
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	// Setup routes
	routes.SetupRoutes(r)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	r.Run(":" + port)
}
