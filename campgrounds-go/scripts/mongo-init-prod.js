// MongoDB initialization script for production
const { ObjectId } = require("mongodb")
const db = db.getSiblingDB("yelpcamp_prod")

// Create collections
db.createCollection("users")
db.createCollection("campgrounds")
db.createCollection("reviews")

// Create indexes for better performance
db.users.createIndex({ username: 1 }, { unique: true })
db.users.createIndex({ email: 1 }, { unique: true })
db.campgrounds.createIndex({ author_id: 1 })
db.campgrounds.createIndex({ location: "text", title: "text", description: "text" })
db.campgrounds.createIndex({ created_at: -1 })
db.reviews.createIndex({ campground_id: 1 })
db.reviews.createIndex({ author_id: 1 })
db.reviews.createIndex({ created_at: -1 })

// Create admin user (password should be changed immediately in production)
const adminUser = {
  _id: ObjectId(),
  username: "admin",
  email: "admin@yelpcamp.com",
  password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // Change this password!
  created_at: new Date(),
  updated_at: new Date(),
}

db.users.insertOne(adminUser)

// Insert a few sample campgrounds for production (optional)
const sampleCampgrounds = [
  {
    _id: ObjectId(),
    title: "Welcome Campground",
    description:
      "A beautiful starter campground to showcase YelpCamp features. Feel free to add your own campgrounds and make this community grow!",
    location: "Yosemite National Park, CA",
    price: 25.0,
    images: [
      {
        url: "https://images.unsplash.com/photo-1504851149312-7a075b496cc7?w=800&h=600&fit=crop",
        filename: "welcome-campground.jpg",
        key: "yelpcamp/welcome-campground",
      },
    ],
    author_id: adminUser._id,
    created_at: new Date(),
    updated_at: new Date(),
  },
]

db.campgrounds.insertMany(sampleCampgrounds)

print("✅ Production database initialized")
print("⚠️  IMPORTANT: Change the admin password immediately!")
print("📊 Created:")
print("   - 1 admin user")
print("   - 1 sample campground")
print("")
print("🔐 Admin credentials:")
print("   Username: admin")
print("   Password: password123 (CHANGE THIS!)")
