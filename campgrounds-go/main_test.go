package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	setupRoutes(r)
	return r
}

func TestHealthEndpoint(t *testing.T) {
	router := setupTestRouter()

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
}

func TestGetCampgrounds(t *testing.T) {
	router := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/campgrounds", nil)
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

func TestCreateCampground(t *testing.T) {
	router := setupTestRouter()

	newCampground := map[string]interface{}{
		"title":       "Test Campground",
		"description": "A test campground",
		"location":    "Test Location",
		"price":       29.99,
		"images":      []string{"test1.jpg", "test2.jpg"},
	}

	jsonData, _ := json.Marshal(newCampground)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/campgrounds", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != 201 {
		t.Errorf("Expected status code 201, got %d", w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse JSON response: %v", err)
	}

	campground, exists := response["campground"]
	if !exists {
		t.Error("Expected 'campground' field in response")
	}

	campgroundData, ok := campground.(map[string]interface{})
	if !ok {
		t.Error("Expected campground to be an object")
	}

	if campgroundData["title"] != "Test Campground" {
		t.Errorf("Expected title 'Test Campground', got %v", campgroundData["title"])
	}
}
