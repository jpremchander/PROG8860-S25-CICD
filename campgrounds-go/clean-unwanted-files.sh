#!/bin/bash

echo "🧹 Cleaning up unwanted Next.js/React files from Go project..."

# Remove Next.js/React directories
echo "Removing Next.js/React directories..."
rm -rf app/
rm -rf components/
rm -rf hooks/
rm -rf lib/
rm -rf styles/

# Remove Next.js/React config files
echo "Removing Next.js/React config files..."
rm -f next.config.mjs
rm -f package.json
rm -f pnpm-lock.yaml
rm -f postcss.config.mjs
rm -f tailwind.config.ts
rm -f tsconfig.json
rm -f components.json

# Remove duplicate files
echo "Removing duplicate files..."
rm -f SETUP_GUIDE.md  # Keep SETUP-GUIDE.md
rm -f Jenkinsfile-OLD
rm -f Jenkinsfile-v2

# Keep only essential public assets, remove the rest
echo "Cleaning up public directory..."
if [ -d "public" ]; then
    # Keep only placeholder images that might be used by Go templates
    mkdir -p temp_public
    cp public/placeholder*.* temp_public/ 2>/dev/null || true
    rm -rf public/
    mv temp_public public/
fi

# Remove any build artifacts
echo "Removing build artifacts..."
rm -f yelpcamp-go
rm -f yelpcamp-app

echo "✅ Cleanup completed!"
echo ""
echo "📋 Remaining Go project structure:"
echo "├── Go source files (*.go)"
echo "├── Templates (HTML)"
echo "├── Static assets (CSS, images)"
echo "├── Docker configuration"
echo "├── Database scripts"
echo "└── Shell scripts for deployment"
echo ""
echo "🚀 Your Go project is now clean and ready!"
