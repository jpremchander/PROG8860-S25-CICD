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

func SetupWebRoutes(r *gin.Engine) {
	// Initialize controllers
	authController := controllers.NewAuthController()
	campgroundController := controllers.NewCampgroundController()
	reviewController := controllers.NewReviewController()

	// Homepage
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":   "YelpCamp Go",
			"message": "Welcome to YelpCamp - Discover Amazing Campgrounds!",
		})
	})

	// Campgrounds listing
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
			"title":       "All Campgrounds",
			"campgrounds": campgrounds,
		})
	})

	// Individual campground
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
			"title":      "Campground Details",
			"campground": campground,
			"reviews":    reviews,
		})
	})

	// New campground form (protected)
	r.GET("/campgrounds/new", middleware.WebAuthRequired(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "new.html", gin.H{
			"title": "Add New Campground",
		})
	})

	// Create campground (protected)
	r.POST("/campgrounds", middleware.WebAuthRequired(), campgroundController.CreateWeb)

	// Edit campground form (protected)
	r.GET("/campgrounds/:id/edit", middleware.WebAuthRequired(), func(c *gin.Context) {
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

		// Check if user owns this campground
		userID := c.MustGet("user_id").(primitive.ObjectID)
		if campground.AuthorID != userID {
			c.HTML(http.StatusForbidden, "error.html", gin.H{
				"title": "Error",
				"error": "Not authorized to edit this campground",
			})
			return
		}

		c.HTML(http.StatusOK, "edit.html", gin.H{
			"title":      "Edit Campground",
			"campground": campground,
		})
	})

	// Update campground (protected)
	r.PUT("/campgrounds/:id", middleware.WebAuthRequired(), campgroundController.UpdateWeb)

	// Delete campground (protected)
	r.DELETE("/campgrounds/:id", middleware.WebAuthRequired(), campgroundController.DeleteWeb)

	// User profile
	r.GET("/users/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := primitive.ObjectIDFromHex(idParam)
		if err != nil {
			c.HTML(http.StatusBadRequest, "error.html", gin.H{
				"title": "Error",
				"error": "Invalid user ID",
			})
			return
		}

		db := config.GetDB()
		user, err := models.FindUserByID(db, id)
		if err != nil {
			c.HTML(http.StatusNotFound, "error.html", gin.H{
				"title": "Error",
				"error": "User not found",
			})
			return
		}

		campgrounds, _ := models.FindCampgroundsByAuthor(db, id)

		c.HTML(http.StatusOK, "profile.html", gin.H{
			"title":       "User Profile",
			"user":        user,
			"campgrounds": campgrounds,
		})
	})

	// Authentication pages
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

	// Authentication actions
	r.POST("/register", authController.RegisterWeb)
	r.POST("/login", authController.LoginWeb)
	r.POST("/logout", func(c *gin.Context) {
		c.SetCookie("token", "", -1, "/", "", false, true)
		c.Redirect(http.StatusSeeOther, "/")
	})

	// Review actions (protected)
	r.POST("/campgrounds/:id/reviews", middleware.WebAuthRequired(), reviewController.CreateWeb)
	r.DELETE("/campgrounds/:id/reviews/:reviewId", middleware.WebAuthRequired(), reviewController.DeleteWeb)
}
