package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Image struct {
	URL      string `json:"url" bson:"url"`
	Filename string `json:"filename" bson:"filename"`
	Key      string `json:"key" bson:"key"`
}

type Campground struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	Location    string             `json:"location" bson:"location"`
	Price       float64            `json:"price" bson:"price"`
	Images      []Image            `json:"images" bson:"images"`
	AuthorID    primitive.ObjectID `json:"author_id" bson:"author_id"`
	Author      *User              `json:"author,omitempty" bson:"author,omitempty"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

type CampgroundInput struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Location    string  `json:"location" binding:"required"`
	Price       float64 `json:"price" binding:"required,min=0"`
}

func (c *Campground) Create(db *mongo.Database) error {
	c.ID = primitive.NewObjectID()
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()

	collection := db.Collection("campgrounds")
	_, err := collection.InsertOne(context.Background(), c)
	return err
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
	filter := bson.M{"_id": c.ID}

	_, err := collection.DeleteOne(context.Background(), filter)
	return err
}

func FindAllCampgrounds(db *mongo.Database) ([]Campground, error) {
	collection := db.Collection("campgrounds")
	
	// Create aggregation pipeline to join with users collection
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

	var campgrounds []Campground
	if err = cursor.All(context.Background(), &campgrounds); err != nil {
		return nil, err
	}

	return campgrounds, nil
}

func FindCampgroundByID(db *mongo.Database, id primitive.ObjectID) (*Campground, error) {
	collection := db.Collection("campgrounds")
	
	// Create aggregation pipeline to join with users collection
	pipeline := []bson.M{
		{
			"$match": bson.M{"_id": id},
		},
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
	}

	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var campgrounds []Campground
	if err = cursor.All(context.Background(), &campgrounds); err != nil {
		return nil, err
	}

	if len(campgrounds) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return &campgrounds[0], nil
}

func FindCampgroundsByAuthor(db *mongo.Database, authorID primitive.ObjectID) ([]Campground, error) {
	collection := db.Collection("campgrounds")
	
	pipeline := []bson.M{
		{
			"$match": bson.M{"author_id": authorID},
		},
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

	var campgrounds []Campground
	if err = cursor.All(context.Background(), &campgrounds); err != nil {
		return nil, err
	}

	return campgrounds, nil
}
