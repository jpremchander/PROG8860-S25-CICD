package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Review struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Rating       int                `json:"rating" bson:"rating"`
	Body         string             `json:"body" bson:"body"`
	AuthorID     primitive.ObjectID `json:"author_id" bson:"author_id"`
	Author       *User              `json:"author,omitempty" bson:"author,omitempty"`
	CampgroundID primitive.ObjectID `json:"campground_id" bson:"campground_id"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}

type ReviewInput struct {
	Rating int    `json:"rating" binding:"required,min=1,max=5"`
	Body   string `json:"body" binding:"required,min=5"`
}

func (r *Review) Create(db *mongo.Database) error {
	r.CreatedAt = time.Now()
	r.UpdatedAt = time.Now()
	
	collection := db.Collection("reviews")
	result, err := collection.InsertOne(context.Background(), r)
	if err != nil {
		return err
	}
	
	r.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *Review) Delete(db *mongo.Database) error {
	collection := db.Collection("reviews")
	_, err := collection.DeleteOne(context.Background(), bson.M{"_id": r.ID})
	return err
}

func FindReviewByID(db *mongo.Database, id primitive.ObjectID) (*Review, error) {
	var review Review
	collection := db.Collection("reviews")
	
	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&review)
	if err != nil {
		return nil, err
	}
	
	return &review, nil
}

func FindReviewsByCampground(db *mongo.Database, campgroundID primitive.ObjectID) ([]Review, error) {
	var reviews []Review
	collection := db.Collection("reviews")
	
	// Aggregation pipeline to populate author
	pipeline := []bson.M{
		{"$match": bson.M{"campground_id": campgroundID}},
		{
			"$lookup": bson.M{
				"from":         "users",
				"localField":   "author_id",
				"foreignField": "_id",
				"as":           "author",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$author",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$sort": bson.M{"created_at": -1},
		},
	}
	
	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	
	if err = cursor.All(context.Background(), &reviews); err != nil {
		return nil, err
	}
	
	return reviews, nil
}
