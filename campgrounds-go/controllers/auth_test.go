package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"yelpcamp-go/config"
	"yelpcamp-go/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() {
	var err error
	config.DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test database")
	}
	
	// Auto-migrate test models
	config.DB.AutoMigrate(&models.User{}, &models.Campground{}, &models.Review{}, &models.Image{})
}

func TestAuthController_Register(t *testing.T) {
	// Setup
	setupTestDB()
	os.Setenv("JWT_SECRET", "test-secret")
	
	gin.SetMode(gin.TestMode)
	router := gin.New()
	authController := NewAuthController()
	router.POST("/register", authController.Register)

	// Test valid registration
	user := models.UserInput{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	
	jsonData, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusCreated, w.Code)
	
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "User registered successfully", response["message"])
	assert.NotEmpty(t, response["token"])
}

func TestAuthController_Login(t *testing.T) {
	// Setup
	setupTestDB()
	os.Setenv("JWT_SECRET", "test-secret")
	
	// Create a test user first
	user := models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	user.HashPassword()
	config.DB.Create(&user)
	
	gin.SetMode(gin.TestMode)
	router := gin.New()
	authController := NewAuthController()
	router.POST("/login", authController.Login)

	// Test valid login
	loginData := models.LoginInput{
		Username: "testuser",
		Password: "password123",
	}
	
	jsonData, _ := json.Marshal(loginData)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Login successful", response["message"])
	assert.NotEmpty(t, response["token"])
}
