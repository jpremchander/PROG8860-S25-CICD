package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})
	
	api := r.Group("/api")
	{
		api.GET("/campgrounds", getCampgrounds)
		api.GET("/campgrounds/:id", getCampgroundByID)  // ✅ Add this route
		api.POST("/campgrounds", createCampground)
		api.PUT("/campgrounds/:id", updateCampground)
		api.DELETE("/campgrounds/:id", deleteCampground)
	}
	
	return r
}

func TestHealthEndpoint(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

func TestGetCampgrounds(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/campgrounds", nil)
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

func TestCreateCampground(t *testing.T) {
	router := setupRouter()
	
	campground := map[string]interface{}{
		"title":       "Test Camp",
		"description": "Test Description",
		"location":    "Test Location",
		"price":       25.99,
	}
	
	jsonData, _ := json.Marshal(campground)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/campgrounds", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != 201 {
		t.Errorf("Expected 201, got %d", w.Code)
	}
}

func TestUpdateCampground(t *testing.T) {
	router := setupRouter()
	
	// Test updating an existing campground (ID 1 exists in sample data)
	updateData := map[string]interface{}{
		"title":       "Updated Test Camp",
		"description": "Updated Description",
		"location":    "Updated Location",
		"price":       35.99,
	}
	
	jsonData, _ := json.Marshal(updateData)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/campgrounds/1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected 200, got %d", w.Code)
	}
	
	// Verify response contains updated data
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	
	if campground, ok := response["campground"].(map[string]interface{}); ok {
		if title, ok := campground["title"].(string); ok && title != "Updated Test Camp" {
			t.Errorf("Expected title to be 'Updated Test Camp', got %s", title)
		}
	}
}

func TestDeleteCampground(t *testing.T) {
	router := setupRouter()
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/campgrounds/1", nil)
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

func TestInvalidCampgroundID(t *testing.T) {
	router := setupRouter()
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/campgrounds/invalid", nil)
	router.ServeHTTP(w, req)

	// Now that we have the proper route, it should return 400 for invalid ID format
	if w.Code != 400 {
		t.Errorf("Expected 400 for invalid ID, got %d", w.Code)
	}
}
