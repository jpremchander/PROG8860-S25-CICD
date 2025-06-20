package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"yelpcamp-go/config"
	"yelpcamp-go/models"
)

type ReviewController struct{}

func NewReviewController() *ReviewController {
	return &ReviewController{}
}

func (rc *ReviewController) Create(c *gin.Context) {
	campgroundIDParam := c.Param("id")
	campgroundID, err := primitive.ObjectIDFromHex(campgroundIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campground ID"})
		return
	}

	var input models.ReviewInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	// Check if campground exists
	db := config.GetDB()
	_, err = models.FindCampgroundByID(db, campgroundID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campground not found"})
		return
	}

	review := models.Review{
		Rating:       input.Rating,
		Body:         input.Body,
		AuthorID:     userID,
		CampgroundID: campgroundID,
	}

	if err := review.Create(db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create review"})
		return
	}

	// Load the author information
	user, _ := models.FindUserByID(db, userID)
	review.Author = user

	c.JSON(http.StatusCreated, gin.H{
		"message": "Review created successfully",
		"review":  review,
	})
}

func (rc *ReviewController) CreateWeb(c *gin.Context) {
	campgroundIDParam := c.Param("id")
	campgroundID, err := primitive.ObjectIDFromHex(campgroundIDParam)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"title": "Error",
			"error": "Invalid campground ID",
		})
		return
	}

	ratingStr := c.PostForm("rating")
	body := c.PostForm("body")

	if ratingStr == "" || body == "" {
		c.Redirect(http.StatusSeeOther, "/campgrounds/"+campgroundIDParam+"?error=All fields are required")
		return
	}

	rating, err := strconv.Atoi(ratingStr)
	if err != nil || rating < 1 || rating > 5 {
		c.Redirect(http.StatusSeeOther, "/campgrounds/"+campgroundIDParam+"?error=Invalid rating")
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	// Check if campground exists
	db := config.GetDB()
	_, err = models.FindCampgroundByID(db, campgroundID)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"title": "Error",
			"error": "Campground not found",
		})
		return
	}

	review := models.Review{
		Rating:       rating,
		Body:         body,
		AuthorID:     userID,
		CampgroundID: campgroundID,
	}

	if err := review.Create(db); err != nil {
		c.Redirect(http.StatusSeeOther, "/campgrounds/"+campgroundIDParam+"?error=Could not create review")
		return
	}

	c.Redirect(http.StatusSeeOther, "/campgrounds/"+campgroundIDParam)
}

func (rc *ReviewController) Delete(c *gin.Context) {
	campgroundIDParam := c.Param("id")
	campgroundID, err := primitive.ObjectIDFromHex(campgroundIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campground ID"})
		return
	}

	reviewIDParam := c.Param("reviewId")
	reviewID, err := primitive.ObjectIDFromHex(reviewIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	db := config.GetDB()
	review, err := models.FindReviewByID(db, reviewID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	if review.CampgroundID != campgroundID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Review does not belong to this campground"})
		return
	}

	if review.AuthorID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to delete this review"})
		return
	}

	if err := review.Delete(db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete review"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review deleted successfully"})
}

func (rc *ReviewController) DeleteWeb(c *gin.Context) {
	campgroundIDParam := c.Param("id")
	reviewIDParam := c.Param("reviewId")
	
	campgroundID, err := primitive.ObjectIDFromHex(campgroundIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campground ID"})
		return
	}

	reviewID, err := primitive.ObjectIDFromHex(reviewIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	db := config.GetDB()
	review, err := models.FindReviewByID(db, reviewID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	if review.AuthorID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	if err := review.Delete(db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete review"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review deleted successfully"})
}
