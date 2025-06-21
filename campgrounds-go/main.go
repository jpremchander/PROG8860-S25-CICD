package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"yelpcamp-go/config"
	"yelpcamp-go/controllers"
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

	// Seed sample data for development
	if os.Getenv("GIN_MODE") != "release" {
		log.Println("🌱 Development mode: seeding sample data...")
		models.SeedData()
	}

	// Initialize Gin router
	router := gin.Default()

	// Load HTML templates
	router.SetHTMLTemplate(template.Must(template.ParseGlob("templates/*")))

	// Serve static files
	router.Static("/static", "./static")
	router.Static("/public", "./static")

	// CORS middleware
	router.Use(middleware.CORSMiddleware())

	// Basic routes
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "YelpCamp - Discover Amazing Campgrounds",
		})
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"message": "YelpCamp Go server is running",
		})
	})

	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"message":   "YelpCamp API is running",
			"timestamp": "2024-01-01T00:00:00Z",
		})
	})

	// Simple API endpoint for campgrounds
	router.GET("/api/campgrounds", func(c *gin.Context) {
		db := config.GetDB()
		campgrounds, err := models.FindAllCampgrounds(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch campgrounds"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"campgrounds": campgrounds,
			"count":       len(campgrounds),
		})
	})

	router.GET("/api/campgrounds/:id", func(c *gin.Context) {
		// Simple campground by ID endpoint
		c.JSON(http.StatusOK, gin.H{
			"message": "Campground details endpoint",
			"id":      c.Param("id"),
		})
	})

	// Initialize controllers
	authController := controllers.NewAuthController()
	campgroundController := controllers.NewCampgroundController()
	reviewController := controllers.NewReviewController()

	// Setup routes
	routes.SetupWebRoutes(router, authController, campgroundController, reviewController)
	routes.SetupAPIRoutes(router, authController, campgroundController, reviewController)

	// Get port from environment or default to 3000
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("🚀 YelpCamp Go Server starting on port %s", port)
	log.Printf("🌐 Frontend: http://localhost:%s", port)
	log.Printf("📡 API: http://localhost:%s/api/health", port)
	log.Printf("🏕️  Campgrounds: http://localhost:%s/api/campgrounds", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
