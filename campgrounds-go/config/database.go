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
	DB          *mongo.Database
	MongoClient *mongo.Client
)

func ConnectMongoDB() {
	mongoHost := getEnvOrDefault("MONGO_HOST", "localhost")
	mongoPort := getEnvOrDefault("MONGO_PORT", "27017")
	mongoDatabase := getEnvOrDefault("MONGO_DATABASE", "yelpcamp_dev")

	// Simple connection string without authentication for Docker
	connectionString := fmt.Sprintf("mongodb://%s:%s/%s", mongoHost, mongoPort, mongoDatabase)
	
	log.Printf("🔗 Connecting to MongoDB: %s", connectionString)

	clientOptions := options.Client().ApplyURI(connectionString)
	clientOptions.SetMaxPoolSize(20)
	clientOptions.SetMinPoolSize(5)
	clientOptions.SetMaxConnIdleTime(30 * time.Second)
	clientOptions.SetServerSelectionTimeout(10 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Printf("❌ MongoDB connection failed: %v", err)
		log.Fatal("Could not connect to MongoDB")
	}

	// Test the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("❌ MongoDB ping failed: %v", err)
		log.Fatal("Could not ping MongoDB")
	}

	MongoClient = client
	DB = client.Database(mongoDatabase)
	
	log.Printf("✅ MongoDB connected successfully to database: %s", mongoDatabase)
}

func GetDB() *mongo.Database {
	return DB
}

func DisconnectMongoDB() {
	if MongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		MongoClient.Disconnect(ctx)
		log.Println("🔌 MongoDB disconnected")
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
