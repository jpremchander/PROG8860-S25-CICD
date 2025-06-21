package models

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
		log.Println("Sample data already exists, skipping seed")
		return
	}
	
	log.Println("🌱 Seeding sample data...")
	
	// Create sample users
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
	
	// Create sample campgrounds
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
			Description: "Experience breathtaking mountain views and pristine wilderness at this secluded campground. Perfect for hiking enthusiasts and nature photographers.",
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
			Description: "Wake up to stunning lake views and enjoy swimming, fishing, and kayaking. This family-friendly campground offers clean restrooms and showers.",
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
	}
	
	// Insert campgrounds
	var campgroundIDs []primitive.ObjectID
	for i := range campgrounds {
		// Set creation time to simulate different dates
		campgrounds[i].CreatedAt = time.Now().AddDate(0, 0, -(i+1)*2) // 2, 4, 6, 8, 10 days ago
		campgrounds[i].UpdatedAt = campgrounds[i].CreatedAt
		
		if err := campgrounds[i].Create(db); err != nil {
			log.Printf("Error creating campground %s: %v", campgrounds[i].Title, err)
			continue
		}
		campgroundIDs = append(campgroundIDs, campgrounds[i].ID)
		log.Printf("✅ Created campground: %s", campgrounds[i].Title)
	}
	
	// Create sample reviews
	if len(campgroundIDs) > 0 {
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
		}
		
		for i := range reviews {
			reviews[i].CreatedAt = time.Now().AddDate(0, 0, -(i+1)) // 1, 2, 3, 4, 5 days ago
			reviews[i].UpdatedAt = reviews[i].CreatedAt
			
			if err := reviews[i].Create(db); err != nil {
				log.Printf("Error creating review: %v", err)
				continue
			}
			log.Printf("✅ Created review for campground")
		}
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
