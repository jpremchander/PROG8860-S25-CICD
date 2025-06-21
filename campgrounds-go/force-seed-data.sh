#!/bin/bash

echo "🌱 Force Seeding YelpCamp Database"
echo "=================================="

# Check if MongoDB is running
if ! docker ps | grep -q mongo; then
    echo "Starting MongoDB container..."
    docker run -d --name yelpcamp-mongo \
        -p 27017:27017 \
        -e MONGO_INITDB_ROOT_USERNAME=yelpcamp_dev \
        -e MONGO_INITDB_ROOT_PASSWORD=dev_password_123 \
        -e MONGO_INITDB_DATABASE=yelpcamp_dev \
        mongo:7.0
    sleep 15
fi

MONGO_CONTAINER=$(docker ps --format "table {{.Names}}" | grep mongo | head -1)

# Clear existing data
echo "1. Clearing existing data..."
docker exec $MONGO_CONTAINER mongosh yelpcamp_dev --eval "
db.reviews.deleteMany({});
db.campgrounds.deleteMany({});
db.users.deleteMany({});
print('✅ Cleared existing data');
" 2>/dev/null

# Create indexes
echo ""
echo "2. Creating database indexes..."
docker exec $MONGO_CONTAINER mongosh yelpcamp_dev --eval "
db.users.createIndex({ username: 1 }, { unique: true });
db.users.createIndex({ email: 1 }, { unique: true });
db.campgrounds.createIndex({ author_id: 1 });
db.campgrounds.createIndex({ created_at: -1 });
db.reviews.createIndex({ campground_id: 1 });
db.reviews.createIndex({ author_id: 1 });
print('✅ Created indexes');
" 2>/dev/null

# Insert users with hashed passwords
echo ""
echo "3. Creating sample users..."
docker exec $MONGO_CONTAINER mongosh yelpcamp_dev --eval "
const users = [
  {
    _id: ObjectId(),
    username: 'igoswamik',
    email: 'igor@yelpcamp.com',
    password: '\$2a\$10\$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    _id: ObjectId(),
    username: 'hannah',
    email: 'hannah@yelpcamp.com',
    password: '\$2a\$10\$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    _id: ObjectId(),
    username: 'bob',
    email: 'bob@yelpcamp.com',
    password: '\$2a\$10\$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    _id: ObjectId(),
    username: 'alice',
    email: 'alice@yelpcamp.com',
    password: '\$2a\$10\$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
    created_at: new Date(),
    updated_at: new Date()
  }
];

db.users.insertMany(users);
print('✅ Created ' + users.length + ' users');
" 2>/dev/null

# Get user IDs and create campgrounds
echo ""
echo "4. Creating sample campgrounds..."
docker exec $MONGO_CONTAINER mongosh yelpcamp_dev --eval "
const igoswamik = db.users.findOne({ username: 'igoswamik' })._id;
const hannah = db.users.findOne({ username: 'hannah' })._id;
const bob = db.users.findOne({ username: 'bob' })._id;
const alice = db.users.findOne({ username: 'alice' })._id;

