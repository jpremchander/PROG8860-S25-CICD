package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"yelpcamp-go/config"
	"yelpcamp-go/controllers"
	"yelpcamp-go/middleware"
	"yelpcamp-go/models"
)

// SetupRoutes is the main function to setup all routes
func SetupRoutes(r *gin.Engine) {
	// Initialize controllers
	authController := controllers.NewAuthController()
	campgroundController := controllers.NewCampgroundController()
	reviewController := controllers.NewReviewController()

	// Add health check at root level for convenience
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "ok",
			"message":   "YelpCamp Go API is running",
			"timestamp": "2024-06-21",
			"version":   "1.0.0",
		})
	})

	// Setup web and API routes
	setupWebRoutes(r, authController, campgroundController, reviewController)
	setupAPIRoutes(r)
}

// setupWebRoutes handles all web (HTML) routes
func setupWebRoutes(r *gin.Engine, authController *controllers.AuthController, campgroundController *controllers.CampgroundController, reviewController *controllers.ReviewController) {
	// Apply optional auth to all web routes to check if user is logged in
	r.Use(middleware.WebAuthOptional())

	// Homepage
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":         "YelpCamp Go",
			"message":       "Welcome to YelpCamp - Discover Amazing Campgrounds!",
			"authenticated": c.GetBool("authenticated"),
		})
	})

	// Campgrounds listing - PUBLIC but shows different content based on auth
	r.GET("/campgrounds", func(c *gin.Context) {
		db := config.GetDB()
		campgrounds, err := models.FindAllCampgrounds(db)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"title": "Error",
				"error": "Could not load campgrounds",
			})
			return
		}

		c.HTML(http.StatusOK, "campgrounds.html", gin.H{
			"title":         "All Campgrounds",
			"campgrounds":   campgrounds,
			"authenticated": c.GetBool("authenticated"),
		})
	})

	// Individual campground - PUBLIC but shows different content based on auth
	r.GET("/campgrounds/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := primitive.ObjectIDFromHex(idParam)
		if err != nil {
			c.HTML(http.StatusBadRequest, "error.html", gin.H{
				"title": "Error",
				"error": "Invalid campground ID",
			})
			return
		}

		db := config.GetDB()
		campground, err := models.FindCampgroundByID(db, id)
		if err != nil {
			c.HTML(http.StatusNotFound, "error.html", gin.H{
				"title": "Error",
				"error": "Campground not found",
			})
			return
		}

		// Get reviews for this campground
		reviews, _ := models.FindReviewsByCampground(db, id)

		c.HTML(http.StatusOK, "show.html", gin.H{
			"title":         "Campground Details",
			"campground":    campground,
			"reviews":       reviews,
			"authenticated": c.GetBool("authenticated"),
		})
	})

	// Authentication pages and actions
	r.GET("/register", func(c *gin.Context) {
		if c.GetBool("authenticated") {
			c.Redirect(http.StatusSeeOther, "/campgrounds")
			return
		}
		c.HTML(http.StatusOK, "register.html", gin.H{
			"title": "Register - YelpCamp",
		})
	})

	r.GET("/login", func(c *gin.Context) {
		if c.GetBool("authenticated") {
			c.Redirect(http.StatusSeeOther, "/campgrounds")
			return
		}
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title": "Login - YelpCamp",
		})
	})

	// Authentication actions
	r.POST("/register", authController.RegisterWeb)
	r.POST("/login", authController.LoginWeb)
	r.GET("/logout", authController.Logout)
	r.POST("/logout", authController.Logout)

	// Protected routes
	r.GET("/campgrounds/new", middleware.WebAuthRequired(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "new.html", gin.H{
			"title":         "Add New Campground",
			"authenticated": true,
		})
	})

	r.POST("/campgrounds", middleware.WebAuthRequired(), campgroundController.CreateWeb)
	r.POST("/campgrounds/:id/reviews", middleware.WebAuthRequired(), reviewController.CreateWeb)
}

// setupAPIRoutes handles all API (JSON) routes
func setupAPIRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		// Health check
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":    "ok",
				"message":   "YelpCamp Go API is running",
				"timestamp": "2024-06-21",
				"version":   "1.0.0",
			})
		})

		// Test endpoint
		api.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Test endpoint working",
				"data": gin.H{
					"server_time": "2024-06-21T12:00:00Z",
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

		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", controllers.NewAuthController().Register)
			auth.POST("/login", controllers.NewAuthController().Login)
		}

		// Campground routes
		campgrounds := api.Group("/campgrounds")
		{
			campgrounds.GET("", func(c *gin.Context) {
				db := config.GetDB()
				campgroundsList, err := models.FindAllCampgrounds(db)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch campgrounds"})
					return
				}

				c.JSON(http.StatusOK, gin.H{
					"message":     "Campgrounds retrieved successfully",
					"campgrounds": campgroundsList,
					"count":       len(campgroundsList),
				})
			})

			campgrounds.GET("/:id", func(c *gin.Context) {
				idParam := c.Param("id")
				id, err := primitive.ObjectIDFromHex(idParam)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campground ID"})
					return
				}

				db := config.GetDB()
				campground, err := models.FindCampgroundByID(db, id)
				if err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "Campground not found"})
					return
				}

				c.JSON(http.StatusOK, gin.H{
					"message":    "Campground retrieved successfully",
					"campground": campground,
				})
			})

			// Protected routes
			campgrounds.POST("", middleware.AuthRequired(), controllers.NewCampgroundController().Create)
			campgrounds.PUT("/:id", middleware.AuthRequired(), controllers.NewCampgroundController().Update)
			campgrounds.DELETE("/:id", middleware.AuthRequired(), controllers.NewCampgroundController().Delete)
		}
	}
}
