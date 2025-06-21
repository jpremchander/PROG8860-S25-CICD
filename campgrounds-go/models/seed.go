package models

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"yelpcamp-go/config"
)

// ForceSeedData clears existing data and creates fresh sample data
func ForceSeedData() {
	db := config.GetDB()
	if db == nil {
		log.Println("❌ Database connection is nil, cannot force seed")
		return
	}
	
	log.Println("🧹 Force seeding: Clearing existing data...")
	
	// Clear existing data
	ctx := context.Background()
	_, err1 := db.Collection("reviews").DeleteMany(ctx, map[string]interface{}{})
	_, err2 := db.Collection("campgrounds").DeleteMany(ctx, map[string]interface{}{})
	_, err3 := db.Collection("users").DeleteMany(ctx, map[string]interface{}{})
	
	if err1 != nil || err2 != nil || err3 != nil {
		log.Printf("⚠️ Some errors during cleanup: reviews=%v, campgrounds=%v, users=%v", err1, err2, err3)
	}
	
	log.Println("✅ Cleared existing data")
	
	// Seed fresh data
	SeedData()
}

// SeedData creates sample users, campgrounds, and reviews for development
func SeedData() {
	db := config.GetDB()
	if db == nil {
		log.Println("❌ Database connection is nil, cannot seed data")
		return
	}

	log.Println("🌱 Starting sample data seeding...")
	
	// Create sample users with proper password hashing
	users := []User{
		{
			Username: "igoswamik",
			Email:    "igor@yelpcamp.com",
			Password: "password123",
		},
		{
			Username: "hannah",
			Email:    "hannah@yelpcamp.com",
			Password: "password123",
		},
		{
			Username: "bob",
			Email:    "bob@yelpcamp.com",
			Password: "password123",
		},
		{
			Username: "alice",
			Email:    "alice@yelpcamp.com",
			Password: "password123",
		},
	}
	
	// Insert users
	var userIDs []primitive.ObjectID
	log.Println("👥 Creating sample users...")
	for i := range users {
		log.Printf("Creating user: %s", users[i].Username)
		if err := users[i].Create(db); err != nil {
			log.Printf("❌ Error creating user %s: %v", users[i].Username, err)
			continue
		}
		userIDs = append(userIDs, users[i].ID)
		log.Printf("✅ Created user: %s (ID: %s)", users[i].Username, users[i].ID.Hex())
	}
	
	if len(userIDs) == 0 {
		log.Println("❌ No users created, skipping campgrounds")
		return
	}
	
	log.Printf("✅ Successfully created %d users", len(userIDs))
	
	// Create sample campgrounds with realistic data
	campgrounds := []Campground{
		{
			Title:       "Redwood National Park",
			Description: "Experience the majesty of the world's tallest trees in this pristine wilderness setting. Perfect for hiking, photography, and connecting with nature.",
			Location:    "Crescent City, California",
			Price:       25.00,
			AuthorID:    userIDs[0], // igoswamik
			Images: []Image{
				{
					URL:      "https://images.unsplash.com/photo-1441974231531-c6227db76b6e?w=800&h=600&fit=crop",
					Filename: "redwood-1.jpg",
					Key:      "yelpcamp/redwood-1",
				},
			},
			CreatedAt: time.Now().AddDate(0, 0, -10),
			UpdatedAt: time.Now().AddDate(0, 0, -10),
		},
		{
			Title:       "Grand Canyon South Rim",
			Description: "Wake up to breathtaking views of one of the world's natural wonders. This campground offers unparalleled sunrise and sunset viewing opportunities.",
			Location:    "Grand Canyon, Arizona",
			Price:       35.00,
			AuthorID:    userIDs[1], // hannah
			Images: []Image{
				{
					URL:      "https://images.unsplash.com/photo-1506905925346-21bda4d32df4?w=800&h=600&fit=crop",
					Filename: "grand-canyon-1.jpg",
					Key:      "yelpcamp/grand-canyon-1",
				},
			},
			CreatedAt: time.Now().AddDate(0, 0, -8),
			UpdatedAt: time.Now().AddDate(0, 0, -8),
		},
		{
			Title:       "Yellowstone Lake Lodge",
			Description: "Camp beside the pristine waters of Yellowstone Lake with opportunities for fishing, boating, and wildlife viewing. Hot springs nearby!",
			Location:    "Yellowstone National Park, Wyoming",
			Price:       40.00,
			AuthorID:    userIDs[2], // bob
			Images: []Image{
				{
					URL:      "https://images.unsplash.com/photo-1504851149312-7a075b496cc7?w=800&h=600&fit=crop",
					Filename: "yellowstone-1.jpg",
					Key:      "yelpcamp/yellowstone-1",
				},
			},
			CreatedAt: time.Now().AddDate(0, 0, -6),
			UpdatedAt: time.Now().AddDate(0, 0, -6),
		},
		{
			Title:       "Yosemite Valley Floor",
			Description: "Camp in the heart of Yosemite Valley with iconic views of El Capitan and Half Dome. Rock climbing and hiking trails accessible from camp.",
			Location:    "Yosemite National Park, California",
			Price:       45.00,
			AuthorID:    userIDs[3], // alice
			Images: []Image{
				{
					URL:      "https://images.unsplash.com/photo-1504280390367-361c6d9f38f4?w=800&h=600&fit=crop",
					Filename: "yosemite-1.jpg",
					Key:      "yelpcamp/yosemite-1",
				},
			},
			CreatedAt: time.Now().AddDate(0, 0, -4),
			UpdatedAt: time.Now().AddDate(0, 0, -4),
		},
		{
			Title:       "Glacier Point Overlook",
			Description: "Spectacular mountain camping with panoramic views of snow-capped peaks and alpine lakes. Perfect for stargazing and photography.",
			Location:    "Glacier National Park, Montana",
			Price:       30.00,
			AuthorID:    userIDs[0], // igoswamik
			Images: []Image{
				{
					URL:      "https://images.unsplash.com/photo-1445308394109-4ec2920981b1?w=800&h=600&fit=crop",
					Filename: "glacier-1.jpg",
					Key:      "yelpcamp/glacier-1",
				},
			},
			CreatedAt: time.Now().AddDate(0, 0, -2),
			UpdatedAt: time.Now().AddDate(0, 0, -2),
		},
		{
			Title:       "Zion Canyon Riverside",
			Description: "Camp along the Virgin River with towering red rock formations surrounding you. Easy access to hiking trails and the famous Narrows.",
			Location:    "Zion National Park, Utah",
			Price:       32.00,
			AuthorID:    userIDs[1], // hannah
			Images: []Image{
				{
					URL:      "https://images.unsplash.com/photo-1487730116645-74489c95b41b?w=800&h=600&fit=crop",
					Filename: "zion-1.jpg",
					Key:      "yelpcamp/zion-1",
				},
			},
			CreatedAt: time.Now().AddDate(0, 0, -1),
			UpdatedAt: time.Now().AddDate(0, 0, -1),
		},
	}
	
	// Insert campgrounds
	var campgroundIDs []primitive.ObjectID
	log.Println("🏕️ Creating sample campgrounds...")
	for i := range campgrounds {
		log.Printf("Creating campground: %s", campgrounds[i].Title)
		
		// Use direct MongoDB insertion to ensure it works
		collection := db.Collection("campgrounds")
		campgrounds[i].ID = primitive.NewObjectID()
		
		_, err := collection.InsertOne(context.Background(), campgrounds[i])
		if err != nil {
			log.Printf("❌ Error creating campground %s: %v", campgrounds[i].Title, err)
			continue
		}
		
		campgroundIDs = append(campgroundIDs, campgrounds[i].ID)
		log.Printf("✅ Created campground: %s (ID: %s)", campgrounds[i].Title, campgrounds[i].ID.Hex())
	}
	
	log.Printf("✅ Successfully created %d campgrounds", len(campgroundIDs))
	
	// Create sample reviews
	if len(campgroundIDs) > 0 && len(userIDs) > 0 {
		log.Println("⭐ Creating sample reviews...")
		reviews := []Review{
			{
				Rating:       5,
				Body:         "Absolutely breathtaking! The redwoods are magnificent and the campground is well-maintained.",
				AuthorID:     userIDs[1], // hannah
				CampgroundID: campgroundIDs[0], // Redwood
				CreatedAt:    time.Now().AddDate(0, 0, -5),
				UpdatedAt:    time.Now().AddDate(0, 0, -5),
			},
			{
				Rating:       5,
				Body:         "Best camping experience ever! The sunrise over the canyon was unforgettable.",
				AuthorID:     userIDs[0], // igoswamik
				CampgroundID: campgroundIDs[1], // Grand Canyon
				CreatedAt:    time.Now().AddDate(0, 0, -4),
				UpdatedAt:    time.Now().AddDate(0, 0, -4),
			},
			{
				Rating:       4,
				Body:         "Great location with amazing wildlife viewing opportunities. Saw elk and bison!",
				AuthorID:     userIDs[2], // bob
				CampgroundID: campgroundIDs[2], // Yellowstone
				CreatedAt:    time.Now().AddDate(0, 0, -3),
				UpdatedAt:    time.Now().AddDate(0, 0, -3),
			},
		}
		
		reviewCollection := db.Collection("reviews")
		for i := range reviews {
			reviews[i].ID = primitive.NewObjectID()
			
			_, err := reviewCollection.InsertOne(context.Background(), reviews[i])
			if err != nil {
				log.Printf("❌ Error creating review: %v", err)
				continue
			}
			log.Printf("✅ Created review for campground: %s", campgroundIDs[i%len(campgroundIDs)].Hex())
		}
		log.Printf("✅ Successfully created %d reviews", len(reviews))
	}
	
	// Verify the data was created
	finalCount, err := db.Collection("campgrounds").CountDocuments(context.Background(), map[string]interface{}{})
	if err != nil {
		log.Printf("❌ Error verifying campground count: %v", err)
	} else {
		log.Printf("✅ Final campground count: %d", finalCount)
	}
	
	userCount, err := db.Collection("users").CountDocuments(context.Background(), map[string]interface{}{})
	if err != nil {
		log.Printf("❌ Error verifying user count: %v", err)
	} else {
		log.Printf("✅ Final user count: %d", userCount)
	}
	
	log.Println("🎉 Sample data seeding completed successfully!")
	log.Println("📊 Summary:")
	log.Printf("   - %d users created", len(userIDs))
	log.Printf("   - %d campgrounds created", len(campgroundIDs))
	log.Printf("   - Multiple reviews created")
	log.Println("")
	log.Println("🔐 Demo login credentials:")
	log.Println("   Username: igoswamik")
	log.Println("   Password: password123")
	log.Println("")
	log.Println("🌐 Visit http://localhost:3000/campgrounds to see the campgrounds!")
}
