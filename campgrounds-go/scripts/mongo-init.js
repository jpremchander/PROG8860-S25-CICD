// MongoDB initialization script for development
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

// Insert sample data for development
db.users.insertMany([
  {
    _id: ObjectId(),
    username: "devuser1",
    email: "dev1@yelpcamp.com",
    password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password123
    created_at: new Date(),
    updated_at: new Date(),
  },
  {
    _id: ObjectId(),
    username: "devuser2",
    email: "dev2@yelpcamp.com",
    password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password123
    created_at: new Date(),
    updated_at: new Date(),
  },
])

// Insert sample campgrounds
const user1Id = db.users.findOne({ username: "devuser1" })._id
const user2Id = db.users.findOne({ username: "devuser2" })._id

db.campgrounds.insertMany([
  {
    _id: ObjectId(),
    title: "Sunset Valley Campground",
    description:
      "A beautiful campground with stunning sunset views over the valley. Perfect for families and nature lovers.",
    location: "Yosemite National Park, CA",
    price: 35.99,
    images: [
      {
        url: "https://res.cloudinary.com/demo/image/upload/sample.jpg",
        filename: "sunset-valley-1.jpg",
        public_id: "yelpcamp/sunset-valley-1",
      },
    ],
    author_id: user1Id,
    created_at: new Date(),
    updated_at: new Date(),
  },
  {
    _id: ObjectId(),
    title: "Mountain Peak Retreat",
    description:
      "High altitude camping with breathtaking mountain views. Ideal for experienced campers seeking adventure.",
    location: "Rocky Mountain National Park, CO",
    price: 45.5,
    images: [
      {
        url: "https://res.cloudinary.com/demo/image/upload/mountain.jpg",
        filename: "mountain-peak-1.jpg",
        public_id: "yelpcamp/mountain-peak-1",
      },
    ],
    author_id: user2Id,
    created_at: new Date(),
    updated_at: new Date(),
  },
])

print("✅ Development database initialized with sample data")
