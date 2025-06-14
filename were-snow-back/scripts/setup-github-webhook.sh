#!/bin/bash

echo "🔗 Setting up GitHub Webhook for Jenkins"
echo "========================================"

# Get user input
echo "Please provide the following information:"
read -p "GitHub Username: " GITHUB_USER
read -p "Repository Name: " REPO_NAME
read -s -p "GitHub Personal Access Token: " GITHUB_TOKEN
echo

# GitHub API URL
GITHUB_API="https://api.github.com"
REPO_URL="$GITHUB_API/repos/$GITHUB_USER/$REPO_NAME"

# Jenkins webhook URL (using ngrok for localhost)
echo "🌐 Setting up ngrok for webhook access..."

# Check if ngrok is installed
if ! command -v ngrok &> /dev/null; then
    echo "📥 Installing ngrok..."
    # Download and install ngrok
    wget -q https://bin.equinox.io/c/4VmDzA7iaHb/ngrok-stable-linux-amd64.zip
    unzip -q ngrok-stable-linux-amd64.zip
    sudo mv ngrok /usr/local/bin/
    rm ngrok-stable-linux-amd64.zip
fi

# Start ngrok in background
echo "🚀 Starting ngrok tunnel..."
ngrok http 8080 > /dev/null 2>&1 &
NGROK_PID=$!

# Wait for ngrok to start
sleep 5

# Get ngrok public URL
NGROK_URL=$(curl -s http://localhost:4040/api/tunnels | grep -o 'https://[^"]*\.ngrok\.io')

if [ -z "$NGROK_URL" ]; then
    echo "❌ Failed to get ngrok URL. Please check ngrok installation."
    exit 1
fi

WEBHOOK_URL="$NGROK_URL/github-webhook/"

echo "✅ Ngrok tunnel created: $NGROK_URL"
echo "🔗 Webhook URL: $WEBHOOK_URL"

# Create webhook payload
WEBHOOK_PAYLOAD=$(cat << EOF
{
  "name": "web",
  "active": true,
  "events": ["push", "pull_request"],
  "config": {
    "url": "$WEBHOOK_URL",
    "content_type": "json",
    "insecure_ssl": "0"
  }
}
EOF
)

# Create webhook
echo "📡 Creating GitHub webhook..."
RESPONSE=$(curl -s -X POST \
  -H "Authorization: token $GITHUB_TOKEN" \
  -H "Accept: application/vnd.github.v3+json" \
  -d "$WEBHOOK_PAYLOAD" \
  "$REPO_URL/hooks")

# Check if webhook was created successfully
if echo "$RESPONSE" | grep -q '"id"'; then
    echo "✅ GitHub webhook created successfully!"
    WEBHOOK_ID=$(echo "$RESPONSE" | grep -o '"id":[0-9]*' | cut -d':' -f2)
    echo "🆔 Webhook ID: $WEBHOOK_ID"
else
    echo "❌ Failed to create webhook. Response:"
    echo "$RESPONSE"
    exit 1
fi

# Save configuration
cat > .jenkins-config << EOF
GITHUB_USER=$GITHUB_USER
REPO_NAME=$REPO_NAME
NGROK_URL=$NGROK_URL
WEBHOOK_URL=$WEBHOOK_URL
WEBHOOK_ID=$WEBHOOK_ID
NGROK_PID=$NGROK_PID
EOF

echo "💾 Configuration saved to .jenkins-config"
echo "================================================="
echo "✅ Setup Complete!"
echo "🌐 Jenkins: http://localhost:8080"
echo "🔗 Public URL: $NGROK_URL"
echo "📡 Webhook: $WEBHOOK_URL"
echo "================================================="
echo "⚠️  Important: Keep this terminal open to maintain ngrok tunnel"
echo "📋 Next: Create Jenkins pipeline job"
