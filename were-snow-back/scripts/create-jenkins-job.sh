#!/bin/bash

echo "🏗️ Creating Jenkins Pipeline Job"
echo "================================"

# Get configuration
read -p "Jenkins URL (default: http://localhost:8080): " JENKINS_URL
JENKINS_URL=${JENKINS_URL:-http://localhost:8080}

read -p "Jenkins Username: " JENKINS_USER
read -s -p "Jenkins Password/Token: " JENKINS_PASSWORD
echo

read -p "GitHub Username: " GITHUB_USER
read -p "Repository Name (default: snowbored-game): " REPO_NAME
REPO_NAME=${REPO_NAME:-snowbored-game}

JOB_NAME="SnowBored-Game-Pipeline"

# Get Jenkins crumb for CSRF protection
JENKINS_CRUMB=$(curl -s "$JENKINS_URL/crumbIssuer/api/xml?xpath=concat(//crumbRequestField,\":\",//crumb)" \
  --user "$JENKINS_USER:$JENKINS_PASSWORD")

if [ -z "$JENKINS_CRUMB" ]; then
    echo "❌ Failed to get Jenkins crumb. Check credentials."
    exit 1
fi

# Update job config with actual GitHub URL
sed "s/YOUR_USERNAME/$GITHUB_USER/g" jenkins/job-config.xml > /tmp/job-config.xml

# Create the job
echo "🏗️ Creating pipeline job..."
RESPONSE=$(curl -s -X POST "$JENKINS_URL/createItem?name=$JOB_NAME" \
  --user "$JENKINS_USER:$JENKINS_PASSWORD" \
  -H "$JENKINS_CRUMB" \
  -H "Content-Type: application/xml" \
  --data-binary @/tmp/job-config.xml)

if [ $? -eq 0 ]; then
    echo "✅ Pipeline job '$JOB_NAME' created successfully!"
    echo "🌐 Access at: $JENKINS_URL/job/$JOB_NAME/"
    
    # Trigger initial build
    echo "🚀 Triggering initial build..."
    curl -s -X POST "$JENKINS_URL/job/$JOB_NAME/build" \
      --user "$JENKINS_USER:$JENKINS_PASSWORD" \
      -H "$JENKINS_CRUMB"
    
    echo "✅ Initial build triggered!"
else
    echo "❌ Failed to create job. Response: $RESPONSE"
    exit 1
fi

# Clean up
rm /tmp/job-config.xml

echo "================================================"
echo "✅ Setup Complete!"
echo "🌐 Jenkins Job: $JENKINS_URL/job/$JOB_NAME/"
echo "📋 Next: Push code to trigger automatic builds"
