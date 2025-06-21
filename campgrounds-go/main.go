package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Campground represents a campground entity
type Campground struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Price       float64   `json:"price"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	Images      []string  `json:"images"`
	CreatedAt   time.Time `json:"created_at"`
}

// In-memory storage for demonstration
var campgrounds = []Campground{
	{
		ID:          1,
		Title:       "Sunset Valley Campground",
		Price:       25.99,
		Description: "A beautiful campground with stunning sunset views",
		Location:    "Rocky Mountains, Colorado",
		Images:      []string{"https://example.com/sunset1.jpg", "https://example.com/sunset2.jpg"},
		CreatedAt:   time.Now(),
	},
	{
		ID:          2,
		Title:       "Lakeside Retreat",
		Price:       35.50,
		Description: "Peaceful lakeside camping with fishing opportunities",
		Location:    "Lake Tahoe, California",
		Images:      []string{"https://example.com/lake1.jpg"},
		CreatedAt:   time.Now(),
	},
	{
		ID:          3,
		Title:       "Forest Haven",
		Price:       20.00,
		Description: "Deep forest camping for nature lovers",
		Location:    "Olympic National Park, Washington",
		Images:      []string{"https://example.com/forest1.jpg", "https://example.com/forest2.jpg"},
		CreatedAt:   time.Now(),
	},
}

var nextID = 4

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Setup Gin router
	r := gin.Default()

	// Middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// CORS middleware for frontend integration
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// Routes
	setupRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("🚀 YelpCamp Go server starting on port %s", port)
	log.Printf("📍 Environment: %s", os.Getenv("GIN_MODE"))
	log.Printf("🔗 Health check: http://localhost:%s/health", port)
	log.Printf("🏕️  Campgrounds API: http://localhost:%s/api/campgrounds", port)
	
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func setupRoutes(r *gin.Engine) {
	// Home route
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":     "Welcome to YelpCamp Go! 🏕️",
			"status":      "running",
			"version":     "1.0.0",
			"environment": os.Getenv("GIN_MODE"),
			"endpoints": map[string]string{
				"health":      "/health",
				"campgrounds": "/api/campgrounds",
				"docs":        "/api/docs",
			},
		})
	})

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"service":   "campgrounds-go",
			"timestamp": time.Now().UTC(),
			"uptime":    "running",
			"version":   "1.0.0",
		})
	})

	// API documentation
	r.GET("/api/docs", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"api_version": "1.0.0",
			"endpoints": map[string]interface{}{
				"GET /": "Welcome message",
				"GET /health": "Health check",
				"GET /api/campgrounds": "List all campgrounds",
				"GET /api/campgrounds/:id": "Get campground by ID",
				"POST /api/campgrounds": "Create new campground",
				"PUT /api/campgrounds/:id": "Update campground",
				"DELETE /api/campgrounds/:id": "Delete campground",
			},
			"example_campground": Campground{
				ID:          1,
				Title:       "Example Campground",
				Price:       25.99,
				Description: "A sample campground",
				Location:    "Sample Location",
				Images:      []string{"image1.jpg"},
				CreatedAt:   time.Now(),
			},
		})
	})

	// API routes
	api := r.Group("/api")
	{
		// Get all campgrounds
		api.GET("/campgrounds", getCampgrounds)
		
		// Get campground by ID
		api.GET("/campgrounds/:id", getCampgroundByID)
		
		// Create new campground
		api.POST("/campgrounds", createCampground)
		
		// Update campground
		api.PUT("/campgrounds/:id", updateCampground)
		
		// Delete campground
		api.DELETE("/campgrounds/:id", deleteCampground)
	}
}

func getCampgrounds(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"campgrounds": campgrounds,
		"total":       len(campgrounds),
		"message":     "Campgrounds retrieved successfully",
	})
}

func getCampgroundByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campground ID"})
		return
	}

	for _, campground := range campgrounds {
		if campground.ID == id {
			c.JSON(http.StatusOK, gin.H{
				"campground": campground,
				"message":    "Campground found",
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Campground not found"})
}

func createCampground(c *gin.Context) {
	var newCampground Campground
	if err := c.ShouldBindJSON(&newCampground); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate required fields
	if newCampground.Title == "" || newCampground.Description == "" || newCampground.Location == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title, description, and location are required"})
		return
	}

	if newCampground.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Price must be greater than 0"})
		return
	}

	// Set ID and timestamp
	newCampground.ID = nextID
	nextID++
	newCampground.CreatedAt = time.Now()

	// Add to storage
	campgrounds = append(campgrounds, newCampground)

	c.JSON(http.StatusCreated, gin.H{
		"campground": newCampground,
		"message":    "Campground created successfully",
	})
}

func updateCampground(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campground ID"})
		return
	}

	var updateData Campground
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, campground := range campgrounds {
		if campground.ID == id {
			// Update fields
			if updateData.Title != "" {
				campgrounds[i].Title = updateData.Title
			}
			if updateData.Description != "" {
				campgrounds[i].Description = updateData.Description
			}
			if updateData.Location != "" {
				campgrounds[i].Location = updateData.Location
			}
			if updateData.Price > 0 {
				campgrounds[i].Price = updateData.Price
			}
			if len(updateData.Images) > 0 {
				campgrounds[i].Images = updateData.Images
			}

			c.JSON(http.StatusOK, gin.H{
				"campground": campgrounds[i],
				"message":    "Campground updated successfully",
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Campground not found"})
}

func deleteCampground(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campground ID"})
		return
	}

	for i, campground := range campgrounds {
		if campground.ID == id {
			// Remove from slice
			campgrounds = append(campgrounds[:i], campgrounds[i+1:]...)
			c.JSON(http.StatusOK, gin.H{
				"message": fmt.Sprintf("Campground '%s' deleted successfully", campground.Title),
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Campground not found"})
}
