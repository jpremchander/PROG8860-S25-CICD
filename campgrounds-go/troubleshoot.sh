#!/bin/bash

echo "🔍 YelpCamp Troubleshooting Script"
echo "=================================="

# Check current directory
echo "📁 Current directory:"
pwd
ls -la

echo ""
echo "🔍 Checking for Go installation:"
go version

echo ""
echo "🔍 Checking Docker:"
docker --version
docker-compose --version

echo ""
echo "🔍 Checking running containers:"
docker ps

echo ""
echo "🔍 Checking all containers (including stopped):"
docker ps -a

echo ""
echo "🔍 Checking ports in use:"
netstat -tlnp | grep -E ':3000|:8080|:8081|:27017' || echo "No ports 3000, 8080, 8081, or 27017 in use"

echo ""
echo "🔍 Checking if campgrounds-go directory exists:"
if [ -d "campgrounds-go" ]; then
    echo "✅ campgrounds-go directory found"
    cd campgrounds-go
    echo "📁 Contents of campgrounds-go:"
    ls -la
    
    echo ""
    echo "🔍 Checking go.mod:"
    if [ -f "go.mod" ]; then
        echo "✅ go.mod found"
        head -10 go.mod
    else
        echo "❌ go.mod not found"
    fi
    
    echo ""
    echo "🔍 Checking .env file:"
    if [ -f ".env" ]; then
        echo "✅ .env file found"
        echo "Environment variables (without sensitive data):"
        grep -E "^[A-Z_]+" .env | grep -v -E "(PASSWORD|SECRET|KEY)" || echo "No environment variables found"
    else
        echo "❌ .env file not found"
    fi
    
    echo ""
    echo "🔍 Checking main.go:"
    if [ -f "main.go" ]; then
        echo "✅ main.go found"
    else
        echo "❌ main.go not found"
    fi
    
else
    echo "❌ campgrounds-go directory not found"
    echo "Available directories:"
    ls -la
fi

echo ""
echo "🔍 Checking Docker Compose files:"
find . -name "docker-compose*.yml" -type f

echo ""
echo "=================================="
echo "🎯 Next Steps:"
echo "1. Navigate to campgrounds-go directory: cd campgrounds-go"
echo "2. Check if .env file exists and has correct values"
echo "3. Try running: go mod tidy"
echo "4. Try running: docker-compose -f docker-compose.dev.yml up --build"
echo "5. Or try running directly: go run ."
