package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"campgrounds-app/config"
	"campgrounds-app/models"
)

type CampgroundController struct{}

func NewCampgroundController() *CampgroundController {
	return &CampgroundController{}
}

func (cc *CampgroundController) GetAll(c *gin.Context) {
	db := config.GetDB()
	campgrounds, err := models.FindAllCampgrounds(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch campgrounds"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"campgrounds": campgrounds})
}

func (cc *CampgroundController) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campground ID"})
		return
	}

	db := config.GetDB()
	campground, err := models.FindCampgroundByID(db, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Campground not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch campground"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"campground": campground})
}

func (cc *CampgroundController) Create(c *gin.Context) {
	var input models.CampgroundInput
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	campground := models.Campground{
		Title:       input.Title,
		Description: input.Description,
		Location:    input.Location,
		Price:       input.Price,
		AuthorID:    userID,
	}

	db := config.GetDB()
	if err := campground.Create(db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create campground"})
		return
	}

	// Fetch the created campground with author info
	createdCampground, _ := models.FindCampgroundByID(db, campground.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Campground created successfully",
		"campground": createdCampground,
	})
}

func (cc *CampgroundController) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campground ID"})
		return
	}

	var input models.CampgroundInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	db := config.GetDB()
	campground, err := models.FindCampgroundByID(db, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Campground not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch campground"})
		}
		return
	}

	if campground.AuthorID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to update this campground"})
		return
	}

	campground.Title = input.Title
	campground.Description = input.Description
	campground.Location = input.Location
	campground.Price = input.Price

	if err := campground.Update(db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update campground"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Campground updated successfully",
		"campground": campground,
	})
}

func (cc *CampgroundController) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campground ID"})
		return
	}

	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	db := config.GetDB()
	campground, err := models.FindCampgroundByID(db, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Campground not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch campground"})
		}
		return
	}

	if campground.AuthorID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to delete this campground"})
		return
	}

	if err := campground.Delete(db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete campground"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Campground deleted successfully"})
}

func (cc *CampgroundController) UploadImages(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campground ID"})
		return
	}

	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	db := config.GetDB()
	campground, err := models.FindCampgroundByID(db, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campground not found"})
		return
	}

	if campground.AuthorID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to upload images for this campground"})
		return
	}

	// Handle multipart form
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not parse form"})
		return
	}

	files := form.File["images"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No images provided"})
		return
	}

	// Upload images to Cloudinary
	cloudinaryService, err := NewCloudinaryService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not initialize image service"})
		return
	}

	var uploadedImages []models.Image
	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			continue
		}
		defer src.Close()

		result, err := cloudinaryService.UploadImage(src, "yelpcamp")
		if err != nil {
			continue
		}

		image := models.Image{
			URL:      result.SecureURL,
			Filename: file.Filename,
			PublicID: result.PublicID,
		}
		uploadedImages = append(uploadedImages, image)
	}

	if len(uploadedImages) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not upload any images"})
		return
	}

	// Add images to campground
	campground.Images = append(campground.Images, uploadedImages...)
	if err := campground.Update(db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update campground with images"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Images uploaded successfully",
		"images":  uploadedImages,
	})
}

// Add these methods to the existing CampgroundController

func (cc *CampgroundController) CreateWeb(c *gin.Context) {
	title := c.PostForm("title")
	description := c.PostForm("description")
	location := c.PostForm("location")
	priceStr := c.PostForm("price")

	if title == "" || description == "" || location == "" || priceStr == "" {
		c.HTML(http.StatusBadRequest, "new.html", gin.H{
			"title": "Add New Campground",
			"error": "All fields are required",
		})
		return
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil || price < 0 {
		c.HTML(http.StatusBadRequest, "new.html", gin.H{
			"title": "Add New Campground",
			"error": "Invalid price",
		})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	campground := models.Campground{
		Title:       title,
		Description: description,
		Location:    location,
		Price:       price,
		AuthorID:    userID,
	}

	db := config.GetDB()
	if err := campground.Create(db); err != nil {
		c.HTML(http.StatusInternalServerError, "new.html", gin.H{
			"title": "Add New Campground",
			"error": "Could not create campground",
		})
		return
	}

	c.Redirect(http.StatusSeeOther, "/campgrounds/"+campground.ID.Hex())
}

func (cc *CampgroundController) UpdateWeb(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
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
	if err != nil || price < 0 {
		c.HTML(http.StatusBadRequest, "edit.html", gin.H{
			"title": "Edit Campground",
			"error": "Invalid price",
		})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)
	db := config.GetDB()
	
	campground, err := models.FindCampgroundByID(db, id)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"title": "Error",
			"error": "Campground not found",
		})
		return
	}

	if campground.AuthorID != userID {
		c.HTML(http.StatusForbidden, "error.html", gin.H{
			"title": "Error",
			"error": "Not authorized to update this campground",
		})
		return
	}

	campground.Title = title
	campground.Description = description
	campground.Location = location
	campground.Price = price

	if err := campground.Update(db); err != nil {
		c.HTML(http.StatusInternalServerError, "edit.html", gin.H{
			"title": "Edit Campground",
			"error": "Could not update campground",
		})
		return
	}

	c.Redirect(http.StatusSeeOther, "/campgrounds/"+campground.ID.Hex())
}

func (cc *CampgroundController) DeleteWeb(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campground ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)
	db := config.GetDB()
	
	campground, err := models.FindCampgroundByID(db, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campground not found"})
		return
	}

	if campground.AuthorID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	if err := campground.Delete(db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete campground"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Campground deleted successfully"})
}
