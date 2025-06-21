#!/bin/bash
# Development deployment script

set -e

echo "🚀 Deploying YelpCamp Go - DEVELOPMENT Environment"

# Build development image
docker build --target development -t yelp-camp:dev .

# Stop existing dev container
docker stop campgrounds-dev 2>/dev/null || true
docker rm campgrounds-dev 2>/dev/null || true

# Start development container with hot reload
docker run -d \
    --name campgrounds-dev \
    -p 3001:3001 \
    -e PORT=3001 \
    -e GIN_MODE=debug \
    -v $(pwd):/app \
    --restart unless-stopped \
    yelp-camp:dev

echo "⏳ Waiting for development server to start..."
sleep 10

# Health check
if curl -f http://localhost:3001/health; then
    echo "✅ Development deployment successful!"
    echo "🌐 Application: http://localhost:3001"
    echo "🔄 Hot reload enabled - changes will auto-reload"
else
    echo "❌ Development deployment failed!"
    docker logs campgrounds-dev
    exit 1
fi
