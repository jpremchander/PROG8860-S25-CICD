package main

import (
	"encoding/json"
	"fmt"
	"html/template"
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
		Description: "A beautiful campground with stunning sunset views over the Rocky Mountains",
		Location:    "Rocky Mountains, Colorado",
		Images:      []string{"https://images.unsplash.com/photo-1504280390367-361c6d9f38f4?w=400", "https://images.unsplash.com/photo-1441974231531-c6227db76b6e?w=400"},
		CreatedAt:   time.Now(),
	},
	{
		ID:          2,
		Title:       "Lakeside Retreat",
		Price:       35.50,
		Description: "Peaceful lakeside camping with fishing opportunities and crystal clear waters",
		Location:    "Lake Tahoe, California",
		Images:      []string{"https://images.unsplash.com/photo-1506905925346-21bda4d32df4?w=400"},
		CreatedAt:   time.Now(),
	},
	{
		ID:          3,
		Title:       "Forest Haven",
		Price:       20.00,
		Description: "Deep forest camping for nature lovers with hiking trails and wildlife viewing",
		Location:    "Olympic National Park, Washington",
		Images:      []string{"https://images.unsplash.com/photo-1551632811-561732d1e306?w=400", "https://images.unsplash.com/photo-1486022119026-a0a0b4b7b7b7?w=400"},
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

	// Load HTML templates
	r.SetHTMLTemplate(template.Must(template.New("").ParseGlob("templates/*")))
	
	// Serve static files
	r.Static("/static", "./static")

	// Middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// CORS middleware for API
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

	// Setup routes
	setupRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("🚀 YelpCamp Go server starting on port %s", port)
	log.Printf("📍 Environment: %s", os.Getenv("GIN_MODE"))
	log.Printf("🌐 Web UI: http://localhost:%s", port)
	log.Printf("🔗 Health check: http://localhost:%s/health", port)
	log.Printf("🏕️  API: http://localhost:%s/api/campgrounds", port)
	
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func setupRoutes(r *gin.Engine) {
	// Web UI Routes
	r.GET("/", homePage)
	r.GET("/campgrounds", campgroundsPage)
	r.GET("/campgrounds/new", newCampgroundPage)
	r.POST("/campgrounds", createCampgroundWeb)
	r.GET("/campgrounds/:id", showCampgroundPage)
	r.GET("/campgrounds/:id/edit", editCampgroundPage)
	r.POST("/campgrounds/:id", updateCampgroundWeb)
	r.POST("/campgrounds/:id/delete", deleteCampgroundWeb)

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

	// API Routes
	api := r.Group("/api")
	{
		api.GET("/docs", apiDocs)
		api.GET("/campgrounds", getCampgrounds)
		api.GET("/campgrounds/:id", getCampgroundByID)
		api.POST("/campgrounds", createCampground)
		api.PUT("/campgrounds/:id", updateCampground)
		api.DELETE("/campgrounds/:id", deleteCampground)
	}
}

// Web UI Handlers
func homePage(c *gin.Context) {
	c.HTML(http.StatusOK, "home.html", gin.H{
		"title": "YelpCamp - Find Your Perfect Campground",
	})
}

func campgroundsPage(c *gin.Context) {
	c.HTML(http.StatusOK, "campgrounds.html", gin.H{
		"title":       "All Campgrounds",
		"campgrounds": campgrounds,
	})
}

func newCampgroundPage(c *gin.Context) {
	c.HTML(http.StatusOK, "new.html", gin.H{
		"title": "Add New Campground",
	})
}

func showCampgroundPage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"title": "Error",
			"error": "Invalid campground ID",
		})
		return
	}

	for _, campground := range campgrounds {
		if campground.ID == id {
			c.HTML(http.StatusOK, "show.html", gin.H{
				"title":      campground.Title,
				"campground": campground,
			})
			return
		}
	}

	c.HTML(http.StatusNotFound, "error.html", gin.H{
		"title": "Not Found",
		"error": "Campground not found",
	})
}

