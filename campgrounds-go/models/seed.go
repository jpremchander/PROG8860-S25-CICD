package models

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"yelpcamp-go/config"
)

// SeedData creates sample users, campgrounds, and reviews for development
func SeedData() {
	db := config.GetDB()
	
	// Check if data already exists
	count, err := db.Collection("campgrounds").CountDocuments(context.Background(), map[string]interface{}{})
	if err != nil {
		log.Printf("Error checking existing data: %v", err)
		return
	}
	
	if count > 0 {
		log.Printf("Sample data already exists (%d campgrounds), skipping seed", count)
		return
	}
	
	log.Println("🌱 Seeding sample data...")
	
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
	for i := range users {
		if err := users[i].Create(db); err != nil {
			log.Printf("Error creating user %s: %v", users[i].Username, err)
			continue
		}
		userIDs = append(userIDs, users[i].ID)
		log.Printf("✅ Created user: %s", users[i].Username)
	}
	
	if len(userIDs) == 0 {
		log.Println("❌ No users created, skipping campgrounds")
		return
	}
	
	// Create sample campgrounds with realistic data
	campgrounds := []Campground{
		{
			Title:       "Redwood, Flats",
			Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Esse ipsum incidunt repellendus corporis corrupti harum eum cum recusandae, eveniet distinctio a, saepe voluptatibus! Repudiandae, eosi Ea aliquid iure nihil id?",
			Location:    "Walnut Creek, California",
			Price:       14.00,
			AuthorID:    userIDs[0], // igoswamik
			Images: []Image{
				{
					URL:      "https://images.unsplash.com/photo-1504851149312-7a075b496cc7?w=800&h=600&fit=crop",
					Filename: "redwood-flats-1.jpg",
					Key:      "yelpcamp/redwood-flats-1",
				},
			},
		},
		{
			Title:       "Cascade, Cliffs",
			Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Esse ipsum incidunt repellendus corporis corrupti harum eum cum recusandae, eveniet distinctio a, saepe voluptatibus! Repudiandae, eosi Ea aliquid iure nihil id?",
			Location:    "Roswell, New Mexico",
			Price:       22.50,
			AuthorID:    userIDs[1], // hannah
			Images: []Image{
				{
					URL:      "https://images.unsplash.com/photo-1445308394109-4ec2920981b1?w=800&h=600&fit=crop",
					Filename: "cascade-cliffs-1.jpg",
					Key:      "yelpcamp/cascade-cliffs-1",
				},
			},
		},
		{
			Title:       "Petrified, Creekside",
			Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Esse ipsum incidunt repellendus corporis corrupti harum eum cum recusandae, eveniet distinctio a, saepe voluptatibus! Repudiandae, eosi Ea aliquid iure nihil id?",
			Location:    "Pensacola, Florida",
			Price:       18.75,
			AuthorID:    userIDs[2], // bob
			Images: []Image{
				{
					URL:      "https://images.unsplash.com/photo-1487730116645-74489c95b41b?w=800&h=600&fit=crop",
					Filename: "petrified-creekside-1.jpg",
					Key:      "yelpcamp/petrified-creekside-1",
				},
			},
		},
		{
			Title:       "Mountain View Retreat",
			Description: "Experience breathtaking mountain views and pristine wilderness at this secluded campground. Perfect for hiking enthusiasts and nature photographers. Features include fire pits, picnic tables, and access to hiking trails.",
			Location:    "Aspen, Colorado",
			Price:       35.00,
			AuthorID:    userIDs[3], // alice
			Images: []Image{
				{
					URL:      "https://images.unsplash.com/photo-1504280390367-361c6d9f38f4?w=800&h=600&fit=crop",
					Filename: "mountain-view-1.jpg",
					Key:      "yelpcamp/mountain-view-1",
				},
			},
		},
		{
			Title:       "Lakeside Paradise",
			Description: "Wake up to stunning lake views and enjoy swimming, fishing, and kayaking. This family-friendly campground offers clean restrooms, showers, and a camp store. Perfect for a relaxing weekend getaway.",
			Location:    "Lake Tahoe, Nevada",
			Price:       28.00,
			AuthorID:    userIDs[0], // igoswamik
			Images: []Image{
				{
					URL:      "https://images.unsplash.com/photo-1441974231531-c6227db76b6e?w=800&h=600&fit=crop",
					Filename: "lakeside-paradise-1.jpg",
					Key:      "yelpcamp/lakeside-paradise-1",
				},
			},
		},
		{
			Title:       "Desert Oasis",
			Description: "Discover the beauty of the desert landscape under star-filled skies. This unique campground offers a peaceful escape with stunning sunsets and sunrise views. Ideal for stargazing and desert photography.",
			Location:    "Sedona, Arizona",
			Price:       25.50,
			AuthorID:    userIDs[1], // hannah
			Images: []Image{
				{
					URL:      "https://images.unsplash.com/photo-1506905925346-21bda4d32df4?w=800&h=600&fit=crop",
					Filename: "desert-oasis-1.jpg",
					Key:      "yelpcamp/desert-oasis-1",
				},
			},
		},
		{
			Title:       "Forest Haven",
			Description: "Immerse yourself in old-growth forest with towering trees and peaceful hiking trails. This campground offers a true back-to-nature experience with minimal amenities for the adventurous camper.",
			Location:    "Olympic National Park, Washington",
			Price:       20.00,
			AuthorID:    userIDs[2], // bob
			Images: []Image{
				{
					URL:      "https://images.unsplash.com/photo-1571863533956-01c88e79957e?w=800&h=600&fit=crop",
					Filename: "forest-haven-1.jpg",
					Key:      "yelpcamp/forest-haven-1",
				},
			},
		},
		{
			Title:       "Coastal Bluffs",
			Description: "Camp on dramatic coastal bluffs with panoramic ocean views. Listen to the waves crash below while enjoying spectacular sunsets. Features include wind-resistant fire pits and ocean access trails.",
			Location:    "Big Sur, California",
			Price:       42.00,
			AuthorID:    userIDs[3], // alice
			Images: []Image{
				{
					URL:      "https://images.unsplash.com/photo-1571863533956-01c88e79957e?w=800&h=600&fit=crop",
					Filename: "coastal-bluffs-1.jpg",
					Key:      "yelpcamp/coastal-bluffs-1",
				},
			},
		},
	}
	
	// Insert campgrounds with different creation dates
	var campgroundIDs []primitive.ObjectID
	for i := range campgrounds {
		// Set creation time to simulate different dates
		campgrounds[i].CreatedAt = time.Now().AddDate(0, 0, -(i+1)*2) // 2, 4, 6, 8, 10, 12, 14, 16 days ago
		campgrounds[i].UpdatedAt = campgrounds[i].CreatedAt
		
		if err := campgrounds[i].Create(db); err != nil {
			log.Printf("Error creating campground %s: %v", campgrounds[i].Title, err)
			continue
		}
		campgroundIDs = append(campgroundIDs, campgrounds[i].ID)
		log.Printf("✅ Created campground: %s", campgrounds[i].Title)
	}
	
	// Create sample reviews
	if len(campgroundIDs) > 0 && len(userIDs) > 0 {
		reviews := []Review{
			{
				Rating:       5,
				Body:         "I love this place.",
				AuthorID:     userIDs[1], // hannah
				CampgroundID: campgroundIDs[0], // Redwood, Flats
			},
			{
				Rating:       5,
				Body:         "Really cool place to visited. I would love to go again.",
				AuthorID:     userIDs[0], // igoswamik
				CampgroundID: campgroundIDs[0], // Redwood, Flats
			},
			{
				Rating:       5,
				Body:         "One of my favorite places to visit!",
				AuthorID:     userIDs[2], // bob
				CampgroundID: campgroundIDs[0], // Redwood, Flats
			},
			{
				Rating:       4,
				Body:         "Beautiful scenery and well-maintained facilities. The hiking trails are amazing!",
				AuthorID:     userIDs[3], // alice
				CampgroundID: campgroundIDs[1], // Cascade, Cliffs
			},
			{
				Rating:       5,
				Body:         "Perfect spot for a peaceful getaway. The creek sounds were so relaxing.",
				AuthorID:     userIDs[0], // igoswamik
				CampgroundID: campgroundIDs[2], // Petrified, Creekside
			},
			{
				Rating:       5,
				Body:         "Absolutely stunning mountain views! Worth every penny.",
				AuthorID:     userIDs[1], // hannah
				CampgroundID: campgroundIDs[3], // Mountain View Retreat
			},
			{
				Rating:       4,
				Body:         "Great for families! Kids loved swimming in the lake.",
				AuthorID:     userIDs[2], // bob
				CampgroundID: campgroundIDs[4], // Lakeside Paradise
			},
		}
		
		for i := range reviews {
			reviews[i].CreatedAt = time.Now().AddDate(0, 0, -(i+1)) // 1, 2, 3, 4, 5, 6, 7 days ago
			reviews[i].UpdatedAt = reviews[i].CreatedAt
			
			if err := reviews[i].Create(db); err != nil {
				log.Printf("Error creating review: %v", err)
				continue
			}
		}
		log.Printf("✅ Created %d reviews", len(reviews))
	}
	
	log.Println("🎉 Sample data seeding completed!")
	log.Println("📊 Created:")
	log.Printf("   - %d users", len(userIDs))
	log.Printf("   - %d campgrounds", len(campgroundIDs))
	log.Println("   - Multiple reviews")
	log.Println("")
	log.Println("🔐 All sample users have password: password123")
	log.Println("🌐 Visit /campgrounds to see the sample data!")
}

// ForceSeedData clears existing data and creates fresh sample data
func ForceSeedData() {
	db := config.GetDB()
	
	log.Println("🧹 Clearing existing data...")
	
	// Clear existing data
	db.Collection("reviews").DeleteMany(context.Background(), map[string]interface{}{})
	db.Collection("campgrounds").DeleteMany(context.Background(), map[string]interface{}{})
	db.Collection("users").DeleteMany(context.Background(), map[string]interface{}{})
	
	log.Println("✅ Cleared existing data")
	
	// Seed fresh data
	SeedData()
}
