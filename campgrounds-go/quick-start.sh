#!/bin/bash

echo "🧹 Clean Build for YelpCamp Go"
echo "=============================="

# Load environment variables
if [ -f ".env" ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Stop and remove containers
echo "🛑 Stopping containers..."
docker stop yelpcamp-mongo 2>/dev/null || true
docker rm yelpcamp-mongo 2>/dev/null || true

# Clean Go artifacts
echo "🧹 Cleaning Go artifacts..."
rm -f yelpcamp-app
rm -f go.sum
go clean -modcache

# Remove any Cloudinary references from go.mod
echo "📝 Cleaning go.mod..."
rm -f go.mod go.sum

# Initialize fresh Go module
go mod init yelpcamp-go

# Add only required dependencies
echo "📦 Adding dependencies..."
go get github.com/gin-gonic/gin@v1.9.1
go get github.com/joho/godotenv@v1.5.1
go get go.mongodb.org/mongo-driver@v1.13.1
go get golang.org/x/crypto@v0.17.0
go get github.com/golang-jwt/jwt/v4@v4.5.0
go get github.com/aws/aws-sdk-go@v1.55.7

# Tidy modules
go mod tidy

# Start MongoDB
echo "🗄️  Starting MongoDB..."
docker run -d --name yelpcamp-mongo \
    -p 27017:27017 \
    -e MONGO_INITDB_ROOT_USERNAME=yelpcamp_dev \
    -e MONGO_INITDB_ROOT_PASSWORD=dev_password_123 \
    -e MONGO_INITDB_DATABASE=yelpcamp_dev \
    mongo:7.0

# Wait for MongoDB
echo "⏳ Waiting for MongoDB..."
sleep 15

# Update .env for localhost
sed -i 's/MONGO_HOST=mongo/MONGO_HOST=localhost/' .env

# Build
echo "🏗️  Building..."
go build -o yelpcamp-app .

if [ -f "./yelpcamp-app" ]; then
    echo "✅ Build successful!"
    echo ""
    echo "🚀 Starting YelpCamp..."
    echo "📱 Web: http://localhost:3000"
    echo "📡 API: http://localhost:3000/api/health"
    echo "☁️  S3: ${S3_BUCKET}"
    echo ""
    ./yelpcamp-app
else
    echo "❌ Build failed!"
    go mod verify
fi
