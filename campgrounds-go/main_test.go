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
		api.POST("/campgrounds", createCampground)
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