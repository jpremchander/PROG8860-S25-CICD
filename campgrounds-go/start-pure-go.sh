#!/bin/bash

echo "🚀 Starting Pure Go YelpCamp Application"
echo "========================================"

# Navigate to project directory
if [ -d "campgrounds-go" ]; then
    cd campgrounds-go
    echo "✅ Navigated to campgrounds-go directory"
else
    echo "❌ campgrounds-go directory not found!"
    exit 1
fi

# Create .env file for pure Go setup
echo "📝 Creating .env file for pure Go setup..."
cat > .env << 'EOF'
# Environment Configuration
GIN_MODE=debug

# MongoDB Configuration
MONGO_HOST=localhost
MONGO_PORT=27017
MONGO_USER=yelpcamp_dev
MONGO_PASSWORD=dev_password_123
MONGO_DATABASE=yelpcamp_dev

# JWT Secret
JWT_SECRET=your_super_secret_jwt_key_here

# Cloudinary Configuration
CLOUDINARY_CLOUD_NAME=dl0zvwr9i
CLOUDINARY_API_KEY=719144713242759
CLOUDINARY_API_SECRET=tCIGtTFz8vvEem2Vl1HuOd_7xhE

# Server Configuration
PORT=3000

# AWS Configuration
AWS_REGION=us-east-1
S3_BUCKET=yelpcamp-artifacts-bucket
EOF

# Clean up any existing containers
echo "🧹 Cleaning up existing containers..."
docker stop yelpcamp-mongo 2>/dev/null || true
docker rm yelpcamp-mongo 2>/dev/null || true

# Start MongoDB only
echo "🗄️  Starting MongoDB container..."
docker run -d --name yelpcamp-mongo \
    -p 27017:27017 \
    -e MONGO_INITDB_ROOT_USERNAME=yelpcamp_dev \
    -e MONGO_INITDB_ROOT_PASSWORD=dev_password_123 \
    -e MONGO_INITDB_DATABASE=yelpcamp_dev \
    mongo:7.0

# Wait for MongoDB to start
echo "⏳ Waiting for MongoDB to start..."
sleep 15

# Download Go dependencies
echo "📦 Installing Go dependencies..."
go mod tidy

# Build the application
echo "🏗️  Building Go application..."
go build -o yelpcamp-app .

# Start the application
echo "🚀 Starting YelpCamp Go application..."
echo ""
echo "📱 Application will be available at:"
echo "   🌐 Web Interface: http://localhost:3000"
echo "   📡 API Health: http://localhost:3000/api/health"
echo "   🗄️  MongoDB: localhost:27017"
echo ""
echo "🛑 Press Ctrl+C to stop the application"
echo ""

./yelpcamp-app
