#!/bin/bash

echo "🔄 Restarting YelpCamp with Fresh Data"
echo "======================================"

# Stop application if running
echo "1. Stopping application..."
pkill -f yelpcamp-app 2>/dev/null || echo "   No application process found"

# Stop and remove MongoDB container
echo ""
echo "2. Stopping MongoDB container..."
docker stop yelpcamp-mongo 2>/dev/null || echo "   No container to stop"
docker rm yelpcamp-mongo 2>/dev/null || echo "   No container to remove"

# Start fresh MongoDB
echo ""
echo "3. Starting fresh MongoDB container..."
docker run -d --name yelpcamp-mongo \
    -p 27017:27017 \
    -e MONGO_INITDB_ROOT_USERNAME=yelpcamp_dev \
    -e MONGO_INITDB_ROOT_PASSWORD=dev_password_123 \
    -e MONGO_INITDB_DATABASE=yelpcamp_dev \
    mongo:7.0

echo "   Waiting for MongoDB to start..."
sleep 15

# Force seed data
echo ""
echo "4. Seeding fresh data..."
./force-seed-data.sh

# Update .env for localhost
echo ""
echo "5. Updating .env for localhost..."
sed -i 's/MONGO_HOST=mongo/MONGO_HOST=localhost/' .env

# Build and start application
echo ""
echo "6. Building and starting application..."
go build -o yelpcamp-app .

if [ -f "./yelpcamp-app" ]; then
    echo "✅ Build successful!"
    echo ""
    echo "🚀 Starting YelpCamp..."
    echo "📱 Application will be available at:"
    echo "   🌐 Homepage: http://localhost:3000"
    echo "   🏕️  Campgrounds: http://localhost:3000/campgrounds"
    echo "   📡 API: http://localhost:3000/api/health"
    echo ""
    echo "🛑 Press Ctrl+C to stop"
    echo ""
    ./yelpcamp-app
else
    echo "❌ Build failed!"
fi
