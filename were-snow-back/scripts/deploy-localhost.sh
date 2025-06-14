#!/bin/bash

echo "🚀 Deploying SnowBored Game to Localhost"
echo "========================================"

# Build the application
echo "🏗️ Building application..."
npm run build

# Build Docker image
echo "🐳 Building Docker image..."
docker build -t snowbored-game:latest .

# Stop existing container
echo "🛑 Stopping existing container..."
docker stop snowbored-game 2>/dev/null || true
docker rm snowbored-game 2>/dev/null || true

# Run new container
echo "🚀 Starting new container..."
docker run -d \
    --name snowbored-game \
    -p 3000:3000 \
    --restart unless-stopped \
    -e NODE_ENV=production \
    snowbored-game:latest

# Wait for container to start
echo "⏳ Waiting for application to start..."
sleep 10

# Health check
if curl -f http://localhost:3000 > /dev/null 2>&1; then
    echo "✅ Deployment successful!"
    echo "🌐 Application available at: http://localhost:3000"
    echo "🐳 Container status:"
    docker ps | grep snowbored-game
else
    echo "❌ Deployment failed! Check container logs:"
    docker logs snowbored-game
    exit 1
fi
