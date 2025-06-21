package routes

import (
	"github.com/gin-gonic/gin"
	"yelpcamp-go/controllers"
	"yelpcamp-go/middleware"
)

	// Initialize controllers
	authController := controllers.NewAuthController()
	campgroundController := controllers.NewCampgroundController()
	reviewController := controllers.NewReviewController()

	// API routes
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

		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
		}

		// Campground routes
		campgrounds := api.Group("/campgrounds")
		{
			campgrounds.GET("", campgroundController.GetAll)
			campgrounds.GET("/:id", campgroundController.GetByID)
			
			// Protected routes
			campgrounds.POST("", middleware.AuthRequired(), campgroundController.Create)
			campgrounds.PUT("/:id", middleware.AuthRequired(), campgroundController.Update)
			campgrounds.DELETE("/:id", middleware.AuthRequired(), campgroundController.Delete)
			campgrounds.POST("/:id/images", middleware.AuthRequired(), campgroundController.UploadImages)
			
			// Review routes
			campgrounds.POST("/:id/reviews", middleware.AuthRequired(), reviewController.Create)
			campgrounds.DELETE("/:id/reviews/:reviewId", middleware.AuthRequired(), reviewController.Delete)
		}

		// User routes
		users := api.Group("/users")
		{
			users.GET("/:id", func(c *gin.Context) {
				// Get user profile API endpoint
				c.JSON(200, gin.H{"message": "User profile endpoint"})
			})
		}
	}
