#!/bin/bash

echo "🔍 YelpCamp Database Debug Script"
echo "================================="

# Check if app is running
echo "1. Checking if application is running..."
if curl -s http://localhost:3000/health > /dev/null; then
    echo "✅ Application is running on port 3000"
else
    echo "❌ Application is not responding on port 3000"
    echo "Please start the application first"
    exit 1
fi

# Check MongoDB container
echo ""
echo "2. Checking MongoDB container..."
if docker ps | grep -q mongo; then
    MONGO_CONTAINER=$(docker ps --format "table {{.Names}}" | grep mongo | head -1)
    echo "✅ MongoDB container found: $MONGO_CONTAINER"
else
    echo "❌ No MongoDB container running"
    echo "Starting MongoDB container..."
    docker run -d --name yelpcamp-mongo \
        -p 27017:27017 \
        -e MONGO_INITDB_ROOT_USERNAME=yelpcamp_dev \
        -e MONGO_INITDB_ROOT_PASSWORD=dev_password_123 \
        -e MONGO_INITDB_DATABASE=yelpcamp_dev \
        mongo:7.0
    sleep 15
    MONGO_CONTAINER="yelpcamp-mongo"
fi

# Test MongoDB connection
echo ""
echo "3. Testing MongoDB connection..."
docker exec $MONGO_CONTAINER mongosh --eval "db.adminCommand('ping')" 2>/dev/null && echo "✅ MongoDB is responding" || echo "❌ MongoDB connection failed"

# Check database and collections
echo ""
echo "4. Checking database collections..."
echo "   Available databases:"
docker exec $MONGO_CONTAINER mongosh --eval "show dbs" 2>/dev/null || echo "   Could not list databases"

echo ""
echo "   Checking yelpcamp_dev database:"
docker exec $MONGO_CONTAINER mongosh yelpcamp_dev --eval "show collections" 2>/dev/null || echo "   Could not access yelpcamp_dev database"

# Check collection counts
echo ""
echo "5. Checking collection data counts..."
echo "   Users:"
docker exec $MONGO_CONTAINER mongosh yelpcamp_dev --eval "db.users.countDocuments()" 2>/dev/null || echo "   Could not count users"

echo "   Campgrounds:"
docker exec $MONGO_CONTAINER mongosh yelpcamp_dev --eval "db.campgrounds.countDocuments()" 2>/dev/null || echo "   Could not count campgrounds"

echo "   Reviews:"
docker exec $MONGO_CONTAINER mongosh yelpcamp_dev --eval "db.reviews.countDocuments()" 2>/dev/null || echo "   Could not count reviews"

# Show sample data
echo ""
echo "6. Sample campground data:"
docker exec $MONGO_CONTAINER mongosh yelpcamp_dev --eval "db.campgrounds.find().limit(3).forEach(printjson)" 2>/dev/null || echo "   Could not fetch campground data"

# Test API endpoints
echo ""
echo "7. Testing API endpoints..."
echo "   /api/health:"
curl -s http://localhost:3000/api/health | jq '.status' 2>/dev/null || echo "   API health check failed"

echo "   /api/campgrounds:"
CAMPGROUND_COUNT=$(curl -s http://localhost:3000/api/campgrounds | jq '.count' 2>/dev/null)
echo "   API reports $CAMPGROUND_COUNT campgrounds"

# Test web pages
echo ""
echo "8. Testing web pages..."
echo "   Homepage (/):"
curl -s -o /dev/null -w "%{http_code}" http://localhost:3000/

echo "   Campgrounds page (/campgrounds):"
curl -s -o /dev/null -w "%{http_code}" http://localhost:3000/campgrounds

echo ""
echo "================================="
echo "🎯 Diagnosis Summary:"

if [ "$CAMPGROUND_COUNT" = "0" ] || [ -z "$CAMPGROUND_COUNT" ]; then
    echo "❌ ISSUE FOUND: No campgrounds in database"
    echo ""
    echo "💡 Solutions:"
    echo "1. Run the force-seed script: ./force-seed-data.sh"
    echo "2. Restart the application to trigger auto-seeding"
    echo "3. Check MongoDB authentication settings"
else
    echo "✅ Database has $CAMPGROUND_COUNT campgrounds"
    echo "✅ The issue might be in the web interface display"
fi

echo ""
echo "🔧 Next steps:"
echo "1. Run: ./force-seed-data.sh"
echo "2. Visit: http://localhost:3000/campgrounds"
echo "3. Check browser console for JavaScript errors"
