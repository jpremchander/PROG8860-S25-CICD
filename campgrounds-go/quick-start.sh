#!/bin/bash

echo "🔧 Quick Fix for YelpCamp Campgrounds Issue"
echo "==========================================="

# Make scripts executable
chmod +x debug-database.sh
chmod +x force-seed-data.sh
chmod +x restart-with-fresh-data.sh

echo "1. Running database debug..."
./debug-database.sh

echo ""
echo "2. Force seeding data..."
./force-seed-data.sh

echo ""
echo "3. Testing endpoints..."
sleep 3

echo "   API campgrounds:"
curl -s http://localhost:3000/api/campgrounds | jq '.count' 2>/dev/null || echo "   Could not fetch count"

echo ""
echo "   Sample campground titles:"
curl -s http://localhost:3000/api/campgrounds | jq '.data[].title' 2>/dev/null || echo "   Could not fetch titles"

echo ""
echo "✅ Quick fix completed!"
echo ""
echo "🌐 Now visit: http://localhost:3000/campgrounds"
echo "📱 You should see 8 sample campgrounds with images"
