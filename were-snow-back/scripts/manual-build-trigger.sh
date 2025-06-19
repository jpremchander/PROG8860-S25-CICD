#!/bin/bash

echo "🚀 Manual Build Trigger with Proper Authentication"
echo "================================================="

# Load configuration
if [ -f ".jenkins-config" ]; then
    source .jenkins-config
else
    echo "❌ Configuration not found"
    exit 1
fi

echo "🔐 Getting fresh CSRF token..."
JENKINS_CRUMB=$(curl -s "$JENKINS_URL/crumbIssuer/api/xml?xpath=concat(//crumbRequestField,\":\",//crumb)" \
  --user "$JENKINS_USER:$JENKINS_PASSWORD" \
  -H "ngrok-skip-browser-warning: true")

if [ -z "$JENKINS_CRUMB" ]; then
    echo "❌ Failed to get CSRF token"
    exit 1
fi

echo "✅ CSRF token obtained"

# Check job exists
echo "🔍 Verifying job exists..."
JOB_CHECK=$(curl -s -w "%{http_code}" -o /dev/null \
  --user "$JENKINS_USER:$JENKINS_PASSWORD" \
  -H "ngrok-skip-browser-warning: true" \
  "$JENKINS_URL/job/SnowBored-Game-Pipeline/api/json")

if [ "$JOB_CHECK" != "200" ]; then
    echo "❌ Job not accessible (HTTP $JOB_CHECK)"
    echo "🌐 Please check: $JENKINS_URL/job/SnowBored-Game-Pipeline/"
    exit 1
fi

echo "✅ Job verified"

# Trigger build
echo "🚀 Triggering build..."
BUILD_RESPONSE=$(curl -s -w "%{http_code}" -X POST \
  "$JENKINS_URL/job/SnowBored-Game-Pipeline/build" \
  --user "$JENKINS_USER:$JENKINS_PASSWORD" \
  -H "$JENKINS_CRUMB" \
  -H "ngrok-skip-browser-warning: true")

echo "📊 Build response: $BUILD_RESPONSE"

if [[ "$BUILD_RESPONSE" == *"201"* ]] || [[ "$BUILD_RESPONSE" == *"200"* ]]; then
    echo "✅ Build triggered successfully!"
    
    # Wait and get build number
    sleep 3
    echo "🔢 Getting build number..."
    
    BUILD_INFO=$(curl -s --user "$JENKINS_USER:$JENKINS_PASSWORD" \
      -H "ngrok-skip-browser-warning: true" \
      "$JENKINS_URL/job/SnowBored-Game-Pipeline/api/json")
    
    LAST_BUILD=$(echo "$BUILD_INFO" | grep -o '"lastBuild":{"number":[0-9]*' | grep -o '[0-9]*')
    
    if [ -n "$LAST_BUILD" ]; then
        echo "🎯 Build #$LAST_BUILD started"
        echo "📊 Monitor: $JENKINS_URL/job/SnowBored-Game-Pipeline/$LAST_BUILD/"
        echo "📈 Console: $JENKINS_URL/job/SnowBored-Game-Pipeline/$LAST_BUILD/console"
        echo "🎨 Blue Ocean: $JENKINS_URL/blue/organizations/jenkins/SnowBored-Game-Pipeline/detail/SnowBored-Game-Pipeline/$LAST_BUILD/pipeline"
    fi
    
else
    echo "❌ Build trigger failed"
    echo "Response: $BUILD_RESPONSE"
fi

echo "================================================"
echo "🌐 Jenkins Dashboard: $JENKINS_URL"
echo "🏗️ Pipeline Job: $JENKINS_URL/job/SnowBored-Game-Pipeline/"
echo "================================================"
