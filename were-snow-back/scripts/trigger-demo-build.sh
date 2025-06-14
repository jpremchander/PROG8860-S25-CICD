#!/bin/bash

echo "🎬 Triggering Demo Build for Midterm Presentation"
echo "================================================"

# Load configuration
if [ -f ".jenkins-config" ]; then
    source .jenkins-config
else
    echo "❌ Configuration not found. Run setup first."
    exit 1
fi

JOB_NAME="SnowBored-Game-Pipeline"

echo "🚀 Triggering demo build..."
echo "🌐 Jenkins: $JENKINS_URL"
echo "🏗️ Job: $JOB_NAME"

# Get Jenkins crumb
JENKINS_CRUMB=$(curl -s "$JENKINS_URL/crumbIssuer/api/xml?xpath=concat(//crumbRequestField,\":\",//crumb)" \
  --user "$JENKINS_USER:$JENKINS_PASSWORD" \
  -H "ngrok-skip-browser-warning: true")

# Trigger build with parameters
BUILD_RESPONSE=$(curl -s -w "%{http_code}" -X POST \
  "$JENKINS_URL/job/$JOB_NAME/buildWithParameters" \
  --user "$JENKINS_USER:$JENKINS_PASSWORD" \
  -H "$JENKINS_CRUMB" \
  -H "ngrok-skip-browser-warning: true" \
  -d "BRANCH=main&DEPLOY=true")

if [ "$BUILD_RESPONSE" = "201" ]; then
    echo "✅ Demo build triggered successfully!"
    
    # Wait a moment and get build number
    sleep 3
    BUILD_NUMBER=$(curl -s --user "$JENKINS_USER:$JENKINS_PASSWORD" \
      -H "ngrok-skip-browser-warning: true" \
      "$JENKINS_URL/job/$JOB_NAME/api/json" | \
      grep -o '"nextBuildNumber":[0-9]*' | cut -d':' -f2)
    
    if [ -n "$BUILD_NUMBER" ]; then
        CURRENT_BUILD=$((BUILD_NUMBER - 1))
        echo "🔢 Build Number: #$CURRENT_BUILD"
        echo "📊 Monitor: $JENKINS_URL/job/$JOB_NAME/$CURRENT_BUILD/"
        echo "📈 Console: $JENKINS_URL/job/$JOB_NAME/$CURRENT_BUILD/console"
        echo "🎨 Blue Ocean: $JENKINS_URL/blue/organizations/jenkins/$JOB_NAME/detail/$JOB_NAME/$CURRENT_BUILD/pipeline"
    fi
    
    echo "================================================"
    echo "🎯 Demo Build Started!"
    echo "⏱️  Expected Duration: 3-5 minutes"
    echo "📋 Pipeline Stages:"
    echo "   1. ✅ Checkout - Clone repository"
    echo "   2. 🏗️ Build - Install deps, lint, build"
    echo "   3. 🧪 Test - Run unit tests with coverage"
    echo "   4. 🔒 Security - npm audit scan"
    echo "   5. 🐳 Docker - Build container image"
    echo "   6. 🚀 Deploy - Deploy to localhost:3000"
    echo "   7. 📦 Archive - Store artifacts"
    echo "================================================"
    echo "🌐 Access your app: http://localhost:3000"
    echo "📊 Jenkins Dashboard: $JENKINS_URL"
    
else
    echo "❌ Failed to trigger build (HTTP $BUILD_RESPONSE)"
fi
