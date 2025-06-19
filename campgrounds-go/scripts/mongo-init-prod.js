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
db.users.insertOne({
  _id: ObjectId(),
  username: "admin",
  email: "admin@yelpcamp.com",
  password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // Change this password!
  created_at: new Date(),
  updated_at: new Date(),
})

print("✅ Production database initialized")
print("⚠️  IMPORTANT: Change the admin password immediately!")
