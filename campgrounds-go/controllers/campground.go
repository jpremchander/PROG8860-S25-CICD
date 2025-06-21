package controllers

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"yelpcamp-go/config"
	"yelpcamp-go/models"
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

	c.JSON(http.StatusOK, gin.H{
		"message":     "Campgrounds retrieved successfully",
		"data":        campgrounds,
		"count":       len(campgrounds),
		"campgrounds": campgrounds, // For backward compatibility
	})
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

	c.JSON(http.StatusOK, gin.H{
		"message":    "Campground retrieved successfully",
		"campground": campground,
	})
}

// Helper function to save uploaded images to local storage
func (cc *CampgroundController) saveUploadedImage(fileHeader *multipart.FileHeader, userID string, index int) (*models.Image, error) {
	// Create uploads directory if it doesn't exist
	uploadDir := "static/uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %v", err)
	}

	// Generate unique filename
	timestamp := time.Now().Unix()
	ext := filepath.Ext(fileHeader.Filename)
	if ext == "" {
		ext = ".jpg" // Default extension
	}
	
	// Clean the extension and make it lowercase
	ext = strings.ToLower(ext)
	filename := fmt.Sprintf("%s_%d_%d%s", userID, timestamp, index, ext)
	filePath := filepath.Join(uploadDir, filename)

	// Open uploaded file
	src, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %v", err)
	}
	defer dst.Close()

	// Copy file content
	if _, err := io.Copy(dst, src); err != nil {
		return nil, fmt.Errorf("failed to save file: %v", err)
	}

	// Return image info with correct URL path
	return &models.Image{
		URL:      fmt.Sprintf("/static/uploads/%s", filename),
		Filename: fileHeader.Filename,
		Key:      filename,
	}, nil
}

// Web methods for form handling
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
		Images:      []models.Image{}, // Initialize empty images array
	}

	// Handle image uploads
	form, err := c.MultipartForm()
	if err == nil && form.File["images"] != nil {
		files := form.File["images"]
		fmt.Printf("📸 Processing %d uploaded images...\n", len(files))
		
		for i, file := range files {
			if file.Size > 0 { // Only process non-empty files
				fmt.Printf("📷 Processing image %d: %s (size: %d bytes)\n", i+1, file.Filename, file.Size)
				
				// Save the actual uploaded image
				savedImage, err := cc.saveUploadedImage(file, userID.Hex(), i)
				if err != nil {
					fmt.Printf("❌ Error saving image %s: %v\n", file.Filename, err)
					continue
				}
				
				fmt.Printf("✅ Successfully saved image: %s -> %s\n", file.Filename, savedImage.URL)
				campground.Images = append(campground.Images, *savedImage)
			}
		}
	}

	// If no images were uploaded or all failed, add a default one
	if len(campground.Images) == 0 {
		fmt.Println("📷 No images uploaded, using default placeholder")
		campground.Images = []models.Image{
			{
				URL:      "https://images.unsplash.com/photo-1504851149312-7a075b496cc7?w=800&h=600&fit=crop",
				Filename: "default-campground.jpg",
				Key:      fmt.Sprintf("default-%s", userID.Hex()),
			},
		}
	} else {
		fmt.Printf("✅ Successfully processed %d images for campground\n", len(campground.Images))
	}

	db := config.GetDB()
	if err := campground.Create(db); err != nil {
		c.HTML(http.StatusInternalServerError, "new.html", gin.H{
			"title": "Add New Campground",
			"error": "Could not create campground: " + err.Error(),
		})
		return
	}

	fmt.Printf("🎉 Campground created successfully with %d images\n", len(campground.Images))
	c.Redirect(http.StatusSeeOther, "/campgrounds/"+campground.ID.Hex())
}

func (cc *CampgroundController) Create(c *gin.Context) {
	var input models.CampgroundInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	campground := models.Campground{
		Title:       input.Title,
		Description: input.Description,
		Location:    input.Location,
		Price:       input.Price,
		AuthorID:    userID,
		Images:      []models.Image{}, // Initialize empty images array
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

	userID := c.MustGet("user_id").(primitive.ObjectID)

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

	userID := c.MustGet("user_id").(primitive.ObjectID)

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

func (cc *CampgroundController) UploadImages(c *gin.Context) {
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

	var uploadedImages []models.Image
	for i, file := range files {
		// Save uploaded image
		savedImage, err := cc.saveUploadedImage(file, userID.Hex(), i)
		if err != nil {
			fmt.Printf("Error saving image %s: %v\n", file.Filename, err)
			continue
		}
		uploadedImages = append(uploadedImages, *savedImage)
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
