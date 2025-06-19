# Localhost Deployment Guide

Complete guide to set up Jenkins CI/CD pipeline with GitHub webhooks on localhost.

## 🚀 Quick Start

### 1. Prerequisites
\`\`\`bash
# Install required tools
sudo apt update
sudo apt install -y docker.io docker-compose curl wget unzip

# Start Docker service
sudo systemctl start docker
sudo systemctl enable docker

# Add user to docker group
sudo usermod -aG docker $USER
newgrp docker
\`\`\`

### 2. Setup Jenkins
\`\`\`bash
# Make scripts executable
chmod +x scripts/*.sh

# Run Jenkins setup
./scripts/setup-jenkins.sh
\`\`\`

### 3. Configure Jenkins
\`\`\`bash
# Wait for Jenkins to start, then run configuration
./scripts/configure-jenkins.sh
\`\`\`

### 4. Setup GitHub Webhook
\`\`\`bash
# Configure GitHub webhook with ngrok
./scripts/setup-github-webhook.sh
\`\`\`

### 5. Create Pipeline Job
1. Open Jenkins: http://localhost:8080
2. Go to **Manage Jenkins > Script Console**
3. Copy and paste the content from `jenkins/pipeline-job-config.groovy`
4. Update `YOUR_USERNAME` with your GitHub username
5. Click **Run**

## 📋 Step-by-Step Instructions

### Step 1: Initial Setup
\`\`\`bash
# Clone your repository
git clone https://github.com/YOUR_USERNAME/snowbored-game.git
cd snowbored-game

# Install dependencies
npm install

# Test local build
npm run build
npm run test
\`\`\`

### Step 2: Jenkins Installation
\`\`\`bash
# Run the setup script
./scripts/setup-jenkins.sh

# Wait for Jenkins to start (2-3 minutes)
# Access Jenkins at http://localhost:8080
\`\`\`

**Initial Jenkins Setup:**
1. Enter the initial admin password (displayed in terminal)
2. Install suggested plugins
3. Create admin user:
   - Username: `admin`
   - Password: `your-secure-password`
   - Full name: `Jenkins Admin`
   - Email: `your-email@example.com`

### Step 3: Plugin Configuration
\`\`\`bash
# Run configuration script
./scripts/configure-jenkins.sh
\`\`\`

**Manual Plugin Installation (if script fails):**
1. Go to **Manage Jenkins > Manage Plugins**
2. Install these plugins:
   - Git Plugin
   - GitHub Plugin
   - Pipeline Plugin
   - Docker Pipeline Plugin
   - NodeJS Plugin
   - HTML Publisher Plugin
   - Blue Ocean

### Step 4: Global Tool Configuration
1. Go to **Manage Jenkins > Global Tool Configuration**
2. **NodeJS Installations:**
   - Name: `NodeJS-18`
   - Install automatically: ✅
   - Version: `NodeJS 18.19.0`
3. **Docker Installations:**
   - Name: `Docker`
   - Installation root: `/usr/bin/docker`

### Step 5: GitHub Integration

#### Create Personal Access Token
1. Go to GitHub **Settings > Developer settings > Personal access tokens**
2. Generate new token with these scopes:
   - `repo` (Full control of private repositories)
   - `admin:repo_hook` (Full control of repository hooks)
3. Copy the token

#### Setup Webhook
\`\`\`bash
# Run webhook setup script
./scripts/setup-github-webhook.sh

# Enter your GitHub credentials when prompted
\`\`\`

**Manual Webhook Setup (if script fails):**
1. Go to your GitHub repository
2. **Settings > Webhooks > Add webhook**
3. **Payload URL:** `http://YOUR_NGROK_URL/github-webhook/`
4. **Content type:** `application/json`
5. **Events:** Push events, Pull requests
6. **Active:** ✅

### Step 6: Create Pipeline Job

#### Option 1: Script Console (Recommended)
1. **Manage Jenkins > Script Console**
2. Paste content from `jenkins/pipeline-job-config.groovy`
3. Update `YOUR_USERNAME` with your GitHub username
4. Click **Run**

#### Option 2: Manual Creation
1. **New Item > Pipeline**
2. **Name:** `SnowBored-Game-Pipeline`
3. **Pipeline:**
   - Definition: `Pipeline script from SCM`
   - SCM: `Git`
   - Repository URL: `https://github.com/YOUR_USERNAME/snowbored-game.git`
   - Branch: `*/main`
   - Script Path: `Jenkinsfile`
4. **Build Triggers:**
   - ✅ GitHub hook trigger for GITScm polling

### Step 7: Test the Pipeline

#### Trigger Build
\`\`\`bash
# Make a test commit
echo "# Test commit" >> README.md
git add README.md
git commit -m "test: trigger Jenkins pipeline"
git push origin main
\`\`\`

#### Monitor Build
1. Go to Jenkins dashboard
2. Click on `SnowBored-Game-Pipeline`
3. Watch the build progress
4. Check **Blue Ocean** for better visualization

## 🔧 Troubleshooting

### Common Issues

#### Jenkins Won't Start
\`\`\`bash
# Check Docker status
docker ps
docker logs jenkins-snowbored

# Restart Jenkins
docker-compose -f docker-compose.jenkins.yml restart
\`\`\`

#### Webhook Not Triggering
\`\`\`bash
# Check ngrok status
curl http://localhost:4040/api/tunnels

# Restart ngrok
pkill ngrok
ngrok http 8080 &

# Update webhook URL in GitHub
\`\`\`

#### Docker Permission Issues
\`\`\`bash
# Fix Docker permissions
sudo chmod 666 /var/run/docker.sock
sudo usermod -aG docker jenkins
docker restart jenkins-snowbored
\`\`\`

#### Build Failures
\`\`\`bash
# Check Jenkins logs
docker logs jenkins-snowbored

# Check build console output in Jenkins UI
# Go to Build > Console Output
\`\`\`

### Debug Commands
\`\`\`bash
# Check running containers
docker ps

# Check Jenkins logs
docker logs jenkins-snowbored

# Check application logs
docker logs snowbored-game

# Test webhook manually
curl -X POST http://localhost:8080/github-webhook/ \
  -H "Content-Type: application/json" \
  -d '{"ref":"refs/heads/main"}'
\`\`\`

## 🌐 Access URLs

- **Jenkins:** http://localhost:8080
- **Application:** http://localhost:3000
- **Ngrok Dashboard:** http://localhost:4040
- **Public Webhook:** Check `.jenkins-config` file

## 📊 Pipeline Stages

1. **Checkout:** Clone repository
2. **Build:** Install dependencies, lint, build
3. **Test:** Run unit tests with coverage
4. **Security Scan:** npm audit
5. **Docker Build:** Create container image
6. **Deploy:** Deploy to localhost
7. **Archive:** Store build artifacts

## 🎯 Success Indicators

✅ **Jenkins accessible at localhost:8080**
✅ **Pipeline job created successfully**
✅ **GitHub webhook configured**
✅ **Build triggers on git push**
✅ **All pipeline stages pass**
✅ **Application deployed at localhost:3000**
✅ **Docker container running**

## 🚀 Next Steps

1. **Test the complete flow:**
   - Make code changes
   - Push to GitHub
   - Verify automatic build trigger
   - Check deployment

2. **Monitor and optimize:**
   - Review build times
   - Check test coverage
   - Monitor resource usage

3. **Enhance pipeline:**
   - Add more test stages
   - Implement notifications
   - Add deployment environments
