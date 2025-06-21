#!/bin/bash

echo "🔍 YelpCamp Troubleshooting Script"
echo "=================================="

# Check if containers are running
echo "📦 Checking Docker containers..."
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

echo ""
echo "🗄️ Checking MongoDB connection..."
docker exec yelpcamp-mongo mongosh yelpcamp_dev --eval "db.runCommand('ping')" 2>/dev/null || echo "❌ MongoDB connection failed"

echo ""
echo "📊 Checking database collections..."
docker exec yelpcamp-mongo mongosh yelpcamp_dev --eval "
  print('Campgrounds count:', db.campgrounds.countDocuments({}));
  print('Users count:', db.users.countDocuments({}));
  print('Reviews count:', db.reviews.countDocuments({}));
" 2>/dev/null || echo "❌ Database query failed"

echo ""
echo "🏕️ Sample campgrounds in database..."
docker exec yelpcamp-mongo mongosh yelpcamp_dev --eval "
  db.campgrounds.find({}, {title: 1, location: 1, price: 1}).limit(3).forEach(printjson);
" 2>/dev/null || echo "❌ Campgrounds query failed"

echo ""
echo "🌐 Testing API endpoints..."
echo "Health check:"
curl -s http://localhost:3000/api/health | jq '.' 2>/dev/null || echo "❌ Health endpoint failed"

echo ""
echo "Campgrounds API:"
curl -s http://localhost:3000/api/campgrounds | jq '.count, .campgrounds[0].title' 2>/dev/null || echo "❌ Campgrounds API failed"

echo ""
echo "Debug DB endpoint:"
curl -s http://localhost:3000/debug/db | jq '.database_stats' 2>/dev/null || echo "❌ Debug endpoint failed"

echo ""
echo "✅ Troubleshooting complete!"
echo "Visit http://localhost:3000/campgrounds to see the web interface"
echo "Visit http://localhost:3000/debug/db to see database contents"
