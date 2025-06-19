package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"campgrounds-app/config"
	"campgrounds-app/models"
)

type ReviewController struct{}

func NewReviewController() *ReviewController {
	return &ReviewController{}
}

func (rc *ReviewController) Create(c *gin.Context) {
	campgroundID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campground ID"})
		return
	}

	var input models.ReviewInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Check if campground exists
	var campground models.Campground
	db := config.GetDB()
	if err := db.First(&campground, campgroundID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campground not found"})
		return
	}

	review := models.Review{
		Rating:       input.Rating,
		Body:         input.Body,
		AuthorID:     userID.(uint),
		CampgroundID: uint(campgroundID),
	}

	if err := db.Create(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create review"})
		return
	}

	// Load the author information
	db.Preload("Author").First(&review, review.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Review created successfully",
		"review":  review,
	})
}

func (rc *ReviewController) Delete(c *gin.Context) {
	campgroundID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campground ID"})
		return
	}

	reviewID, err := strconv.ParseUint(c.Param("reviewId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var review models.Review
	db := config.GetDB()
	
	if err := db.Where("id = ? AND campground_id = ?", reviewID, campgroundID).First(&review).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	if review.AuthorID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to delete this review"})
		return
	}

	if err := db.Delete(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete review"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review deleted successfully"})
}
