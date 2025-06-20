#!/bin/bash

echo "🏃 Running YelpCamp Locally (without Docker)"
echo "============================================"

# Navigate to project directory
if [ -d "campgrounds-go" ]; then
    cd campgrounds-go
    echo "✅ Navigated to campgrounds-go directory"
else
    echo "❌ campgrounds-go directory not found!"
    exit 1
fi

# Check if .env exists
if [ ! -f ".env" ]; then
    echo "📝 Creating .env file for local development..."
    cat > .env << 'EOF'
# Environment Configuration
GIN_MODE=debug

# MongoDB Configuration (for local MongoDB)
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
IAM_ROLE=arn:aws:iam::054037100649:role/yelpcamp-s3-role
EOF
    echo "✅ .env file created for local development"
fi

# Download dependencies
echo "📦 Installing Go dependencies..."
go mod tidy

# Check if MongoDB is running locally
echo "🔍 Checking for local MongoDB..."
if ! pgrep mongod > /dev/null; then
    echo "⚠️  MongoDB not running locally. Starting MongoDB with Docker..."
    docker run -d --name yelpcamp-mongo \
        -p 27017:27017 \
        -e MONGO_INITDB_ROOT_USERNAME=yelpcamp_dev \
        -e MONGO_INITDB_ROOT_PASSWORD=dev_password_123 \
        -e MONGO_INITDB_DATABASE=yelpcamp_dev \
        mongo:7.0
    
    echo "⏳ Waiting for MongoDB to start..."
    sleep 10
fi

# Build and run the application
echo "🏗️  Building application..."
go build -o yelpcamp-app .

echo "🚀 Starting YelpCamp application..."
echo "📱 Application will be available at: http://localhost:3000"
echo "🛑 Press Ctrl+C to stop"
echo ""

./yelpcamp-app
