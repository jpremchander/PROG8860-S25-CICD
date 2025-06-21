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
	"campgrounds-go/utils"
)

func GetCampgrounds(c *gin.Context) {
	collection := config.DB.Collection("campgrounds")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch campgrounds"})
		return
	}
	defer cursor.Close(ctx)

	var campgrounds []models.Campground
	if err = cursor.All(ctx, &campgrounds); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode campgrounds"})
		return
	}

	c.HTML(http.StatusOK, "campgrounds/index.html", gin.H{
		"campgrounds": campgrounds,
		"title":       "All Campgrounds",
	})
}

func ShowCampground(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campground ID"})
		return
	}

	collection := config.DB.Collection("campgrounds")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var campground models.Campground
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&campground)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campground not found"})
		return
	}

	c.HTML(http.StatusOK, "campgrounds/show.html", gin.H{
		"campground": campground,
		"title":      campground.Title,
	})
}

func NewCampground(c *gin.Context) {
	c.HTML(http.StatusOK, "campgrounds/new.html", gin.H{
		"title": "New Campground",
	})
}

func CreateCampground(c *gin.Context) {
	var campground models.Campground
	if err := c.ShouldBind(&campground); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Handle file uploads
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
		return
	}

	files := form.File["images"]
	for _, file := range files {
		// Upload to S3
		imageURL, filename, err := utils.UploadToS3(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image"})
			return
		}

		campground.Images = append(campground.Images, models.Image{
			URL:      imageURL,
			Filename: filename,
		})
	}

	// Get user ID from session (implement session middleware)
	userID, _ := c.Get("userID")
	if userID != nil {
		campground.Author = userID.(primitive.ObjectID)
	}

	campground.CreatedAt = time.Now()
	campground.UpdatedAt = time.Now()

	collection := config.DB.Collection("campgrounds")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, campground)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create campground"})
		return
	}

	c.Redirect(http.StatusSeeOther, "/campgrounds/"+result.InsertedID.(primitive.ObjectID).Hex())
}

func EditCampground(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campground ID"})
		return
	}

	collection := config.DB.Collection("campgrounds")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var campground models.Campground
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&campground)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campground not found"})
		return
	}

	c.HTML(http.StatusOK, "campgrounds/edit.html", gin.H{
		"campground": campground,
		"title":      "Edit " + campground.Title,
	})
}

func UpdateCampground(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campground ID"})
		return
	}

	var updateData models.Campground
	if err := c.ShouldBind(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateData.UpdatedAt = time.Now()

	collection := config.DB.Collection("campgrounds")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"title":       updateData.Title,
			"price":       updateData.Price,
			"description": updateData.Description,
			"location":    updateData.Location,
			"updated_at":  updateData.UpdatedAt,
		},
	}

	_, err = collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update campground"})
		return
	}

	c.Redirect(http.StatusSeeOther, "/campgrounds/"+id)
}

func DeleteCampground(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campground ID"})
		return
	}

	collection := config.DB.Collection("campgrounds")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete campground"})
		return
	}

	c.Redirect(http.StatusSeeOther, "/campgrounds")
}
