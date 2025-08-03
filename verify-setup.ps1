# Azure Functions CI/CD Setup Verification Script

Write-Host "=== Azure Functions CI/CD Setup Verification ===" -ForegroundColor Cyan
Write-Host ""

# Check Node.js
Write-Host "Checking Node.js installation..." -ForegroundColor Yellow
try {
    $nodeVersion = node --version
    Write-Host "✅ Node.js version: $nodeVersion" -ForegroundColor Green
} catch {
    Write-Host "❌ Node.js not found. Please install Node.js 18+" -ForegroundColor Red
    exit 1
}

# Check npm
Write-Host "Checking npm installation..." -ForegroundColor Yellow
try {
    $npmVersion = npm --version
    Write-Host "✅ npm version: $npmVersion" -ForegroundColor Green
} catch {
    Write-Host "❌ npm not found. Please install npm" -ForegroundColor Red
    exit 1
}

# Install dependencies
Write-Host ""
Write-Host "Installing dependencies..." -ForegroundColor Yellow
npm install

if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Dependencies installed successfully" -ForegroundColor Green
} else {
    Write-Host "❌ Failed to install dependencies" -ForegroundColor Red
    exit 1
}

# Run tests
Write-Host ""
Write-Host "Running tests..." -ForegroundColor Yellow
npm test

if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ All tests passed!" -ForegroundColor Green
} else {
    Write-Host "❌ Some tests failed" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "=== Setup verification completed successfully! ===" -ForegroundColor Green
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Cyan
Write-Host "1. Set up your Azure Function App"
Write-Host "2. Create a Service Principal for Jenkins"
Write-Host "3. Configure Jenkins with your GitHub repository"
Write-Host "4. Update the Jenkinsfile with your Azure details"
Write-Host "5. Run the Jenkins pipeline"
Write-Host ""
Write-Host "See AZURE_SETUP.md for detailed instructions." -ForegroundColor Yellow
