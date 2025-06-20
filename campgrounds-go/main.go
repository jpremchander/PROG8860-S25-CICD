package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize Gin router
	r := gin.Default()

	// Load HTML templates
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	// Basic middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Web Routes (HTML pages)
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":   "YelpCamp Go",
			"message": "Welcome to YelpCamp - Discover Amazing Campgrounds!",
		})
	})

	r.GET("/campgrounds", func(c *gin.Context) {
		// Sample campgrounds data with working placeholder images
		campgrounds := []gin.H{
			{
				"id":          1,
				"title":       "Sunset Valley Campground",
				"description": "A beautiful campground with stunning sunset views over the valley. Perfect for families and nature lovers.",
				"location":    "Yosemite National Park, CA",
				"price":       35.99,
				"image":       "https://images.unsplash.com/photo-1504280390367-361c6d9f38f4?w=400&h=300&fit=crop",
			},
			{
				"id":          2,
				"title":       "Mountain Peak Retreat",
				"description": "High altitude camping with breathtaking mountain views. Ideal for experienced campers seeking adventure.",
				"location":    "Rocky Mountain National Park, CO",
				"price":       45.50,
				"image":       "https://images.unsplash.com/photo-1571863533956-01c88e79957e?w=400&h=300&fit=crop",
			},
			{
				"id":          3,
				"title":       "Lakeside Paradise",
				"description": "Peaceful lakeside camping with crystal clear waters. Great for fishing, swimming, and relaxation.",
				"location":    "Lake Tahoe, CA",
				"price":       40.00,
				"image":       "https://images.unsplash.com/photo-1441974231531-c6227db76b6e?w=400&h=300&fit=crop",
			},
		}

		c.HTML(http.StatusOK, "campgrounds.html", gin.H{
			"title":       "All Campgrounds",
			"campgrounds": campgrounds,
		})
	})

	r.GET("/campgrounds/:id", func(c *gin.Context) {
		id := c.Param("id")
		
		// Sample campground detail (in real app, fetch from database)
		campground := gin.H{
			"id":          id,
			"title":       "Sunset Valley Campground",
			"description": "A beautiful campground with stunning sunset views over the valley. Perfect for families and nature lovers. This campground offers modern amenities including clean restrooms, hot showers, fire pits, and picnic tables. The site is surrounded by hiking trails and offers easy access to some of the most spectacular viewpoints in the area.",
			"location":    "Yosemite National Park, CA",
			"price":       35.99,
			"image":       "https://images.unsplash.com/photo-1504280390367-361c6d9f38f4?w=800&h=600&fit=crop",
			"amenities":   []string{"Fire Pits", "Restrooms", "Showers", "Picnic Tables", "Hiking Trails"},
		}

		c.HTML(http.StatusOK, "show.html", gin.H{
			"title":      "Campground Details",
			"campground": campground,
		})
	})

	r.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"title": "Register - YelpCamp",
		})
	})

	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title": "Login - YelpCamp",
		})
	})

	// API Routes (JSON responses)
	api := r.Group("/api")
	{
		// Health check
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":    "ok",
				"message":   "YelpCamp Go API is running",
				"timestamp": "2024-06-19",
				"version":   "1.0.0",
			})
		})

		// Campgrounds API
		api.GET("/campgrounds", func(c *gin.Context) {
			campgrounds := []gin.H{
				{
					"id":          1,
					"title":       "Sunset Valley Campground",
					"description": "A beautiful campground with stunning sunset views",
					"location":    "Yosemite National Park, CA",
					"price":       35.99,
					"image":       "https://images.unsplash.com/photo-1504280390367-361c6d9f38f4?w=400&h=300&fit=crop",
				},
				{
					"id":          2,
					"title":       "Mountain Peak Retreat",
					"description": "High altitude camping with breathtaking mountain views",
					"location":    "Rocky Mountain National Park, CO",
					"price":       45.50,
					"image":       "https://images.unsplash.com/photo-1571863533956-01c88e79957e?w=400&h=300&fit=crop",
				},
				{
					"id":          3,
					"title":       "Lakeside Paradise",
					"description": "Peaceful lakeside camping with crystal clear waters",
					"location":    "Lake Tahoe, CA",
					"price":       40.00,
					"image":       "https://images.unsplash.com/photo-1441974231531-c6227db76b6e?w=400&h=300&fit=crop",
				},
			}

			c.JSON(200, gin.H{
				"message": "Campgrounds retrieved successfully",
				"data":    campgrounds,
				"count":   len(campgrounds),
			})
		})

		api.GET("/campgrounds/:id", func(c *gin.Context) {
			id := c.Param("id")
			
			campground := gin.H{
				"id":          id,
				"title":       "Sunset Valley Campground",
				"description": "A beautiful campground with stunning sunset views over the valley",
				"location":    "Yosemite National Park, CA",
				"price":       35.99,
				"image":       "https://images.unsplash.com/photo-1504280390367-361c6d9f38f4?w=800&h=600&fit=crop",
				"amenities":   []string{"Fire Pits", "Restrooms", "Showers", "Picnic Tables", "Hiking Trails"},
			}

			c.JSON(200, gin.H{
				"message":    "Campground retrieved successfully",
				"campground": campground,
			})
		})

		// Auth API
		auth := api.Group("/auth")
		{
			auth.POST("/register", func(c *gin.Context) {
				var user struct {
					Username string `json:"username" binding:"required"`
					Email    string `json:"email" binding:"required,email"`
					Password string `json:"password" binding:"required,min=6"`
				}

				if err := c.ShouldBindJSON(&user); err != nil {
					c.JSON(400, gin.H{
						"error":   "Validation failed",
						"details": err.Error(),
					})
					return
				}

				// Simulate user registration
				c.JSON(201, gin.H{
					"message": "User registered successfully",
					"user": gin.H{
						"id":       123,
						"username": user.Username,
						"email":    user.Email,
					},
					"token": "sample_jwt_token_here",
				})
			})
			
			auth.POST("/login", func(c *gin.Context) {
				var credentials struct {
					Username string `json:"username" binding:"required"`
					Password string `json:"password" binding:"required"`
				}

				if err := c.ShouldBindJSON(&credentials); err != nil {
					c.JSON(400, gin.H{
						"error":   "Validation failed",
						"details": err.Error(),
					})
					return
				}

				// Simulate login validation
				if credentials.Username == "demo" && credentials.Password == "password" {
					c.JSON(200, gin.H{
						"message": "Login successful",
						"user": gin.H{
							"id":       123,
							"username": credentials.Username,
							"email":    "demo@yelpcamp.com",
						},
						"token": "sample_jwt_token_here",
					})
				} else {
					c.JSON(401, gin.H{
						"error": "Invalid credentials",
					})
				}
			})
		}

		// Test endpoints for validation
		api.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Test endpoint working",
				"data": gin.H{
					"server_time": "2024-06-19T10:00:00Z",
					"endpoints": []string{
						"/api/health",
						"/api/campgrounds",
						"/api/campgrounds/:id",
						"/api/auth/register",
						"/api/auth/login",
					},
				},
			})
		})
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("🌐 Web interface: http://localhost:%s", port)
	log.Printf("📡 API endpoints: http://localhost:%s/api", port)
	
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