func editCampgroundPage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"title": "Error",
			"error": "Invalid campground ID",
		})
		return
	}

	for _, campground := range campgrounds {
		if campground.ID == id {
			c.HTML(http.StatusOK, "edit.html", gin.H{
				"title":      "Edit " + campground.Title,
				"campground": campground,
			})
			return
		}
	}

	c.HTML(http.StatusNotFound, "error.html", gin.H{
		"title": "Not Found",
		"error": "Campground not found",
	})
}

func createCampgroundWeb(c *gin.Context) {
	title := c.PostForm("title")
	description := c.PostForm("description")
	location := c.PostForm("location")
	priceStr := c.PostForm("price")
	
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"title": "Error",
			"error": "Invalid price format",
		})
		return
	}

	if title == "" || description == "" || location == "" {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"title": "Error",
			"error": "All fields are required",
		})
		return
	}

	newCampground := Campground{
		ID:          nextID,
		Title:       title,
		Price:       price,
		Description: description,
		Location:    location,
		Images:      []string{"https://images.unsplash.com/photo-1504280390367-361c6d9f38f4?w=400"},
		CreatedAt:   time.Now(),
	}
	nextID++

	campgrounds = append(campgrounds, newCampground)
	c.Redirect(http.StatusSeeOther, "/campgrounds")
}

func updateCampgroundWeb(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"title": "Error",
			"error": "Invalid campground ID",
		})
		return
	}

	title := c.PostForm("title")
	description := c.PostForm("description")
	location := c.PostForm("location")
	priceStr := c.PostForm("price")
	
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"title": "Error",
			"error": "Invalid price format",
		})
		return
	}

	for i, campground := range campgrounds {
		if campground.ID == id {
			campgrounds[i].Title = title
			campgrounds[i].Description = description
			campgrounds[i].Location = location
			campgrounds[i].Price = price
			c.Redirect(http.StatusSeeOther, fmt.Sprintf("/campgrounds/%d", id))
			return
		}
	}

	c.HTML(http.StatusNotFound, "error.html", gin.H{
		"title": "Not Found",
		"error": "Campground not found",
	})
}

func deleteCampgroundWeb(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"title": "Error",
			"error": "Invalid campground ID",
		})
		return
	}

	for i, campground := range campgrounds {
		if campground.ID == id {
			campgrounds = append(campgrounds[:i], campgrounds[i+1:]...)
			c.Redirect(http.StatusSeeOther, "/campgrounds")
			return
		}
	}

	c.HTML(http.StatusNotFound, "error.html", gin.H{
		"title": "Not Found",
		"error": "Campground not found",
	})
}

// API Handlers
func apiDocs(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"api_version": "1.0.0",
		"endpoints": map[string]interface{}{
			"GET /api/campgrounds":     "List all campgrounds",
			"GET /api/campgrounds/:id": "Get campground by ID",
			"POST /api/campgrounds":    "Create new campground",
			"PUT /api/campgrounds/:id": "Update campground",
			"DELETE /api/campgrounds/:id": "Delete campground",
		},
		"web_interface": map[string]string{
			"GET /":                    "Home page",
			"GET /campgrounds":         "View all campgrounds",
			"GET /campgrounds/new":     "Add new campground form",
			"GET /campgrounds/:id":     "View campground details",
			"GET /campgrounds/:id/edit": "Edit campground form",
		},
	})
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

	if newCampground.Title == "" || newCampground.Description == "" || newCampground.Location == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title, description, and location are required"})
		return
	}

	if newCampground.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Price must be greater than 0"})
		return
	}

	newCampground.ID = nextID
	nextID++
	newCampground.CreatedAt = time.Now()
	if len(newCampground.Images) == 0 {
		newCampground.Images = []string{"https://images.unsplash.com/photo-1504280390367-361c6d9f38f4?w=400"}
	}

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
			campgrounds = append(campgrounds[:i], campgrounds[i+1:]...)
			c.JSON(http.StatusOK, gin.H{
				"message": fmt.Sprintf("Campground '%s' deleted successfully", campground.Title),
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Campground not found"})
}
