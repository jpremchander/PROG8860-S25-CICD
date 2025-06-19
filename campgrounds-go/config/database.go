package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DB         *mongo.Database
	MongoClient *mongo.Client
)

func ConnectMongoDB() {
	// MongoDB connection string
	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=admin",
		os.Getenv("MONGO_USER"),
		os.Getenv("MONGO_PASSWORD"),
		os.Getenv("MONGO_HOST"),
		os.Getenv("MONGO_PORT"),
		os.Getenv("MONGO_DATABASE"),
	)

	// Set client options
	clientOptions := options.Client().ApplyURI(mongoURI)
	clientOptions.SetMaxPoolSize(20)
	clientOptions.SetMinPoolSize(5)
	clientOptions.SetMaxConnIdleTime(30 * time.Second)

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Test the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	MongoClient = client
	DB = client.Database(os.Getenv("MONGO_DATABASE"))
	
	log.Println("MongoDB connected successfully")
}

func GetDB() *mongo.Database {
	return DB
}

func DisconnectMongoDB() {
	if MongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		MongoClient.Disconnect(ctx)
	}
}
