#!/bin/bash

echo "=== Azure Functions CI/CD Setup Verification ==="
echo ""

# Check Node.js
echo "Checking Node.js installation..."
if command -v node &> /dev/null; then
    echo "✅ Node.js version: $(node --version)"
else
    echo "❌ Node.js not found. Please install Node.js 18+"
    exit 1
fi

# Check npm
echo "Checking npm installation..."
if command -v npm &> /dev/null; then
    echo "✅ npm version: $(npm --version)"
else
    echo "❌ npm not found. Please install npm"
    exit 1
fi

# Install dependencies
echo ""
echo "Installing dependencies..."
npm install

if [ $? -eq 0 ]; then
    echo "✅ Dependencies installed successfully"
else
    echo "❌ Failed to install dependencies"
    exit 1
fi

# Run tests
echo ""
echo "Running tests..."
npm test

if [ $? -eq 0 ]; then
    echo "✅ All tests passed!"
else
    echo "❌ Some tests failed"
    exit 1
fi

echo ""
echo "=== Setup verification completed successfully! ==="
echo ""
echo "Next steps:"
echo "1. Set up your Azure Function App"
echo "2. Create a Service Principal for Jenkins"
echo "3. Configure Jenkins with your GitHub repository"
echo "4. Update the Jenkinsfile with your Azure details"
echo "5. Run the Jenkins pipeline"
echo ""
echo "See AZURE_SETUP.md for detailed instructions."
