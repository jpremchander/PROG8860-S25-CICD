package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Review struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Body         string             `bson:"body" json:"body" binding:"required"`
	Rating       int                `bson:"rating" json:"rating" binding:"required,min=1,max=5"`
	Author       primitive.ObjectID `bson:"author" json:"author"`
	Campground   primitive.ObjectID `bson:"campground" json:"campground"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}
