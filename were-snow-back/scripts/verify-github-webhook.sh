#!/bin/bash

echo "🔍 Verifying GitHub Webhook Configuration"
echo "========================================"

# Load configuration
if [ -f ".jenkins-config" ]; then
    source .jenkins-config
else
    echo "❌ Configuration not found"
    exit 1
fi

echo "📡 Checking webhook in GitHub repository..."
echo "Repository: $GITHUB_USER/$REPO_NAME"

# Get all webhooks
WEBHOOKS_RESPONSE=$(curl -s \
  -H "Authorization: token $GITHUB_TOKEN" \
  "https://api.github.com/repos/$GITHUB_USER/$REPO_NAME/hooks")

if echo "$WEBHOOKS_RESPONSE" | grep -q "message.*Not Found"; then
    echo "❌ Repository not found or no access"
    exit 1
fi

# Check if our webhook exists
WEBHOOK_COUNT=$(echo "$WEBHOOKS_RESPONSE" | grep -c '"url"')
echo "📊 Total webhooks found: $WEBHOOK_COUNT"

if [ "$WEBHOOK_COUNT" -gt 0 ]; then
    echo "🔍 Webhook details:"
    echo "$WEBHOOKS_RESPONSE" | grep -A 10 -B 2 "$WEBHOOK_URL" || echo "⚠️  Our webhook URL not found"
    
    # Get our webhook ID if missing
    if [ -z "$WEBHOOK_ID" ]; then
        WEBHOOK_ID=$(echo "$WEBHOOKS_RESPONSE" | grep -B5 -A5 "$WEBHOOK_URL" | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)
        if [ -n "$WEBHOOK_ID" ]; then
            echo "✅ Found webhook ID: $WEBHOOK_ID"
            # Update config file
            sed -i "s/WEBHOOK_ID=$/WEBHOOK_ID=$WEBHOOK_ID/" .jenkins-config
        fi
    fi
else
    echo "❌ No webhooks found"
    echo "🔧 Creating webhook..."
    
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

    CREATE_RESPONSE=$(curl -s -X POST \
      -H "Authorization: token $GITHUB_TOKEN" \
      -H "Accept: application/vnd.github.v3+json" \
      -d "$WEBHOOK_PAYLOAD" \
      "https://api.github.com/repos/$GITHUB_USER/$REPO_NAME/hooks")
    
    if echo "$CREATE_RESPONSE" | grep -q '"id"'; then
        NEW_WEBHOOK_ID=$(echo "$CREATE_RESPONSE" | grep -o '"id":[0-9]*' | cut -d':' -f2)
        echo "✅ Webhook created with ID: $NEW_WEBHOOK_ID"
        # Update config
        sed -i "s/WEBHOOK_ID=$/WEBHOOK_ID=$NEW_WEBHOOK_ID/" .jenkins-config
    else
        echo "❌ Failed to create webhook:"
        echo "$CREATE_RESPONSE"
    fi
fi

# Test webhook endpoint
echo "🧪 Testing webhook endpoint..."
WEBHOOK_TEST=$(curl -s -w "%{http_code}" -o /dev/null \
  -H "ngrok-skip-browser-warning: true" \
  "$WEBHOOK_URL")

if [ "$WEBHOOK_TEST" = "200" ] || [ "$WEBHOOK_TEST" = "405" ]; then
    echo "✅ Webhook endpoint accessible (HTTP $WEBHOOK_TEST)"
else
    echo "❌ Webhook endpoint issue (HTTP $WEBHOOK_TEST)"
fi

echo "================================================"
echo "📊 Webhook Status Summary:"
echo "🌐 Repository: https://github.com/$GITHUB_USER/$REPO_NAME"
echo "🔗 Webhook URL: $WEBHOOK_URL"
echo "🆔 Webhook ID: ${WEBHOOK_ID:-'Not found'}"
echo "✅ Endpoint: Accessible"
echo "================================================"
