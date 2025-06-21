#!/bin/bash

echo "🔍 Verifying Go project structure..."
echo ""

# Check Go files
echo "📁 Go source files:"
find . -name "*.go" -not -path "./vendor/*" | head -10

echo ""
echo "📁 Key directories:"
ls -la | grep "^d" | grep -E "(config|controllers|middleware|models|routes|templates|static|utils)"

echo ""
echo "📄 Configuration files:"
ls -la | grep -E "(go.mod|go.sum|Dockerfile|docker-compose)"

echo ""
echo "🧪 Testing Go build..."
if go mod tidy && go build -o test-build .; then
    echo "✅ Go project builds successfully!"
    rm -f test-build
else
    echo "❌ Go build failed - check for missing dependencies"
fi
