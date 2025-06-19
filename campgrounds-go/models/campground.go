package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Campground struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	Location    string             `json:"location" bson:"location"`
	Price       float64            `json:"price" bson:"price"`
	Images      []Image            `json:"images" bson:"images"`
	AuthorID    primitive.ObjectID `json:"author_id" bson:"author_id"`
	Author      *User              `json:"author,omitempty" bson:"author,omitempty"`
	Reviews     []Review           `json:"reviews,omitempty" bson:"reviews,omitempty"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

type CampgroundInput struct {
	Title       string  `json:"title" binding:"required,min=3,max=100"`
	Description string  `json:"description" binding:"required,min=10"`
	Location    string  `json:"location" binding:"required"`
	Price       float64 `json:"price" binding:"required,min=0"`
}

type Image struct {
	URL      string `json:"url" bson:"url"`
	Filename string `json:"filename" bson:"filename"`
	PublicID string `json:"public_id" bson:"public_id"`
}

func (c *Campground) Create(db *mongo.Database) error {
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	c.Images = []Image{} // Initialize empty images array
	
	collection := db.Collection("campgrounds")
	result, err := collection.InsertOne(context.Background(), c)
	if err != nil {
		return err
	}
	
	c.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (c *Campground) Update(db *mongo.Database) error {
	c.UpdatedAt = time.Now()
	
	collection := db.Collection("campgrounds")
	filter := bson.M{"_id": c.ID}
	update := bson.M{"$set": c}
	
	_, err := collection.UpdateOne(context.Background(), filter, update)
	return err
}

func (c *Campground) Delete(db *mongo.Database) error {
	collection := db.Collection("campgrounds")
	_, err := collection.DeleteOne(context.Background(), bson.M{"_id": c.ID})
	return err
}

func FindAllCampgrounds(db *mongo.Database) ([]Campground, error) {
	var campgrounds []Campground
	collection := db.Collection("campgrounds")
	
	// Create aggregation pipeline to populate author
	pipeline := []bson.M{
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
	
	if err = cursor.All(context.Background(), &campgrounds); err != nil {
		return nil, err
	}
	
	return campgrounds, nil
}

func FindCampgroundByID(db *mongo.Database, id primitive.ObjectID) (*Campground, error) {
	var campground Campground
	collection := db.Collection("campgrounds")
	
	// Aggregation pipeline to populate author and reviews
	pipeline := []bson.M{
		{"$match": bson.M{"_id": id}},
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
			"$lookup": bson.M{
				"from":         "reviews",
				"localField":   "_id",
				"foreignField": "campground_id",
				"as":           "reviews",
			},
		},
	}
	
	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	
	if cursor.Next(context.Background()) {
		if err := cursor.Decode(&campground); err != nil {
			return nil, err
		}
		return &campground, nil
	}
	
	return nil, mongo.ErrNoDocuments
}

func FindCampgroundsByAuthor(db *mongo.Database, authorID primitive.ObjectID) ([]Campground, error) {
	var campgrounds []Campground
	collection := db.Collection("campgrounds")
	
	cursor, err := collection.Find(context.Background(), bson.M{"author_id": authorID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	
	if err = cursor.All(context.Background(), &campgrounds); err != nil {
		return nil, err
	}
	
	return campgrounds, nil
}
