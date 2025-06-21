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
		log.Println("No .env file found, using system environment variables")
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

	// ALWAYS force seed data for demo purposes
	log.Println("🌱 Force seeding sample data for demo...")
	models.ForceSeedData()

	// Verify seeding worked with detailed logging
	count, err := db.Collection("campgrounds").CountDocuments(ctx, map[string]interface{}{})
	if err != nil {
		log.Printf("❌ Error checking campground count: %v", err)
	} else {
		log.Printf("📊 Total campgrounds in database: %d", count)
		if count == 0 {
			log.Println("⚠️ WARNING: No campgrounds found after seeding!")
		}
	}

	// Initialize Gin router
	router := gin.Default()

	// Load HTML templates
	router.SetHTMLTemplate(template.Must(template.ParseGlob("templates/*")))

	// Serve static files (including uploads)
	router.Static("/static", "./static")
	router.Static("/uploads", "./static/uploads")

	// CORS middleware
	router.Use(middleware.CORSMiddleware())

	// Basic routes
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "YelpCamp - Discover Amazing Campgrounds",
		})
	})

	router.GET("/health", func(c *gin.Context) {
		// Test database connection
		db := config.GetDB()
		campgroundCount := int64(0)
		userCount := int64(0)
		
		if db != nil {
			count, err := db.Collection("campgrounds").CountDocuments(c.Request.Context(), map[string]interface{}{})
			if err == nil {
				campgroundCount = count
			}
			
			uCount, err := db.Collection("users").CountDocuments(c.Request.Context(), map[string]interface{}{})
			if err == nil {
				userCount = uCount
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"status":      "healthy",
			"message":     "YelpCamp Go server is running",
			"campgrounds": campgroundCount,
			"users":       userCount,
			"timestamp":   time.Now().Format(time.RFC3339),
		})
	})

	router.GET("/api/health", func(c *gin.Context) {
		// Test database connection
		db := config.GetDB()
		campgroundCount := int64(0)
		userCount := int64(0)
		
		if db != nil {
			count, err := db.Collection("campgrounds").CountDocuments(c.Request.Context(), map[string]interface{}{})
			if err == nil {
				campgroundCount = count
			}
			
			uCount, err := db.Collection("users").CountDocuments(c.Request.Context(), map[string]interface{}{})
			if err == nil {
				userCount = uCount
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"status":      "healthy",
			"message":     "YelpCamp API is running",
			"campgrounds": campgroundCount,
			"users":       userCount,
			"timestamp":   time.Now().Format(time.RFC3339),
		})
	})

	// Enhanced API endpoint for campgrounds with better error handling
	router.GET("/api/campgrounds", func(c *gin.Context) {
		db := config.GetDB()
		if db == nil {
			log.Println("❌ Database connection is nil")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":       "Database connection failed",
				"campgrounds": nil,
			})
			return
		}

		log.Println("🔍 Fetching campgrounds from database...")
		campgrounds, err := models.FindAllCampgrounds(db)
		if err != nil {
			log.Printf("❌ Error fetching campgrounds: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":       "Could not fetch campgrounds: " + err.Error(),
				"campgrounds": nil,
			})
			return
		}

		log.Printf("✅ Successfully fetched %d campgrounds", len(campgrounds))
		c.JSON(http.StatusOK, gin.H{
			"campgrounds": campgrounds,
			"count":       len(campgrounds),
			"message":     "Campgrounds retrieved successfully",
		})
	})

	// Debug endpoint to check database contents
	router.GET("/debug/db", func(c *gin.Context) {
		db := config.GetDB()
		if db == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
			return
		}

		// Get collection stats
		campgroundCount, _ := db.Collection("campgrounds").CountDocuments(c.Request.Context(), map[string]interface{}{})
		userCount, _ := db.Collection("users").CountDocuments(c.Request.Context(), map[string]interface{}{})
		reviewCount, _ := db.Collection("reviews").CountDocuments(c.Request.Context(), map[string]interface{}{})

		// Get sample campgrounds
		campgrounds, err := models.FindAllCampgrounds(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"database_stats": gin.H{
				"campgrounds": campgroundCount,
				"users":       userCount,
				"reviews":     reviewCount,
			},
			"sample_campgrounds": campgrounds,
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

	log.Println("🚀 Starting YelpCamp Go Server...")
	log.Printf("🌐 Frontend: http://localhost:%s", port)
	log.Printf("📡 API Health: http://localhost:%s/api/health", port)
	log.Printf("🏕️  Campgrounds API: http://localhost:%s/api/campgrounds", port)
	log.Printf("🔍 Debug DB: http://localhost:%s/debug/db", port)
	log.Printf("🗄️  MongoDB: Connected and seeded")

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
