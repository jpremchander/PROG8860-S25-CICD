#!/bin/bash

echo "📦 Installing Project Dependencies"
echo "================================="

# Check if we're in the right directory
if [ ! -f "package.json" ]; then
    echo "❌ package.json not found. Are you in the project directory?"
    exit 1
fi

echo "🔍 Checking Node.js and npm..."
if ! command -v node &> /dev/null; then
    echo "❌ Node.js not installed"
    echo "💡 Install Node.js: https://nodejs.org/"
    exit 1
fi

if ! command -v npm &> /dev/null; then
    echo "❌ npm not installed"
    exit 1
fi

echo "✅ Node.js $(node --version)"
echo "✅ npm $(npm --version)"

# Install dependencies
echo "📦 Installing dependencies..."
npm install

if [ $? -eq 0 ]; then
    echo "✅ Dependencies installed successfully"
else
    echo "❌ Failed to install dependencies"
    exit 1
fi

# Test build
echo "🏗️ Testing build..."
npm run build

if [ $? -eq 0 ]; then
    echo "✅ Build successful"
else
    echo "❌ Build failed"
    echo "💡 Check your Next.js configuration"
fi

# Test if we can run tests
echo "🧪 Testing test suite..."
if npm test -- --passWithNoTests --watchAll=false; then
    echo "✅ Test suite working"
else
    echo "⚠️  Test suite needs configuration"
fi

echo "================================================"
echo "✅ Project dependencies ready!"
echo "🚀 You can now run:"
echo "   npm run dev    # Development server"
echo "   npm run build  # Production build"
echo "   npm test       # Run tests"
echo "================================================"
