// MongoDB initialization script for development
print("🌱 Initializing YelpCamp MongoDB...")

const { ObjectId } = require("mongodb")
const db = db.getSiblingDB("yelpcamp_dev")

// Create collections
db.createCollection("users")
db.createCollection("campgrounds")
db.createCollection("reviews")

// Create indexes for better performance
db.users.createIndex({ username: 1 }, { unique: true })
db.users.createIndex({ email: 1 }, { unique: true })
db.campgrounds.createIndex({ author_id: 1 })
db.campgrounds.createIndex({ location: "text", title: "text", description: "text" })
db.reviews.createIndex({ campground_id: 1 })
db.reviews.createIndex({ author_id: 1 })

print("✅ Created collections and indexes")

// Insert sample users with hashed passwords (password123)
const users = [
  {
    _id: ObjectId(),
    username: "igoswamik",
    email: "igor@yelpcamp.com",
    password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi",
    created_at: new Date(),
    updated_at: new Date(),
  },
  {
    _id: ObjectId(),
    username: "hannah",
    email: "hannah@yelpcamp.com",
    password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi",
    created_at: new Date(),
    updated_at: new Date(),
  },
  {
    _id: ObjectId(),
    username: "bob",
    email: "bob@yelpcamp.com",
    password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi",
    created_at: new Date(),
    updated_at: new Date(),
  },
  {
    _id: ObjectId(),
    username: "alice",
    email: "alice@yelpcamp.com",
    password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi",
    created_at: new Date(),
    updated_at: new Date(),
  },
]

db.users.insertMany(users)
print("✅ Created " + users.length + " sample users")

// Get user IDs for campground authors
const igoswamik = db.users.findOne({ username: "igoswamik" })._id
const hannah = db.users.findOne({ username: "hannah" })._id
const bob = db.users.findOne({ username: "bob" })._id
const alice = db.users.findOne({ username: "alice" })._id

