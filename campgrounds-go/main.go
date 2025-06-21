package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

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
		log.Println("⚠️ Warning: .env file not found, using system environment variables")
	}

	log.Println("🔗 Connecting to MongoDB...")
	// Connect to MongoDB
	config.ConnectMongoDB()
	defer config.DisconnectMongoDB()

	// Wait for MongoDB to be fully ready
	log.Println("⏳ Waiting for MongoDB to be ready...")
	time.Sleep(5 * time.Second)

	// Test database connection
	db := config.GetDB()
	if db == nil {
		log.Fatal("❌ Database connection failed")
	}

	// Test database operation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	err := db.Client().Ping(ctx, nil)
	if err != nil {
		log.Fatalf("❌ Database ping failed: %v", err)
	}
	log.Println("✅ Database connection verified")

	// **FORCE SEED DATA ON STARTUP**
	log.Println("🌱 Force seeding database with sample data...")
	models.ForceSeedData()

	// Initialize Gin router
	gin.SetMode(gin.DebugMode)
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

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
	reviewController := controllers.NewReviewController()

	// Setup routes
	routes.SetupWebRoutes(router, authController, campgroundController, reviewController)
	routes.SetupAPIRoutes(router, authController, campgroundController, reviewController)

	// Root route - redirect to campgrounds
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusSeeOther, "/campgrounds")
	})

	// Health check routes
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"message":   "YelpCamp Go Server is running",
			"timestamp": time.Now().Format("2006-01-02 15:04:05"),
			"version":   "1.0.0",
		})
	})

	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"message":   "YelpCamp Go API is running",
			"timestamp": time.Now().Format("2006-01-02"),
			"version":   "1.0.0",
		})
	})

	// Public API endpoints (no auth required)
	router.GET("/api/campgrounds", func(c *gin.Context) {
		db := config.GetDB()
		campgrounds, err := models.FindAllCampgrounds(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch campgrounds"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":     "Campgrounds retrieved successfully",
			"count":       len(campgrounds),
			"campgrounds": campgrounds,
		})
	})

	router.GET("/api/campgrounds/:id", func(c *gin.Context) {
		campgroundController.GetByID(c)
	})

	// Debug endpoint to check database contents
	router.GET("/debug/db", func(c *gin.Context) {
		db := config.GetDB()
		
		// Count documents in each collection
		campgroundCount, _ := db.Collection("campgrounds").CountDocuments(context.Background(), map[string]interface{}{})
		userCount, _ := db.Collection("users").CountDocuments(context.Background(), map[string]interface{}{})
		reviewCount, _ := db.Collection("reviews").CountDocuments(context.Background(), map[string]interface{}{})

		// Get sample campgrounds
		campgrounds, _ := models.FindAllCampgrounds(db)

		c.JSON(http.StatusOK, gin.H{
			"database_stats": gin.H{
				"campgrounds": campgroundCount,
				"users":       userCount,
				"reviews":     reviewCount,
			},
			"sample_campgrounds": campgrounds,
			"total_campgrounds":  len(campgrounds),
		})
	})

	// Start server
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
