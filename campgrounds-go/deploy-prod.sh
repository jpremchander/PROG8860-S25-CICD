#!/bin/bash
# Production deployment script

set -e

echo "🚀 Deploying YelpCamp Go - PRODUCTION Environment"

# Build production image
docker build --target production -t yelp-camp:prod .

# Stop existing prod container
docker stop campgrounds-prod 2>/dev/null || true
docker rm campgrounds-prod 2>/dev/null || true

# Start production container
docker run -d \
    --name campgrounds-prod \
    -p 3000:3000 \
    -e PORT=3000 \
    -e GIN_MODE=release \
    --restart unless-stopped \
    yelp-camp:prod

echo "⏳ Waiting for production server to start..."
sleep 10

# Health check
if curl -f http://localhost:3000/health; then
    echo "✅ Production deployment successful!"
    echo "🌐 Application: http://localhost:3000"
    echo "🔒 Optimized for production performance"
else
    echo "❌ Production deployment failed!"
    docker logs campgrounds-prod
    exit 1
fi