// Insert realistic sample campgrounds matching the reference
const campgrounds = [
  {
    _id: ObjectId(),
    title: "Redwood, Flats",
    description:
      "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Esse ipsum incidunt repellendus corporis corrupti harum eum cum recusandae, eveniet distinctio a, saepe voluptatibus! Repudiandae, eosi Ea aliquid iure nihil id?",
    location: "Walnut Creek, California",
    price: 14.0,
    images: [
      {
        url: "https://images.unsplash.com/photo-1504851149312-7a075b496cc7?w=800&h=600&fit=crop",
        filename: "redwood-flats-1.jpg",
        key: "yelpcamp/redwood-flats-1",
      },
    ],
    author_id: igoswamik,
    created_at: new Date(Date.now() - 2 * 24 * 60 * 60 * 1000), // 2 days ago
    updated_at: new Date(Date.now() - 2 * 24 * 60 * 60 * 1000),
  },
  {
    _id: ObjectId(),
    title: "Cascade, Cliffs",
    description:
      "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Esse ipsum incidunt repellendus corporis corrupti harum eum cum recusandae, eveniet distinctio a, saepe voluptatibus! Repudiandae, eosi Ea aliquid iure nihil id?",
    location: "Roswell, New Mexico",
    price: 22.5,
    images: [
      {
        url: "https://images.unsplash.com/photo-1445308394109-4ec2920981b1?w=800&h=600&fit=crop",
        filename: "cascade-cliffs-1.jpg",
        key: "yelpcamp/cascade-cliffs-1",
      },
    ],
    author_id: hannah,
    created_at: new Date(Date.now() - 5 * 24 * 60 * 60 * 1000), // 5 days ago
    updated_at: new Date(Date.now() - 5 * 24 * 60 * 60 * 1000),
  },
  {
    _id: ObjectId(),
    title: "Petrified, Creekside",
    description:
      "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Esse ipsum incidunt repellendus corporis corrupti harum eum cum recusandae, eveniet distinctio a, saepe voluptatibus! Repudiandae, eosi Ea aliquid iure nihil id?",
    location: "Pensacola, Florida",
    price: 18.75,
    images: [
      {
        url: "https://images.unsplash.com/photo-1487730116645-74489c95b41b?w=800&h=600&fit=crop",
        filename: "petrified-creekside-1.jpg",
        key: "yelpcamp/petrified-creekside-1",
      },
    ],
    author_id: bob,
    created_at: new Date(Date.now() - 7 * 24 * 60 * 60 * 1000), // 7 days ago
    updated_at: new Date(Date.now() - 7 * 24 * 60 * 60 * 1000),
  },
  {
    _id: ObjectId(),
    title: "Mountain View Retreat",
    description:
      "Experience breathtaking mountain views and pristine wilderness at this secluded campground. Perfect for hiking enthusiasts and nature photographers. Features include fire pits, picnic tables, and access to hiking trails.",
    location: "Aspen, Colorado",
    price: 35.0,
    images: [
      {
        url: "https://images.unsplash.com/photo-1504280390367-361c6d9f38f4?w=800&h=600&fit=crop",
        filename: "mountain-view-1.jpg",
        key: "yelpcamp/mountain-view-1",
      },
    ],
    author_id: alice,
    created_at: new Date(Date.now() - 10 * 24 * 60 * 60 * 1000), // 10 days ago
    updated_at: new Date(Date.now() - 10 * 24 * 60 * 60 * 1000),
  },
  {
    _id: ObjectId(),
    title: "Lakeside Paradise",
    description:
      "Wake up to stunning lake views and enjoy swimming, fishing, and kayaking. This family-friendly campground offers clean restrooms, showers, and a camp store. Perfect for a relaxing weekend getaway.",
    location: "Lake Tahoe, Nevada",
    price: 28.0,
    images: [
      {
        url: "https://images.unsplash.com/photo-1441974231531-c6227db76b6e?w=800&h=600&fit=crop",
        filename: "lakeside-paradise-1.jpg",
        key: "yelpcamp/lakeside-paradise-1",
      },
    ],
    author_id: igoswamik,
    created_at: new Date(Date.now() - 12 * 24 * 60 * 60 * 1000), // 12 days ago
    updated_at: new Date(Date.now() - 12 * 24 * 60 * 60 * 1000),
  },
  {
    _id: ObjectId(),
    title: "Desert Oasis",
    description:
      "Discover the beauty of the desert landscape under star-filled skies. This unique campground offers a peaceful escape with stunning sunsets and sunrise views. Ideal for stargazing and desert photography.",
    location: "Sedona, Arizona",
    price: 25.5,
    images: [
      {
        url: "https://images.unsplash.com/photo-1506905925346-21bda4d32df4?w=800&h=600&fit=crop",
        filename: "desert-oasis-1.jpg",
        key: "yelpcamp/desert-oasis-1",
      },
    ],
    author_id: hannah,
    created_at: new Date(Date.now() - 15 * 24 * 60 * 60 * 1000), // 15 days ago
    updated_at: new Date(Date.now() - 15 * 24 * 60 * 60 * 1000),
  },
  {
    _id: ObjectId(),
    title: "Forest Haven",
    description:
      "Immerse yourself in old-growth forest with towering trees and peaceful hiking trails. This campground offers a true back-to-nature experience with minimal amenities for the adventurous camper.",
    location: "Olympic National Park, Washington",
    price: 20.0,
    images: [
      {
        url: "https://images.unsplash.com/photo-1571863533956-01c88e79957e?w=800&h=600&fit=crop",
        filename: "forest-haven-1.jpg",
        key: "yelpcamp/forest-haven-1",
      },
    ],
    author_id: bob,
    created_at: new Date(Date.now() - 18 * 24 * 60 * 60 * 1000), // 18 days ago
    updated_at: new Date(Date.now() - 18 * 24 * 60 * 60 * 1000),
  },
  {
    _id: ObjectId(),
    title: "Coastal Bluffs",
    description:
      "Camp on dramatic coastal bluffs with panoramic ocean views. Listen to the waves crash below while enjoying spectacular sunsets. Features include wind-resistant fire pits and ocean access trails.",
    location: "Big Sur, California",
    price: 42.0,
    images: [
      {
        url: "https://images.unsplash.com/photo-1571863533956-01c88e79957e?w=800&h=600&fit=crop",
        filename: "coastal-bluffs-1.jpg",
        key: "yelpcamp/coastal-bluffs-1",
      },
    ],
    author_id: alice,
    created_at: new Date(Date.now() - 20 * 24 * 60 * 60 * 1000), // 20 days ago
    updated_at: new Date(Date.now() - 20 * 24 * 60 * 60 * 1000),
  },
]

