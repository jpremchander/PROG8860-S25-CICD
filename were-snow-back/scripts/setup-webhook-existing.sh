#!/bin/bash

echo "🔗 Setting up GitHub Webhook for Existing Jenkins"
echo "================================================"

# Get configuration
read -p "GitHub Username: " GITHUB_USER
read -p "Repository Name: " REPO_NAME
read -s -p "GitHub Personal Access Token: " GITHUB_TOKEN
echo

read -p "Jenkins URL (default: http://localhost:8080): " JENKINS_URL
JENKINS_URL=${JENKINS_URL:-http://localhost:8080}

# Check if Jenkins is accessible from internet
echo "🌐 Checking Jenkins accessibility..."
if curl -s --connect-timeout 5 "$JENKINS_URL" > /dev/null; then
    echo "✅ Jenkins is accessible"
    WEBHOOK_URL="$JENKINS_URL/github-webhook/"
    USE_NGROK=false
else
    echo "🔒 Jenkins not accessible from internet. Setting up ngrok..."
    USE_NGROK=true
fi

if [ "$USE_NGROK" = true ]; then
    # Check if ngrok is installed
    if ! command -v ngrok &> /dev/null; then
        echo "📥 Installing ngrok..."
        # For Linux
        if [[ "$OSTYPE" == "linux-gnu"* ]]; then
            wget -q https://bin.equinox.io/c/4VmDzA7iaHb/ngrok-stable-linux-amd64.zip
            unzip -q ngrok-stable-linux-amd64.zip
            sudo mv ngrok /usr/local/bin/
            rm ngrok-stable-linux-amd64.zip
        # For macOS
        elif [[ "$OSTYPE" == "darwin"* ]]; then
            brew install ngrok
        fi
    fi

    # Start ngrok
    echo "🚀 Starting ngrok tunnel..."
    ngrok http 8080 > /dev/null 2>&1 &
    NGROK_PID=$!
    sleep 5

    # Get ngrok URL
    NGROK_URL=$(curl -s http://localhost:4040/api/tunnels | grep -o 'https://[^"]*\.ngrok\.io')
    
    if [ -z "$NGROK_URL" ]; then
        echo "❌ Failed to get ngrok URL"
        exit 1
    fi

    WEBHOOK_URL="$NGROK_URL/github-webhook/"
    echo "✅ Ngrok tunnel: $NGROK_URL"
fi

echo "🔗 Webhook URL: $WEBHOOK_URL"

# Create webhook
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

echo "📡 Creating GitHub webhook..."
RESPONSE=$(curl -s -X POST \
  -H "Authorization: token $GITHUB_TOKEN" \
  -H "Accept: application/vnd.github.v3+json" \
  -d "$WEBHOOK_PAYLOAD" \
  "https://api.github.com/repos/$GITHUB_USER/$REPO_NAME/hooks")

if echo "$RESPONSE" | grep -q '"id"'; then
    WEBHOOK_ID=$(echo "$RESPONSE" | grep -o '"id":[0-9]*' | cut -d':' -f2)
    echo "✅ GitHub webhook created successfully!"
    echo "🆔 Webhook ID: $WEBHOOK_ID"
    
    # Save configuration
    cat > .jenkins-webhook-config << EOF
GITHUB_USER=$GITHUB_USER
REPO_NAME=$REPO_NAME
JENKINS_URL=$JENKINS_URL
WEBHOOK_URL=$WEBHOOK_URL
WEBHOOK_ID=$WEBHOOK_ID
USE_NGROK=$USE_NGROK
EOF

    if [ "$USE_NGROK" = true ]; then
        echo "NGROK_URL=$NGROK_URL" >> .jenkins-webhook-config
        echo "NGROK_PID=$NGROK_PID" >> .jenkins-webhook-config
        echo "⚠️  Keep terminal open to maintain ngrok tunnel"
    fi

    echo "💾 Configuration saved to .jenkins-webhook-config"
else
    echo "❌ Failed to create webhook:"
    echo "$RESPONSE"
    exit 1
fi

echo "✅ Webhook setup completed!"
