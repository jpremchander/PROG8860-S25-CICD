package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCampgroundModel(t *testing.T) {
	campground := Campground{
		Title:       "Test Campground",
		Price:       25.99,
		Description: "A beautiful test campground",
		Location:    "Test Location",
		Images:      []Image{},
		Author:      primitive.NewObjectID(),
		Reviews:     []primitive.ObjectID{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	assert.Equal(t, "Test Campground", campground.Title)
	assert.Equal(t, 25.99, campground.Price)
	assert.Equal(t, "A beautiful test campground", campground.Description)
	assert.Equal(t, "Test Location", campground.Location)
	assert.NotNil(t, campground.Author)
	assert.NotNil(t, campground.CreatedAt)
	assert.NotNil(t, campground.UpdatedAt)
}

func TestImageModel(t *testing.T) {
	image := Image{
		URL:      "https://example.com/image.jpg",
		Filename: "image.jpg",
	}

	assert.Equal(t, "https://example.com/image.jpg", image.URL)
	assert.Equal(t, "image.jpg", image.Filename)
}
