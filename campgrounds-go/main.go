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
	if err := config.ConnectDB(); err != nil {
		log.Printf("❌ Failed to connect to database: %v", err)
		os.Exit(1)
	}

	// Seed database with sample data
	db := config.GetDB()
	if err := models.SeedDatabase(db); err != nil {
		log.Printf("⚠️ Warning: Failed to seed database: %v", err)
	}

	// Initialize Gin router
	router := gin.Default()

	// Add CORS middleware
	router.Use(middleware.CORSMiddleware())

	// Serve static files
	router.Static("/static", "./static")
	router.Static("/public", "./public")

	// Load HTML templates
	router.LoadHTMLGlob("templates/*")

	// Initialize controllers
	authController := controllers.NewAuthController()
	campgroundController := controllers.NewCampgroundController()

	// Basic routes
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "YelpCamp - Discover Amazing Campgrounds",
		})
	})

	// Health check endpoints
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": "2025-06-21T12:25:00Z",
			"version":   "1.0.0",
			"jwt_configured": os.Getenv("JWT_SECRET") != "",
		})
	})

	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "API is running",
			"database": "connected",
			"jwt_configured": os.Getenv("JWT_SECRET") != "",
		})
	})

	// API routes
	api := router.Group("/api")
	{
		// Campgrounds API
		api.GET("/campgrounds", campgroundController.GetAllCampgrounds)
		api.GET("/campgrounds/:id", campgroundController.GetCampgroundByID)
		
		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.AuthRequired())
		{
			protected.POST("/campgrounds", campgroundController.CreateCampground)
			protected.PUT("/campgrounds/:id", campgroundController.UpdateCampground)
			protected.DELETE("/campgrounds/:id", campgroundController.DeleteCampground)
		}

		// Auth API routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
		}
	}

	// Web routes
	routes.SetupWebRoutes(router, authController, campgroundController)

	// Get port from environment or default to 3000
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("🚀 YelpCamp Go Server starting on port %s", port)
	log.Printf("🌐 Frontend: http://localhost:%s", port)
	log.Printf("📡 API: http://localhost:%s/api/health", port)
	log.Printf("🏕️  Campgrounds: http://localhost:%s/api/campgrounds", port)

	// Start server
	if err := router.Run(":" + port); err != nil {
		log.Printf("❌ Failed to start server: %v", err)
		os.Exit(1)
	}
}
