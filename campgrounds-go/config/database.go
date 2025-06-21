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
	clientOptions.SetServerSelectionTimeout(15 * time.Second)
	clientOptions.SetConnectTimeout(15 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Printf("❌ MongoDB connection failed: %v", err)
		log.Fatal("Could not connect to MongoDB")
	}

	// Test the connection with retries
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
		err = client.Ping(pingCtx, nil)
		pingCancel()
		
		if err == nil {
			break
		}
		
		log.Printf("❌ MongoDB ping attempt %d failed: %v", i+1, err)
		if i < maxRetries-1 {
			log.Printf("⏳ Retrying in 2 seconds...")
			time.Sleep(2 * time.Second)
		}
	}

	if err != nil {
		log.Printf("❌ MongoDB ping failed after %d attempts: %v", maxRetries, err)
		log.Fatal("Could not ping MongoDB")
	}

	MongoClient = client
	DB = client.Database(mongoDatabase)
	
	log.Printf("✅ MongoDB connected successfully to database: %s", mongoDatabase)
	
	// Test a simple operation
	testCtx, testCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer testCancel()
	
	count, err := DB.Collection("campgrounds").CountDocuments(testCtx, map[string]interface{}{})
	if err != nil {
		log.Printf("⚠️ Warning: Could not test database operation: %v", err)
	} else {
		log.Printf("📊 Current campgrounds in database: %d", count)
	}
}

func GetDB() *mongo.Database {
	if DB == nil {
		log.Println("⚠️ Warning: Database connection is nil")
	}
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
