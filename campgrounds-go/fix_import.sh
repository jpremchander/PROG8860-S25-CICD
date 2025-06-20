#!/bin/bash
echo "🔧 Fixing import paths..."

# Fix import paths from campgrounds-app to yelpcamp-go
find . -name "*.go" -type f -exec sed -i 's|campgrounds-app/|yelpcamp-go/|g' {} \;

echo "✅ Import paths fixed!"
