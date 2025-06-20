#!/bin/bash

echo "🚀 YelpCamp Quick Start Script"
echo "=============================="

# Navigate to project directory
if [ -d "campgrounds-go" ]; then
    cd campgrounds-go
    echo "✅ Navigated to campgrounds-go directory"
else
    echo "❌ campgrounds-go directory not found!"
    echo "Current directory contents:"
    ls -la
    exit 1
fi

# Check if .env exists
if [ ! -f ".env" ]; then
    echo "📝 Creating .env file..."
    cat > .env << 'EOF'
# Environment Configuration
GIN_MODE=debug

# MongoDB Configuration
MONGO_HOST=mongo
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
IAM_ROLE=arn:aws:iam::054037100649:role/yelpcamp-s3-role
EOF
    echo "✅ .env file created"
fi

# Clean up any existing containers
echo "🧹 Cleaning up existing containers..."
docker-compose -f docker-compose.dev.yml down 2>/dev/null || echo "No containers to stop"

# Download dependencies
echo "📦 Downloading Go dependencies..."
go mod tidy

# Start with Docker Compose
echo "🐳 Starting with Docker Compose..."
docker-compose -f docker-compose.dev.yml up --build -d

# Wait for services to start
echo "⏳ Waiting for services to start..."
sleep 30

# Check if services are running
echo "🔍 Checking service status..."
docker-compose -f docker-compose.dev.yml ps

# Test endpoints
echo "🧪 Testing endpoints..."
echo "Testing health endpoint:"
curl -s http://localhost:3000/api/health | jq . 2>/dev/null || curl -s http://localhost:3000/api/health

echo ""
echo "Testing root endpoint:"
curl -s -I http://localhost:3000/ | head -1

echo ""
echo "🎉 Setup complete!"
echo "📱 Access the application:"
echo "  - Web Interface: http://localhost:3000"
echo "  - API Health: http://localhost:3000/api/health"
echo "  - Mongo Express: http://localhost:8081 (admin/admin123)"
echo ""
echo "📋 To view logs: docker-compose -f docker-compose.dev.yml logs -f"
echo "🛑 To stop: docker-compose -f docker-compose.dev.yml down"