db.campgrounds.insertMany(campgrounds)
print("✅ Created " + campgrounds.length + " sample campgrounds")

// Get campground IDs for reviews
const redwoodFlats = db.campgrounds.findOne({ title: "Redwood, Flats" })._id
const cascadeCliffs = db.campgrounds.findOne({ title: "Cascade, Cliffs" })._id
const petrifiedCreekside = db.campgrounds.findOne({ title: "Petrified, Creekside" })._id
const mountainView = db.campgrounds.findOne({ title: "Mountain View Retreat" })._id
const lakesideParadise = db.campgrounds.findOne({ title: "Lakeside Paradise" })._id

// Insert sample reviews
const reviews = [
  {
    _id: ObjectId(),
    rating: 5,
    body: "I love this place.",
    author_id: hannah,
    campground_id: redwoodFlats,
    created_at: new Date(Date.now() - 1 * 24 * 60 * 60 * 1000), // 1 day ago
    updated_at: new Date(Date.now() - 1 * 24 * 60 * 60 * 1000),
  },
  {
    _id: ObjectId(),
    rating: 5,
    body: "Really cool place to visited. I would love to go again.",
    author_id: igoswamik,
    campground_id: redwoodFlats,
    created_at: new Date(Date.now() - 2 * 24 * 60 * 60 * 1000), // 2 days ago
    updated_at: new Date(Date.now() - 2 * 24 * 60 * 60 * 1000),
  },
  {
    _id: ObjectId(),
    rating: 5,
    body: "One of my favorite places to visit!",
    author_id: bob,
    campground_id: redwoodFlats,
    created_at: new Date(Date.now() - 3 * 24 * 60 * 60 * 1000), // 3 days ago
    updated_at: new Date(Date.now() - 3 * 24 * 60 * 60 * 1000),
  },
  {
    _id: ObjectId(),
    rating: 4,
    body: "Beautiful scenery and well-maintained facilities. The hiking trails are amazing!",
    author_id: alice,
    campground_id: cascadeCliffs,
    created_at: new Date(Date.now() - 4 * 24 * 60 * 60 * 1000), // 4 days ago
    updated_at: new Date(Date.now() - 4 * 24 * 60 * 60 * 1000),
  },
  {
    _id: ObjectId(),
    rating: 5,
    body: "Perfect spot for a peaceful getaway. The creek sounds were so relaxing.",
    author_id: igoswamik,
    campground_id: petrifiedCreekside,
    created_at: new Date(Date.now() - 6 * 24 * 60 * 60 * 1000), // 6 days ago
    updated_at: new Date(Date.now() - 6 * 24 * 60 * 60 * 1000),
  },
  {
    _id: ObjectId(),
    rating: 5,
    body: "Absolutely stunning mountain views! Worth every penny.",
    author_id: hannah,
    campground_id: mountainView,
    created_at: new Date(Date.now() - 8 * 24 * 60 * 60 * 1000), // 8 days ago
    updated_at: new Date(Date.now() - 8 * 24 * 60 * 60 * 1000),
  },
  {
    _id: ObjectId(),
    rating: 4,
    body: "Great for families! Kids loved swimming in the lake.",
    author_id: bob,
    campground_id: lakesideParadise,
    created_at: new Date(Date.now() - 11 * 24 * 60 * 60 * 1000), // 11 days ago
    updated_at: new Date(Date.now() - 11 * 24 * 60 * 60 * 1000),
  },
]

db.reviews.insertMany(reviews)
print("✅ Created " + reviews.length + " sample reviews")

print("🎉 YelpCamp MongoDB initialization completed!")
print("📊 Sample data created:")
print("   - " + users.length + " users (igoswamik, hannah, bob, alice)")
print("   - " + campgrounds.length + " campgrounds")
print("   - " + reviews.length + " reviews")
print("")
print("🔐 All users have password: password123")
print("🌐 You can now visit /campgrounds to see the sample data!")
