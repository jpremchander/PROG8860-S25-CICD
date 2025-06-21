package controllers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"campgrounds-go/config"
	"campgrounds-go/models"
)

func CreateReview(c *gin.Context) {
	campgroundID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(campgroundID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campground ID"})
		return
	}

	// Get form data
	body := c.PostForm("body")
	ratingStr := c.PostForm("rating")
	rating, err := strconv.Atoi(ratingStr)
	if err != nil || rating < 1 || rating > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rating"})
		return
	}

	// Get user from session
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	review := models.Review{
		Body:       body,
		Rating:     rating,
		Author:     userID,
		Campground: objectID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	collection := config.DB.Collection("reviews")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, review)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create review"})
		return
	}

	// Add review to campground
	campgroundCollection := config.DB.Collection("campgrounds")
	_, err = campgroundCollection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$push": bson.M{"reviews": result.InsertedID}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update campground"})
		return
	}

	c.Redirect(http.StatusSeeOther, "/campgrounds/"+campgroundID)
}

func DeleteReview(c *gin.Context) {
	campgroundID := c.Param("id")
	reviewID := c.Param("reviewId")

	campgroundObjectID, err := primitive.ObjectIDFromHex(campgroundID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campground ID"})
		return
	}

	reviewObjectID, err := primitive.ObjectIDFromHex(reviewID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	collection := config.DB.Collection("reviews")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Delete review
	_, err = collection.DeleteOne(ctx, bson.M{"_id": reviewObjectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete review"})
		return
	}

	// Remove review from campground
	campgroundCollection := config.DB.Collection("campgrounds")
	_, err = campgroundCollection.UpdateOne(
		ctx,
		bson.M{"_id": campgroundObjectID},
		bson.M{"$pull": bson.M{"reviews": reviewObjectID}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update campground"})
		return
	}

	c.Redirect(http.StatusSeeOther, "/campgrounds/"+campgroundID)
}
