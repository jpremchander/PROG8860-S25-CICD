package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Campground struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title" binding:"required"`
	Price       float64            `bson:"price" json:"price" binding:"required"`
	Description string             `bson:"description" json:"description" binding:"required"`
	Location    string             `bson:"location" json:"location" binding:"required"`
	Images      []Image            `bson:"images" json:"images"`
	Author      primitive.ObjectID `bson:"author" json:"author"`
	Reviews     []primitive.ObjectID `bson:"reviews" json:"reviews"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

type Image struct {
	URL      string `bson:"url" json:"url"`
	Filename string `bson:"filename" json:"filename"`
}
