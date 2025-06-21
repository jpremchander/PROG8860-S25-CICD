package models

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	_, err1 := db.Collection("reviews").DeleteMany(ctx, bson.M{})
	_, err2 := db.Collection("campgrounds").DeleteMany(ctx, bson.M{})
	_, err3 := db.Collection("users").DeleteMany(ctx, bson.M{})
	
	if err1 != nil || err2 != nil || err3 != nil {
		log.Printf("⚠️ Some errors during cleanup: reviews=%v, campgrounds=%v, users=%v", err1, err2, err3)
	}
	
	log.Println("✅ Cleared existing data")
	
	// Seed fresh data
	if err := SeedDatabase(db); err != nil {
		log.Printf("❌ Error seeding database: %v", err)
	}
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
			ID:       primitive.NewObjectID(),
			Username: "igoswamik",
			Email:    "igor@yelpcamp.com",
			Password: "password123",
		},
		{
			ID:       primitive.NewObjectID(),
			Username: "hannah",
			Email:    "hannah@yelpcamp.com",
			Password: "password123",
		},
		{
			ID:       primitive.NewObjectID(),
			Username: "bob",
			Email:    "bob@yelpcamp.com",
			Password: "password123",
		},
		{
			ID:       primitive.NewObjectID(),
			Username: "alice",
			Email:    "alice@yelpcamp.com",
			Password: "password123",
		},
		{
			ID:       primitive.NewObjectID(),
			Username: "charlie",
			Email:    "charlie@yelpcamp.com",
			Password: "password123",
		},
	}
	
	// Insert users directly
	var userIDs []primitive.ObjectID
	log.Println("👥 Creating sample users...")
	
	userCollection := db.Collection("users")
	for i := range users {
		// Hash password
		if err := users[i].HashPassword(); err != nil {
			log.Printf("❌ Error hashing password for user %s: %v", users[i].Username, err)
			continue
		}
		
		users[i].CreatedAt = time.Now()
		users[i].UpdatedAt = time.Now()
		
		_, err := userCollection.InsertOne(context.Background(), users[i])
		if err != nil {
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
	
	// Create sample campgrounds with high-quality images
	campgrounds := []Campground{
		{
			ID:          primitive.NewObjectID(),
			Title:       "Redwood National Park",
			Description: "Experience the majesty of the world's tallest trees in this pristine wilderness setting. Perfect for hiking, photography, and connecting with nature among ancient giants.",
			Location:    "Crescent City, California",
			Price:       45.00,
			AuthorID:    userIDs[0], // igoswamik
			Images: []Image{
				{
					URL:      "https://images.unsplash.com/photo-1441974231531-c6227db76b6e?w=800&h=600&fit=crop&crop=center",
					Filename: "redwood-forest.jpg",
					Key:      "yelpcamp/redwood-1",
				},
			},
			CreatedAt: time.Now().AddDate(0, 0, -15),
			UpdatedAt: time.Now().AddDate(0, 0, -15),
		},
		{
			ID:          primitive.NewObjectID(),
			Title:       "Grand Canyon South Rim",
			Description: "Wake up to breathtaking views of one of the world's natural wonders. This campground offers unparalleled sunrise and sunset viewing opportunities with full amenities.",
			Location:    "Grand Canyon, Arizona",
			Price:       55.00,
			AuthorID:    userIDs[1], // hannah
			Images: []Image{
				{
					URL:      "https://images.unsplash.com/photo-1506905925346-21bda4d32df4?w=800&h=600&fit=crop&crop=center",
					Filename: "grand-canyon-view.jpg",
					Key:      "yelpcamp/grand-canyon-1",
				},
			},
			CreatedAt: time.Now().AddDate(0, 0, -12),
			UpdatedAt: time.Now().AddDate(0, 0, -12),
		},
		{
			ID:          primitive.NewObjectID(),
			Title:       "Yellowstone Lake Lodge",
			Description: "Camp beside the pristine waters of Yellowstone Lake with opportunities for fishing, boating, and wildlife viewing. Hot springs nearby for a relaxing soak!",
			Location:    "Yellowstone National Park, Wyoming",
			Price:       60.00,
			AuthorID:    userIDs[2], // bob
			Images: []Image{
				{
					URL:      "https://images.unsplash.com/photo-1504851149312-7a075b496cc7?w=800&h=600&fit=crop&crop=center",
					Filename: "yellowstone-lake.jpg",
					Key:      "yelpcamp/yellowstone-1",
				},
			},
			CreatedAt: time.Now().AddDate(0, 0, -10),
			UpdatedAt: time.Now().AddDate(0, 0, -10),
		},
		{
			ID:          primitive.NewObjectID(),
			Title:       "Yosemite Valley Floor",
			Description: "Camp in the heart of Yosemite Valley with iconic views of El Capitan and Half Dome. Rock climbing and hiking trails accessible directly from camp.",
			Location:    "Yosemite National Park, California",
			Price:       65.00,
			AuthorID:    userIDs[3], // alice
			Images: []Image{
				{
					URL:      "https://images.unsplash.com/photo-1508873696983-2dfd5898f08b?w=800&h=600&fit=crop&crop=center",
					Filename: "yosemite-valley.jpg",
					Key:      "yelpcamp/yosemite-1",
				},
			},
			CreatedAt: time.Now().AddDate(0, 0, -8),
			UpdatedAt: time.Now().AddDate(0, 0, -8),
		},
		{
			ID:          primitive.NewObjectID(),
			Title:       "Glacier Point Overlook",
			Description: "Spectacular mountain camping with panoramic views of snow-capped peaks and alpine lakes. Perfect for stargazing and photography enthusiasts.",
			Location:    "Glacier National Park, Montana",
			Price:       50.00,
			AuthorID:    userIDs[4], // charlie
			Images: []Image{
				{
					URL:      "https://images.unsplash.com/photo-1445308396983-4ec2920981b1?w=800&h=600&fit=crop&crop=center",
					Filename: "glacier-mountains.jpg",
					Key:      "yelpcamp/glacier-1",
				},
			},
			CreatedAt: time.Now().AddDate(0, 0, -6),
			UpdatedAt: time.Now().AddDate(0, 0, -6),
		},
		{
			ID:          primitive.NewObjectID(),
			Title:       "Zion Canyon Riverside",
			Description: "Camp along the Virgin River with towering red rock formations surrounding you. Easy access to hiking trails and the famous Narrows hike.",
			Location:    "Zion National Park, Utah",
			Price:       48.00,
			AuthorID:    userIDs[0], // igoswamik
			Images: []Image{
				{
					URL:      "https://images.unsplash.com/photo-1487730116645-74489c95b41b?w=800&h=600&fit=crop&crop=center",
					Filename: "zion-canyon.jpg",
					Key:      "yelpcamp/zion-1",
				},
			},
			CreatedAt: time.Now().AddDate(0, 0, -4),
			UpdatedAt: time.Now().AddDate(0, 0, -4),
		},
		{
			ID:          primitive.NewObjectID(),
			Title:       "Arches National Park",
			Description: "Unique desert camping experience surrounded by natural stone arches and formations. Incredible night sky viewing with minimal light pollution.",
			Location:    "Moab, Utah",
			Price:       42.00,
			AuthorID:    userIDs[1], // hannah
			Images: []Image{
				{
					URL:      "https://images.unsplash.com/photo-1434394354979-a235cd36269d?w=800&h=600&fit=crop&crop=center",
					Filename: "arches-desert.jpg",
					Key:      "yelpcamp/arches-1",
				},
			},
			CreatedAt: time.Now().AddDate(0, 0, -2),
			UpdatedAt: time.Now().AddDate(0, 0, -2),
		},
		{
			ID:          primitive.NewObjectID(),
			Title:       "Olympic Peninsula Rainforest",
			Description: "Immerse yourself in the lush temperate rainforest of the Olympic Peninsula. Moss-covered trees and pristine streams create a magical atmosphere.",
			Location:    "Olympic National Park, Washington",
			Price:       38.00,
			AuthorID:    userIDs[2], // bob
			Images: []Image{
				{
					URL:      "https://images.unsplash.com/photo-1518837695005-2083093ee35b?w=800&h=600&fit=crop&crop=center",
					Filename: "olympic-rainforest.jpg",
					Key:      "yelpcamp/olympic-1",
				},
			},
			CreatedAt: time.Now().AddDate(0, 0, -1),
			UpdatedAt: time.Now().AddDate(0, 0, -1),
		},
	}
	
	// Insert campgrounds directly
	var campgroundIDs []primitive.ObjectID
	log.Println("🏕️ Creating sample campgrounds...")
	
	campgroundCollection := db.Collection("campgrounds")
	for i := range campgrounds {
		log.Printf("Creating campground: %s", campgrounds[i].Title)
		
		_, err := campgroundCollection.InsertOne(context.Background(), campgrounds[i])
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
				ID:           primitive.NewObjectID(),
				Rating:       5,
				Body:         "Absolutely breathtaking! The redwoods are magnificent and the campground is well-maintained. Perfect for families.",
				AuthorID:     userIDs[1], // hannah
				CampgroundID: campgroundIDs[0], // Redwood
				CreatedAt:    time.Now().AddDate(0, 0, -5),
				UpdatedAt:    time.Now().AddDate(0, 0, -5),
			},
			{
				ID:           primitive.NewObjectID(),
				Rating:       5,
				Body:         "Best camping experience ever! The sunrise over the canyon was unforgettable. Highly recommend!",
				AuthorID:     userIDs[0], // igoswamik
				CampgroundID: campgroundIDs[1], // Grand Canyon
				CreatedAt:    time.Now().AddDate(0, 0, -4),
				UpdatedAt:    time.Now().AddDate(0, 0, -4),
			},
			{
				ID:           primitive.NewObjectID(),
				Rating:       4,
				Body:         "Great location with amazing wildlife viewing opportunities. Saw elk and bison! Clean facilities.",
				AuthorID:     userIDs[2], // bob
				CampgroundID: campgroundIDs[2], // Yellowstone
				CreatedAt:    time.Now().AddDate(0, 0, -3),
				UpdatedAt:    time.Now().AddDate(0, 0, -3),
			},
			{
				ID:           primitive.NewObjectID(),
				Rating:       5,
				Body:         "Yosemite never disappoints! The views of Half Dome from the campsite are incredible. Will definitely return.",
				AuthorID:     userIDs[3], // alice
				CampgroundID: campgroundIDs[3], // Yosemite
				CreatedAt:    time.Now().AddDate(0, 0, -2),
				UpdatedAt:    time.Now().AddDate(0, 0, -2),
			},
		}
		
		reviewCollection := db.Collection("reviews")
		for i := range reviews {
			_, err := reviewCollection.InsertOne(context.Background(), reviews[i])
			if err != nil {
				log.Printf("❌ Error creating review: %v", err)
				continue
			}
			log.Printf("✅ Created review for campground")
		}
		log.Printf("✅ Successfully created %d reviews", len(reviews))
	}
	
	// Verify the data was created
	finalCount, err := db.Collection("campgrounds").CountDocuments(context.Background(), bson.M{})
	if err != nil {
		log.Printf("❌ Error verifying campground count: %v", err)
	} else {
		log.Printf("✅ Final campground count: %d", finalCount)
	}
	
	userCount, err := db.Collection("users").CountDocuments(context.Background(), bson.M{})
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

