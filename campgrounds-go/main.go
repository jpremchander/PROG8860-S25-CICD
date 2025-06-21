package main

import (
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
	"yelpcamp-go/utils"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("⚠️ Warning: .env file not found, using system environment variables")
	}

	// Ensure JWT secret is configured
	if err := utils.EnsureJWTSecret(); err != nil {
		log.Printf("❌ Failed to configure JWT secret: %v", err)
		os.Exit(1)
	}

	// Initialize database connection
	config.ConnectMongoDB()

	// Seed database with sample data
	db := config.GetDB()
	if err := models.SeedDatabase(db); err != nil {
		log.Printf("⚠️ Warning: Failed to seed database: %v", err)
	}

	// Initialize Gin router
	router := gin.Default()

	// Add CORS middleware
	router.Use(middleware.CORS())
	router.Use(middleware.Logger())

	// Serve static files
	router.Static("/static", "./static")

	// Load HTML templates
	router.LoadHTMLGlob("templates/*")

	// Initialize controllers
	authController := controllers.NewAuthController()
	campgroundController := controllers.NewCampgroundController()
	reviewController := controllers.NewReviewController()

	// Health check endpoints
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":         "healthy",
			"timestamp":      "2025-06-21T14:00:00Z",
			"version":        "1.0.0",
			"jwt_configured": os.Getenv("JWT_SECRET") != "",
		})
	})

	// API routes
	api := router.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":   "API is running",
				"database": "connected",
			})
		})

		// Campgrounds API (using correct method names)
		api.GET("/campgrounds", campgroundController.GetAll)
		api.GET("/campgrounds/:id", campgroundController.GetByID)

		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.AuthRequired())
		{
			protected.POST("/campgrounds", campgroundController.Create)
			protected.PUT("/campgrounds/:id", campgroundController.Update)
			protected.DELETE("/campgrounds/:id", campgroundController.Delete)
		}

		// Auth API routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
		}
	}

	// Web routes
	routes.SetupWebRoutes(router, authController, campgroundController, reviewController)

	// Get port from environment or default to 3000
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("🚀 YelpCamp Go Server starting on port %s", port)
	log.Printf("🌐 Frontend: http://localhost:%s", port)
	log.Printf("📡 API: http://localhost:%s/api/health", port)

	// Start server
	if err := router.Run(":" + port); err != nil {
		log.Printf("❌ Failed to start server: %v", err)
		os.Exit(1)
	}
}
