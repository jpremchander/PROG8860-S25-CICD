package models

import (
	"context"
	"log"
	"yelpcamp-go/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AutoMigrate() {
	db := config.GetDB()
	
	// Create indexes for MongoDB collections
	log.Println("Creating MongoDB indexes...")
	
	// Users collection indexes
	usersCollection := db.Collection("users")
	
	// Username unique index
	usernameIndexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	usersCollection.Indexes().CreateOne(context.Background(), usernameIndexModel)
	
	// Email unique index
	emailIndexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	usersCollection.Indexes().CreateOne(context.Background(), emailIndexModel)
	
	// Campgrounds collection indexes
	campgroundsCollection := db.Collection("campgrounds")
	
	// Author ID index
	authorIndexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "author_id", Value: 1}},
	}
	campgroundsCollection.Indexes().CreateOne(context.Background(), authorIndexModel)
	
	// Created at index
	createdAtIndexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "created_at", Value: -1}},
	}
	campgroundsCollection.Indexes().CreateOne(context.Background(), createdAtIndexModel)
	
	// Reviews collection indexes
	reviewsCollection := db.Collection("reviews")
	
	// Campground ID index
	campgroundIdIndexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "campground_id", Value: 1}},
	}
	reviewsCollection.Indexes().CreateOne(context.Background(), campgroundIdIndexModel)
	
	// Review author ID index
	reviewAuthorIndexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "author_id", Value: 1}},
	}
	reviewsCollection.Indexes().CreateOne(context.Background(), reviewAuthorIndexModel)
	
	log.Println("✅ MongoDB indexes created successfully")
}