func SeedDatabase(db *mongo.Database) error {
	log.Println("🌱 Starting database seeding...")

	// Check if campgrounds already exist
	campgroundsCollection := db.Collection("campgrounds")
	count, err := campgroundsCollection.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		log.Printf("❌ Error counting campgrounds: %v", err)
		return err
	}

	if count > 0 {
		log.Printf("✅ Database already has %d campgrounds, skipping seed", count)
		return nil
	}

	// Create a test user first
	usersCollection := db.Collection("users")
	userCount, err := usersCollection.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		log.Printf("❌ Error counting users: %v", err)
		return err
	}

	var testUserID primitive.ObjectID
	if userCount == 0 {
		// Create test user
		testUser := User{
			ID:       primitive.NewObjectID(),
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123", // This will be hashed by the Create method
		}

		if err := testUser.Create(db); err != nil {
			log.Printf("❌ Error creating test user: %v", err)
			return err
		}
		testUserID = testUser.ID
		log.Printf("✅ Created test user: %s", testUser.Username)
	} else {
		// Get existing user
		var existingUser User
		err := usersCollection.FindOne(context.Background(), bson.M{}).Decode(&existingUser)
		if err != nil {
			log.Printf("❌ Error finding existing user: %v", err)
			return err
		}
		testUserID = existingUser.ID
		log.Printf("✅ Using existing user: %s", existingUser.Username)
	}

	// Sample campgrounds data with better images
	campgrounds := []Campground{
		{
			ID:          primitive.NewObjectID(),
			Title:       "Mountain View Campground",
			Description: "A beautiful campground with stunning mountain views and hiking trails nearby. Perfect for families and outdoor enthusiasts.",
			Location:    "Rocky Mountain National Park, Colorado",
			Price:       45.00,
			AuthorID:    testUserID,
			Images: []Image{
				{
					URL:      "https://images.unsplash.com/photo-1504280390367-361c6d9f38f4?w=800&h=600&fit=crop",
					Filename: "mountain-campground.jpg",
					Key:      "mountain-campground-1",
				},
				{
					URL:      "https://images.unsplash.com/photo-1441974231531-c6227db76b6e?w=800&h=600&fit=crop",
					Filename: "mountain-view.jpg",
					Key:      "mountain-view-1",
				},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:          primitive.NewObjectID(),
			Title:       "Lakeside Paradise",
			Description: "Peaceful lakeside camping with crystal clear water, fishing opportunities, and beautiful sunsets. Ideal for a relaxing getaway.",
			Location:    "Lake Tahoe, California",
			Price:       55.00,
			AuthorID:    testUserID,
			Images: []Image{
				{
					URL:      "https://images.unsplash.com/photo-1506905925346-21bda4d32df4?w=800&h=600&fit=crop",
					Filename: "lakeside-camp.jpg",
					Key:      "lakeside-camp-1",
				},
				{
					URL:      "https://images.unsplash.com/photo-1486022119026-e4b9e0c5b7e9?w=800&h=600&fit=crop",
					Filename: "lake-sunset.jpg",
					Key:      "lake-sunset-1",
				},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:          primitive.NewObjectID(),
			Title:       "Forest Haven",
			Description: "Deep in the forest with towering trees, wildlife viewing, and peaceful nature sounds. A true escape from city life.",
			Location:    "Olympic National Forest, Washington",
			Price:       35.00,
			AuthorID:    testUserID,
			Images: []Image{
				{
					URL:      "https://images.unsplash.com/photo-1508873696983-2dfd5898f08b?w=800&h=600&fit=crop",
					Filename: "forest-camp.jpg",
					Key:      "forest-camp-1",
				},
				{
					URL:      "https://images.unsplash.com/photo-1441974231531-c6227db76b6e?w=800&h=600&fit=crop",
					Filename: "forest-trees.jpg",
					Key:      "forest-trees-1",
				},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:          primitive.NewObjectID(),
			Title:       "Desert Oasis",
			Description: "Unique desert camping experience with stargazing opportunities, cacti gardens, and stunning sunrise views.",
			Location:    "Joshua Tree National Park, California",
			Price:       40.00,
			AuthorID:    testUserID,
			Images: []Image{
				{
					URL:      "https://images.unsplash.com/photo-1506905925346-21bda4d32df4?w=800&h=600&fit=crop",
					Filename: "desert-camp.jpg",
					Key:      "desert-camp-1",
				},
				{
					URL:      "https://images.unsplash.com/photo-1445308396983-4ec2920981b1?w=800&h=600&fit=crop",
					Filename: "desert-stars.jpg",
					Key:      "desert-stars-1",
				},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:          primitive.NewObjectID(),
			Title:       "Riverside Retreat",
			Description: "Camping along a gentle river with fishing, kayaking, and swimming opportunities. Great for water activities.",
			Location:    "Yellowstone National Park, Wyoming",
			Price:       50.00,
			AuthorID:    testUserID,
			Images: []Image{
				{
					URL:      "https://images.unsplash.com/photo-1486022116645-74489c95b41b?w=800&h=600&fit=crop",
					Filename: "riverside-camp.jpg",
					Key:      "riverside-camp-1",
				},
				{
					URL:      "https://images.unsplash.com/photo-1508873696983-2dfd5898f08b?w=800&h=600&fit=crop",
					Filename: "river-view.jpg",
					Key:      "river-view-1",
				},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Insert campgrounds
	var campgroundDocs []interface{}
	for _, campground := range campgrounds {
		campgroundDocs = append(campgroundDocs, campground)
	}

	result, err := campgroundsCollection.InsertMany(context.Background(), campgroundDocs)
	if err != nil {
		log.Printf("❌ Error inserting campgrounds: %v", err)
		return err
	}

	log.Printf("✅ Successfully seeded %d campgrounds", len(result.InsertedIDs))

	// Verify the data was inserted
	finalCount, err := campgroundsCollection.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		log.Printf("⚠️ Warning: Could not verify campground count: %v", err)
	} else {
		log.Printf("🎉 Database now contains %d campgrounds total", finalCount)
	}

	return nil
}
