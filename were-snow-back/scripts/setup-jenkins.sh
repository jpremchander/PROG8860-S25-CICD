#!/bin/bash

echo "🚀 Setting up Jenkins for SnowBored Game CI/CD Pipeline"
echo "=================================================="

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    exit 1
fi

# Create Jenkins directory structure
echo "📁 Creating Jenkins directory structure..."
mkdir -p jenkins_home
mkdir -p jenkins_data/jobs
mkdir -p jenkins_data/plugins
mkdir -p jenkins_data/secrets

# Set proper permissions
sudo chown -R 1000:1000 jenkins_home
sudo chmod -R 755 jenkins_home

echo "✅ Directory structure created successfully!"

# Create Jenkins Docker Compose file
cat > docker-compose.jenkins.yml << 'EOF'
version: '3.8'

services:
  jenkins:
    image: jenkins/jenkins:lts
    container_name: jenkins-snowbored
    restart: unless-stopped
    ports:
      - "8080:8080"
      - "50000:50000"
    volumes:
      - ./jenkins_home:/var/jenkins_home
      - /var/run/docker.sock:/var/run/docker.sock
      - /usr/bin/docker:/usr/bin/docker
    environment:
      - JENKINS_OPTS=--httpPort=8080
      - JAVA_OPTS=-Djenkins.install.runSetupWizard=false
    user: root
    networks:
      - jenkins-network

  jenkins-agent:
    image: jenkins/inbound-agent:latest
    container_name: jenkins-agent-snowbored
    restart: unless-stopped
    environment:
      - JENKINS_URL=http://jenkins:8080
      - JENKINS_AGENT_NAME=docker-agent
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - /usr/bin/docker:/usr/bin/docker
    networks:
      - jenkins-network
    depends_on:
      - jenkins

networks:
  jenkins-network:
    driver: bridge
EOF

echo "✅ Docker Compose configuration created!"

# Start Jenkins
echo "🚀 Starting Jenkins..."
docker-compose -f docker-compose.jenkins.yml up -d

echo "⏳ Waiting for Jenkins to start (this may take 2-3 minutes)..."
sleep 60

# Get initial admin password
echo "🔑 Getting Jenkins initial admin password..."
JENKINS_PASSWORD=$(docker exec jenkins-snowbored cat /var/jenkins_home/secrets/initialAdminPassword 2>/dev/null)

if [ -n "$JENKINS_PASSWORD" ]; then
    echo "✅ Jenkins is starting up!"
    echo "=================================================="
    echo "🌐 Jenkins URL: http://localhost:8080"
    echo "🔑 Initial Admin Password: $JENKINS_PASSWORD"
    echo "=================================================="
    echo "📋 Next Steps:"
    echo "1. Open http://localhost:8080 in your browser"
    echo "2. Enter the password above"
    echo "3. Install suggested plugins"
    echo "4. Create your admin user"
    echo "5. Run the configuration script: ./scripts/configure-jenkins.sh"
else
    echo "⏳ Jenkins is still starting up. Please wait a few more minutes."
    echo "🌐 Check http://localhost:8080 and run: docker exec jenkins-snowbored cat /var/jenkins_home/secrets/initialAdminPassword"
fi
