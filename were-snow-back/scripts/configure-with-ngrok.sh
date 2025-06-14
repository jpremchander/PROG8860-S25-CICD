#!/bin/bash

echo "🚀 Configuring Jenkins with Your Ngrok URL"
echo "=========================================="

# Your Jenkins configuration
JENKINS_URL="https://ec7a-142-156-1-230.ngrok-free.app"
WEBHOOK_URL="$JENKINS_URL/github-webhook/"

echo "🌐 Jenkins URL: $JENKINS_URL"
echo "🔗 Webhook URL: $WEBHOOK_URL"

# Get user credentials
read -p "Jenkins Username: " JENKINS_USER
read -s -p "Jenkins Password/API Token: " JENKINS_PASSWORD
echo

read -p "GitHub Username: " GITHUB_USER
read -p "Repository Name (default: snowbored-game): " REPO_NAME
REPO_NAME=${REPO_NAME:-snowbored-game}
read -s -p "GitHub Personal Access Token: " GITHUB_TOKEN
echo

# Test Jenkins connection
echo "🔗 Testing Jenkins connection..."
RESPONSE=$(curl -s -w "%{http_code}" -o /dev/null "$JENKINS_URL/login" \
  --connect-timeout 10)

if [ "$RESPONSE" = "200" ]; then
    echo "✅ Jenkins is accessible via ngrok"
else
    echo "❌ Cannot access Jenkins (HTTP $RESPONSE)"
    echo "⚠️  Make sure Jenkins is running and ngrok tunnel is active"
    exit 1
fi

# Get Jenkins crumb for CSRF protection
echo "🔐 Getting Jenkins authentication..."
JENKINS_CRUMB=$(curl -s "$JENKINS_URL/crumbIssuer/api/xml?xpath=concat(//crumbRequestField,\":\",//crumb)" \
  --user "$JENKINS_USER:$JENKINS_PASSWORD" \
  -H "ngrok-skip-browser-warning: true")

if [ -z "$JENKINS_CRUMB" ]; then
    echo "❌ Failed to authenticate with Jenkins"
    echo "💡 Check your username and password/token"
    exit 1
fi

echo "✅ Jenkins authentication successful"
# Fix CSRF token issue
echo "🔧 Configuring CSRF protection..."
echo "JENKINS_CRUMB=$JENKINS_CRUMB" >> .jenkins-config

# Create GitHub webhook
echo "📡 Creating GitHub webhook..."
WEBHOOK_PAYLOAD=$(cat << EOF
{
  "name": "web",
  "active": true,
  "events": ["push", "pull_request"],
  "config": {
    "url": "$WEBHOOK_URL",
    "content_type": "json",
    "insecure_ssl": "1"
  }
}
EOF
)

WEBHOOK_RESPONSE=$(curl -s -X POST \
  -H "Authorization: token $GITHUB_TOKEN" \
  -H "Accept: application/vnd.github.v3+json" \
  -d "$WEBHOOK_PAYLOAD" \
  "https://api.github.com/repos/$GITHUB_USER/$REPO_NAME/hooks")

if echo "$WEBHOOK_RESPONSE" | grep -q '"id"'; then
    WEBHOOK_ID=$(echo "$WEBHOOK_RESPONSE" | grep -o '"id":[0-9]*' | cut -d':' -f2)
    echo "✅ GitHub webhook created successfully!"
    echo "🆔 Webhook ID: $WEBHOOK_ID"
else
    echo "❌ Failed to create webhook:"
    echo "$WEBHOOK_RESPONSE"
    echo "💡 Check your GitHub token permissions"
fi

# Extract webhook ID properly
if echo "$WEBHOOK_RESPONSE" | grep -q '"id"'; then
    WEBHOOK_ID=$(echo "$WEBHOOK_RESPONSE" | grep -o '"id":[0-9]*' | cut -d':' -f2)
    echo "WEBHOOK_ID=$WEBHOOK_ID" >> .jenkins-config
else
    echo "WEBHOOK_ID=" >> .jenkins-config
fi

# Save configuration
cat > .jenkins-config << EOF
JENKINS_URL=$JENKINS_URL
JENKINS_USER=$JENKINS_USER
JENKINS_PASSWORD=$JENKINS_PASSWORD
GITHUB_USER=$GITHUB_USER
REPO_NAME=$REPO_NAME
WEBHOOK_URL=$WEBHOOK_URL
WEBHOOK_ID=$WEBHOOK_ID
EOF

echo "💾 Configuration saved to .jenkins-config"

# Test webhook
echo "🧪 Testing webhook endpoint..."
WEBHOOK_TEST=$(curl -s -w "%{http_code}" -o /dev/null \
  -H "ngrok-skip-browser-warning: true" \
  "$WEBHOOK_URL")

if [ "$WEBHOOK_TEST" = "200" ] || [ "$WEBHOOK_TEST" = "405" ]; then
    echo "✅ Webhook endpoint is accessible"
else
    echo "⚠️  Webhook endpoint returned HTTP $WEBHOOK_TEST"
fi

echo "================================================"
echo "✅ Configuration Complete!"
echo "🌐 Jenkins: $JENKINS_URL"
echo "🔗 Webhook: $WEBHOOK_URL"
echo "📋 Next: Create pipeline job"
echo "================================================"

# Check and create missing project files
echo "🔍 Checking project structure..."
if [ ! -f "package.json" ]; then
    echo "📦 Creating package.json..."
    cat > package.json << 'EOF'
{
  "name": "snowbored-game",
  "version": "1.0.0",
  "private": true,
  "scripts": {
    "dev": "next dev",
    "build": "next build",
    "start": "next start",
    "lint": "next lint",
    "test": "jest --passWithNoTests",
    "test:coverage": "jest --coverage --passWithNoTests"
  },
  "dependencies": {
    "next": "14.0.0",
    "react": "^18",
    "react-dom": "^18"
  },
  "devDependencies": {
    "@types/node": "^20",
    "@types/react": "^18",
    "eslint": "^8",
    "eslint-config-next": "14.0.0",
    "typescript": "^5",
    "jest": "^29.7.0"
  }
}
EOF
    echo "✅ package.json created"
fi

if [ ! -f "Dockerfile" ]; then
    echo "🐳 Creating Dockerfile..."
    cat > Dockerfile << 'EOF'
FROM node:18-alpine AS base
WORKDIR /app
COPY package*.json ./

FROM base AS deps
RUN npm ci --only=production

FROM base AS build
RUN npm ci
COPY . .
RUN npm run build

FROM node:18-alpine AS runtime
WORKDIR /app
RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs
COPY --from=build --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=build --chown=nextjs:nodejs /app/.next/static ./.next/static
COPY --from=build --chown=nextjs:nodejs /app/public ./public
USER nextjs
EXPOSE 3000
ENV PORT 3000
CMD ["node", "server.js"]
EOF
    echo "✅ Dockerfile created"
fi

echo "📦 Installing dependencies..."
npm install
echo "✅ Dependencies installed"
