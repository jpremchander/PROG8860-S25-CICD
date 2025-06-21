package routes

import (
	"github.com/gin-gonic/gin"
	"campgrounds-go/controllers"
	"campgrounds-go/middleware"
)

func SetupRoutes(r *gin.Engine) {
	// Apply current user middleware to all routes
	r.Use(middleware.SetCurrentUser())

	// Home route
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "home.html", gin.H{
			"title": "YelpCamp",
		})
	})

	// Authentication routes
	r.GET("/register", controllers.ShowRegister)
	r.POST("/register", controllers.Register)
	r.GET("/login", controllers.ShowLogin)
	r.POST("/login", controllers.Login)
	r.POST("/logout", controllers.Logout)

	// Campground routes
	campgrounds := r.Group("/campgrounds")
	{
		campgrounds.GET("", controllers.GetCampgrounds)
		campgrounds.GET("/new", middleware.RequireAuth(), controllers.NewCampground)
		campgrounds.POST("", middleware.RequireAuth(), controllers.CreateCampground)
		campgrounds.GET("/:id", controllers.ShowCampground)
		campgrounds.GET("/:id/edit", middleware.RequireAuth(), controllers.EditCampground)
		campgrounds.PUT("/:id", middleware.RequireAuth(), controllers.UpdateCampground)
		campgrounds.DELETE("/:id", middleware.RequireAuth(), controllers.DeleteCampground)
		
		// Review routes
		campgrounds.POST("/:id/reviews", middleware.RequireAuth(), controllers.CreateReview)
		campgrounds.DELETE("/:id/reviews/:reviewId", middleware.RequireAuth(), controllers.DeleteReview)
	}

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"service": "campgrounds-go",
		})
	})
}
