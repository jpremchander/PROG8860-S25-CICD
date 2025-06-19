#!/bin/bash

echo "⚙️ Configuring Existing Jenkins for SnowBored Game"
echo "================================================"

# Get Jenkins details
read -p "Jenkins URL (default: http://localhost:8080): " JENKINS_URL
JENKINS_URL=${JENKINS_URL:-http://localhost:8080}

read -p "Jenkins Username: " JENKINS_USER
read -s -p "Jenkins Password/Token: " JENKINS_PASSWORD
echo

# Test Jenkins connection
echo "🔗 Testing Jenkins connection..."
JENKINS_CRUMB=$(curl -s "$JENKINS_URL/crumbIssuer/api/xml?xpath=concat(//crumbRequestField,\":\",//crumb)" \
  --user "$JENKINS_USER:$JENKINS_PASSWORD")

if [ -z "$JENKINS_CRUMB" ]; then
    echo "❌ Failed to connect to Jenkins. Please check credentials."
    exit 1
fi

echo "✅ Connected to Jenkins successfully!"

# Download Jenkins CLI if not exists
if [ ! -f "jenkins-cli.jar" ]; then
    echo "📥 Downloading Jenkins CLI..."
    wget -q "$JENKINS_URL/jnlpJars/jenkins-cli.jar"
fi

# Create Jenkins CLI alias
alias jenkins-cli="java -jar jenkins-cli.jar -s $JENKINS_URL -auth $JENKINS_USER:$JENKINS_PASSWORD"

# Check and install required plugins
echo "🔌 Checking required plugins..."
REQUIRED_PLUGINS=(
    "git"
    "github"
    "pipeline-stage-view"
    "docker-workflow"
    "nodejs"
    "html-publisher"
    "email-ext"
    "build-timeout"
    "timestamper"
    "ws-cleanup"
    "github-branch-source"
    "blue-ocean"
)

INSTALLED_PLUGINS=$(jenkins-cli list-plugins | awk '{print $1}')

for plugin in "${REQUIRED_PLUGINS[@]}"; do
    if echo "$INSTALLED_PLUGINS" | grep -q "^$plugin$"; then
        echo "✅ $plugin already installed"
    else
        echo "📦 Installing $plugin..."
        jenkins-cli install-plugin "$plugin"
        RESTART_NEEDED=true
    fi
done

if [ "$RESTART_NEEDED" = true ]; then
    echo "🔄 Restarting Jenkins to activate new plugins..."
    jenkins-cli restart
    echo "⏳ Waiting for Jenkins to restart..."
    sleep 30
fi

# Configure Global Tools
echo "🛠️ Configuring Global Tools..."

# Check if NodeJS is configured
NODE_CONFIG=$(jenkins-cli get-node localhost | grep -i nodejs || echo "")
if [ -z "$NODE_CONFIG" ]; then
    echo "📦 Configuring NodeJS..."
    # This will be done via Jenkins UI or API
    echo "⚠️  Please configure NodeJS manually in Global Tool Configuration"
fi

echo "✅ Jenkins configuration completed!"
echo "📋 Next steps:"
echo "1. Configure Global Tools (if needed)"
echo "2. Set up GitHub webhook"
echo "3. Create pipeline job"
