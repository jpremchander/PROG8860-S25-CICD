package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to YelpCamp Go!",
			"status":  "running",
		})
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "campgrounds-go",
		})
	})

	r.GET("/campgrounds", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"campgrounds": []map[string]interface{}{
				{
					"id":          1,
					"title":       "Sample Campground",
					"description": "A beautiful campground for testing",
					"price":       25.99,
					"location":    "Test Location",
				},
			},
		})
	})

	return r
}

func TestHealthEndpoint(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse JSON response: %v", err)
	}

	if response["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got %v", response["status"])
	}

	if response["service"] != "campgrounds-go" {
		t.Errorf("Expected service 'campgrounds-go', got %v", response["service"])
	}
}

func TestHomeRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse JSON response: %v", err)
	}

	if response["message"] != "Welcome to YelpCamp Go!" {
		t.Errorf("Expected welcome message, got %v", response["message"])
	}
}

func TestCampgroundsRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/campgrounds", nil)
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse JSON response: %v", err)
	}

	campgrounds, exists := response["campgrounds"]
	if !exists {
		t.Error("Expected 'campgrounds' field in response")
	}

	campgroundsList, ok := campgrounds.([]interface{})
	if !ok || len(campgroundsList) == 0 {
		t.Error("Expected non-empty campgrounds array")
	}
}

func TestInvalidRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/invalid", nil)
	router.ServeHTTP(w, req)

	if w.Code != 404 {
		t.Errorf("Expected status code 404, got %d", w.Code)
	}
}