const campgrounds = [
  {
    _id: ObjectId(),
    title: 'Redwood, Flats',
    description: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Esse ipsum incidunt repellendus corporis corrupti harum eum cum recusandae, eveniet distinctio a, saepe voluptatibus! Repudiandae, eosi Ea aliquid iure nihil id?',
    location: 'Walnut Creek, California',
    price: 14.00,
    images: [{
      url: 'https://images.unsplash.com/photo-1504851149312-7a075b496cc7?w=800&h=600&fit=crop',
      filename: 'redwood-flats-1.jpg',
      key: 'yelpcamp/redwood-flats-1'
    }],
    author_id: igoswamik,
    created_at: new Date(Date.now() - 2 * 24 * 60 * 60 * 1000),
    updated_at: new Date(Date.now() - 2 * 24 * 60 * 60 * 1000)
  },
  {
    _id: ObjectId(),
    title: 'Cascade, Cliffs',
    description: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Esse ipsum incidunt repellendus corporis corrupti harum eum cum recusandae, eveniet distinctio a, saepe voluptatibus! Repudiandae, eosi Ea aliquid iure nihil id?',
    location: 'Roswell, New Mexico',
    price: 22.50,
    images: [{
      url: 'https://images.unsplash.com/photo-1445308394109-4ec2920981b1?w=800&h=600&fit=crop',
      filename: 'cascade-cliffs-1.jpg',
      key: 'yelpcamp/cascade-cliffs-1'
    }],
    author_id: hannah,
    created_at: new Date(Date.now() - 5 * 24 * 60 * 60 * 1000),
    updated_at: new Date(Date.now() - 5 * 24 * 60 * 60 * 1000)
  },
  {
    _id: ObjectId(),
    title: 'Petrified, Creekside',
    description: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Esse ipsum incidunt repellendus corporis corrupti harum eum cum recusandae, eveniet distinctio a, saepe voluptatibus! Repudiandae, eosi Ea aliquid iure nihil id?',
    location: 'Pensacola, Florida',
    price: 18.75,
    images: [{
      url: 'https://images.unsplash.com/photo-1487730116645-74489c95b41b?w=800&h=600&fit=crop',
      filename: 'petrified-creekside-1.jpg',
      key: 'yelpcamp/petrified-creekside-1'
    }],
    author_id: bob,
    created_at: new Date(Date.now() - 7 * 24 * 60 * 60 * 1000),
    updated_at: new Date(Date.now() - 7 * 24 * 60 * 60 * 1000)
  },
  {
    _id: ObjectId(),
    title: 'Mountain View Retreat',
    description: 'Experience breathtaking mountain views and pristine wilderness at this secluded campground. Perfect for hiking enthusiasts and nature photographers. Features include fire pits, picnic tables, and access to hiking trails.',
    location: 'Aspen, Colorado',
    price: 35.00,
    images: [{
      url: 'https://images.unsplash.com/photo-1504280390367-361c6d9f38f4?w=800&h=600&fit=crop',
      filename: 'mountain-view-1.jpg',
      key: 'yelpcamp/mountain-view-1'
    }],
    author_id: alice,
    created_at: new Date(Date.now() - 10 * 24 * 60 * 60 * 1000),
    updated_at: new Date(Date.now() - 10 * 24 * 60 * 60 * 1000)
  },
  {
    _id: ObjectId(),
    title: 'Lakeside Paradise',
    description: 'Wake up to stunning lake views and enjoy swimming, fishing, and kayaking. This family-friendly campground offers clean restrooms, showers, and a camp store. Perfect for a relaxing weekend getaway.',
    location: 'Lake Tahoe, Nevada',
    price: 28.00,
    images: [{
      url: 'https://images.unsplash.com/photo-1441974231531-c6227db76b6e?w=800&h=600&fit=crop',
      filename: 'lakeside-paradise-1.jpg',
      key: 'yelpcamp/lakeside-paradise-1'
    }],
    author_id: igoswamik,
    created_at: new Date(Date.now() - 12 * 24 * 60 * 60 * 1000),
    updated_at: new Date(Date.now() - 12 * 24 * 60 * 60 * 1000)
  },
  {
    _id: ObjectId(),
    title: 'Desert Oasis',
    description: 'Discover the beauty of the desert landscape under star-filled skies. This unique campground offers a peaceful escape with stunning sunsets and sunrise views. Ideal for stargazing and desert photography.',
    location: 'Sedona, Arizona',
    price: 25.50,
    images: [{
      url: 'https://images.unsplash.com/photo-1506905925346-21bda4d32df4?w=800&h=600&fit=crop',
      filename: 'desert-oasis-1.jpg',
      key: 'yelpcamp/desert-oasis-1'
    }],
    author_id: hannah,
    created_at: new Date(Date.now() - 15 * 24 * 60 * 60 * 1000),
    updated_at: new Date(Date.now() - 15 * 24 * 60 * 60 * 1000)
  },
  {
    _id: ObjectId(),
    title: 'Forest Haven',
    description: 'Immerse yourself in old-growth forest with towering trees and peaceful hiking trails. This campground offers a true back-to-nature experience with minimal amenities for the adventurous camper.',
    location: 'Olympic National Park, Washington',
    price: 20.00,
    images: [{
      url: 'https://images.unsplash.com/photo-1571863533956-01c88e79957e?w=800&h=600&fit=crop',
      filename: 'forest-haven-1.jpg',
      key: 'yelpcamp/forest-haven-1'
    }],
    author_id: bob,
    created_at: new Date(Date.now() - 18 * 24 * 60 * 60 * 1000),
    updated_at: new Date(Date.now() - 18 * 24 * 60 * 60 * 1000)
  },
  {
    _id: ObjectId(),
    title: 'Coastal Bluffs',
    description: 'Camp on dramatic coastal bluffs with panoramic ocean views. Listen to the waves crash below while enjoying spectacular sunsets. Features include wind-resistant fire pits and ocean access trails.',
    location: 'Big Sur, California',
    price: 42.00,
    images: [{
      url: 'https://images.unsplash.com/photo-1571863533956-01c88e79957e?w=800&h=600&fit=crop',
      filename: 'coastal-bluffs-1.jpg',
      key: 'yelpcamp/coastal-bluffs-1'
    }],
    author_id: alice,
    created_at: new Date(Date.now() - 20 * 24 * 60 * 60 * 1000),
    updated_at: new Date(Date.now() - 20 * 24 * 60 * 60 * 1000)
  }
];

