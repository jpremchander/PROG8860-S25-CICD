package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"yelpcamp-go/config"
	"yelpcamp-go/controllers"
	"yelpcamp-go/middleware"
	"yelpcamp-go/models"
	"yelpcamp-go/routes"
	"yelpcamp-go/utils"
)

func main() {
	// Load environment variables first
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ Warning: .env file not found, will create one with generated secrets")
	}

	// Auto-generate JWT secret if not present or too short
	jwtSecret := utils.EnsureJWTSecret()
	log.Printf("🔐 JWT Secret configured (%d characters) for %s environment", 
		len(jwtSecret), utils.GetEnvironmentType())

	// Set other required environment variables with defaults
	setEnvDefault("MONGO_HOST", "localhost")
	setEnvDefault("MONGO_PORT", "27017")
	setEnvDefault("MONGO_DATABASE", "yelpcamp_dev")
	setEnvDefault("PORT", "3000")
	setEnvDefault("GIN_MODE", "debug")
	setEnvDefault("UPLOAD_PATH", "./static/uploads")
	setEnvDefault("MAX_UPLOAD_SIZE", "5242880")

	log.Println("🔗 Connecting to MongoDB...")
	// Connect to MongoDB
	config.ConnectMongoDB()
	defer config.DisconnectMongoDB()

	// Wait for MongoDB to be fully ready
	log.Println("⏳ Waiting for MongoDB to be ready...")
	time.Sleep(3 * time.Second)

	// Test database connection
	db := config.GetDB()
	if db == nil {
		log.Fatal("❌ Database connection failed")
	}

	// Test database operation with retry
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err := db.Client().Ping(ctx, nil)
		cancel()
		
		if err == nil {
			log.Println("✅ Database connection verified")
			break
		}
		
		if i == maxRetries-1 {
			log.Fatalf("❌ Database ping failed after %d attempts: %v", maxRetries, err)
		}
		
		log.Printf("⚠️ Database ping attempt %d failed, retrying...", i+1)
		time.Sleep(2 * time.Second)
	}

	// Force clean and seed database
	log.Println("🧹 Cleaning and seeding database...")
	models.ForceSeedData()

	// Verify seeding worked
	campgroundCount, err := db.Collection("campgrounds").CountDocuments(context.Background(), bson.M{})
	if err != nil {
		log.Printf("⚠️ Warning: Could not verify campground count: %v", err)
	} else {
		log.Printf("✅ Database contains %d campgrounds", campgroundCount)
	}

	// Initialize Gin router
	if utils.GetEnvironmentType() == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Add CORS middleware
	router.Use(middleware.CORS())

	// Serve static files with proper headers
	router.Static("/static", "./static")
	router.StaticFS("/uploads", http.Dir("./static/uploads"))

	// Load HTML templates
	router.LoadHTMLGlob("templates/*")

	// Initialize controllers
	authController := controllers.NewAuthController()
	campgroundController := controllers.NewCampgroundController()
	reviewController := controllers.NewReviewController()

	// Setup routes
	routes.SetupWebRoutes(router, authController, campgroundController, reviewController)
	routes.SetupAPIRoutes(router)

	// Root route - redirect to campgrounds
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusSeeOther, "/campgrounds")
	})

	// Health check routes
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":      "ok",
			"message":     "YelpCamp Go Server is running",
			"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
			"version":     "1.0.0",
			"environment": utils.GetEnvironmentType(),
			"campgrounds": campgroundCount,
			"jwt_configured": len(os.Getenv("JWT_SECRET")) >= 32,
		})
	})

	// Debug endpoint (only in development)
	if utils.GetEnvironmentType() == "development" {
		router.GET("/debug", func(c *gin.Context) {
			db := config.GetDB()
			
			campgroundCount, _ := db.Collection("campgrounds").CountDocuments(context.Background(), bson.M{})
			userCount, _ := db.Collection("users").CountDocuments(context.Background(), bson.M{})
			reviewCount, _ := db.Collection("reviews").CountDocuments(context.Background(), bson.M{})

			// Get sample campgrounds
			campgrounds, _ := models.FindAllCampgrounds(db)

			c.JSON(http.StatusOK, gin.H{
				"database_stats": gin.H{
					"campgrounds": campgroundCount,
					"users":       userCount,
					"reviews":     reviewCount,
				},
				"sample_campgrounds": campgrounds,
				"environment":        utils.GetEnvironmentType(),
				"jwt_secret_length":  len(os.Getenv("JWT_SECRET")),
				"auto_generated":     true,
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
	log.Printf("📡 API: http://localhost:%s/health", port)
	log.Printf("🏕️  Campgrounds: http://localhost:%s/campgrounds", port)
	
	if utils.GetEnvironmentType() == "development" {
		log.Printf("🔍 Debug: http://localhost:%s/debug", port)
		log.Println("")
		log.Println("🔐 Demo Login Credentials:")
		log.Println("   Username: igoswamik")
		log.Println("   Password: password123")
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}

// setEnvDefault sets environment variable if not already set
func setEnvDefault(key, defaultValue string) {
	if os.Getenv(key) == "" {
		os.Setenv(key, defaultValue)
	}
}