db.campgrounds.insertMany(campgrounds);
print('✅ Created ' + campgrounds.length + ' campgrounds');
" 2>/dev/null

# Create sample reviews
echo ""
echo "5. Creating sample reviews..."
docker exec $MONGO_CONTAINER mongosh yelpcamp_dev --eval "
const users = db.users.find().toArray();
const campgrounds = db.campgrounds.find().toArray();

const reviews = [
  {
    _id: ObjectId(),
    rating: 5,
    body: 'I love this place.',
    author_id: users[1]._id,
    campground_id: campgrounds[0]._id,
    created_at: new Date(Date.now() - 1 * 24 * 60 * 60 * 1000),
    updated_at: new Date(Date.now() - 1 * 24 * 60 * 60 * 1000)
  },
  {
    _id: ObjectId(),
    rating: 5,
    body: 'Really cool place to visited. I would love to go again.',
    author_id: users[0]._id,
    campground_id: campgrounds[0]._id,
    created_at: new Date(Date.now() - 2 * 24 * 60 * 60 * 1000),
    updated_at: new Date(Date.now() - 2 * 24 * 60 * 60 * 1000)
  },
  {
    _id: ObjectId(),
    rating: 5,
    body: 'One of my favorite places to visit!',
    author_id: users[2]._id,
    campground_id: campgrounds[0]._id,
    created_at: new Date(Date.now() - 3 * 24 * 60 * 60 * 1000),
    updated_at: new Date(Date.now() - 3 * 24 * 60 * 60 * 1000)
  },
  {
    _id: ObjectId(),
    rating: 4,
    body: 'Beautiful scenery and well-maintained facilities. The hiking trails are amazing!',
    author_id: users[3]._id,
    campground_id: campgrounds[1]._id,
    created_at: new Date(Date.now() - 4 * 24 * 60 * 60 * 1000),
    updated_at: new Date(Date.now() - 4 * 24 * 60 * 60 * 1000)
  },
  {
    _id: ObjectId(),
    rating: 5,
    body: 'Perfect spot for a peaceful getaway. The creek sounds were so relaxing.',
    author_id: users[0]._id,
    campground_id: campgrounds[2]._id,
    created_at: new Date(Date.now() - 6 * 24 * 60 * 60 * 1000),
    updated_at: new Date(Date.now() - 6 * 24 * 60 * 60 * 1000)
  }
];

db.reviews.insertMany(reviews);
print('✅ Created ' + reviews.length + ' reviews');
" 2>/dev/null

# Verify data
echo ""
echo "6. Verifying seeded data..."
USER_COUNT=$(docker exec $MONGO_CONTAINER mongosh yelpcamp_dev --eval "db.users.countDocuments()" --quiet 2>/dev/null)
CAMPGROUND_COUNT=$(docker exec $MONGO_CONTAINER mongosh yelpcamp_dev --eval "db.campgrounds.countDocuments()" --quiet 2>/dev/null)
REVIEW_COUNT=$(docker exec $MONGO_CONTAINER mongosh yelpcamp_dev --eval "db.reviews.countDocuments()" --quiet 2>/dev/null)

echo "   Users: $USER_COUNT"
echo "   Campgrounds: $CAMPGROUND_COUNT"
echo "   Reviews: $REVIEW_COUNT"

# Test API
echo ""
echo "7. Testing API after seeding..."
sleep 2
API_COUNT=$(curl -s http://localhost:3000/api/campgrounds | jq '.count' 2>/dev/null)
echo "   API reports: $API_COUNT campgrounds"

echo ""
echo "🎉 Database seeding completed!"
echo ""
echo "📊 Summary:"
echo "   ✅ $USER_COUNT users created (igoswamik, hannah, bob, alice)"
echo "   ✅ $CAMPGROUND_COUNT campgrounds created"
echo "   ✅ $REVIEW_COUNT reviews created"
echo ""
echo "🔐 Login credentials for testing:"
echo "   Username: igoswamik | Password: password123"
echo "   Username: hannah    | Password: password123"
echo "   Username: bob       | Password: password123"
echo "   Username: alice     | Password: password123"
echo ""
echo "🌐 Visit your application:"
echo "   Homepage: http://localhost:3000"
echo "   Campgrounds: http://localhost:3000/campgrounds"
echo "   Login: http://localhost:3000/login"
