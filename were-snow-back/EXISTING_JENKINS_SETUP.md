# Configure Existing Jenkins for SnowBored Game

Quick setup guide for configuring your existing Jenkins instance.

## 🚀 Quick Setup (10 minutes)

### Step 1: Configure Jenkins
\`\`\`bash
# Make scripts executable
chmod +x scripts/*.sh

# Configure existing Jenkins
./scripts/configure-existing-jenkins.sh
\`\`\`

### Step 2: Setup GitHub Webhook
\`\`\`bash
# Setup webhook (will use ngrok if needed)
./scripts/setup-webhook-existing.sh
\`\`\`

### Step 3: Create Pipeline Job
\`\`\`bash
# Create the pipeline job
./scripts/create-jenkins-job.sh
\`\`\`

### Step 4: Test Everything
\`\`\`bash
# Test the complete setup
./scripts/test-pipeline.sh
\`\`\`

## 📋 Manual Configuration (Alternative)

### 1. Install Required Plugins
Go to **Manage Jenkins > Manage Plugins > Available**:
- Git Plugin
- GitHub Plugin  
- Pipeline Plugin
- Docker Pipeline Plugin
- NodeJS Plugin
- HTML Publisher Plugin
- Blue Ocean (optional)

### 2. Configure Global Tools
**Manage Jenkins > Global Tool Configuration:**

**NodeJS:**
- Name: `NodeJS-18`
- Install automatically: ✅
- Version: `NodeJS 18.19.0`

**Docker:**
- Name: `Docker`
- Installation root: `/usr/bin/docker`

### 3. Create Pipeline Job
1. **New Item > Pipeline**
2. **Name:** `SnowBored-Game-Pipeline`
3. **Build Triggers:** ✅ GitHub hook trigger for GITScm polling
4. **Pipeline:**
   - Definition: `Pipeline script from SCM`
   - SCM: `Git`
   - Repository URL: `https://github.com/YOUR_USERNAME/snowbored-game.git`
   - Branch: `*/main`
   - Script Path: `Jenkinsfile`

### 4. Configure GitHub Webhook
**In your GitHub repository:**
1. **Settings > Webhooks > Add webhook**
2. **Payload URL:** `http://YOUR_JENKINS_URL/github-webhook/`
3. **Content type:** `application/json`
4. **Events:** Push events, Pull requests
5. **Active:** ✅

## 🔧 Troubleshooting

### Plugin Installation Issues
\`\`\`bash
# Restart Jenkins after plugin installation
sudo systemctl restart jenkins
# or via UI: Manage Jenkins > Restart
\`\`\`

### Webhook Not Working
\`\`\`bash
# Check if Jenkins is accessible from internet
curl -I http://YOUR_PUBLIC_IP:8080

# If not accessible, use ngrok
ngrok http 8080
# Use the ngrok URL for webhook
\`\`\`

### Permission Issues
\`\`\`bash
# Add Jenkins user to docker group
sudo usermod -aG docker jenkins
sudo systemctl restart jenkins
\`\`\`

## ✅ Success Checklist

- [ ] Jenkins accessible at your URL
- [ ] Required plugins installed
- [ ] Global tools configured
- [ ] Pipeline job created
- [ ] GitHub webhook configured
- [ ] Test build successful
- [ ] Application deploys to localhost:3000

## 🎯 Expected Result

After setup:
- ✅ **Automatic builds** on GitHub push
- ✅ **Complete CI/CD pipeline** with all stages
- ✅ **Application deployed** at http://localhost:3000
- ✅ **Build reports** and artifacts
- ✅ **Professional Jenkins setup**

## 📞 Need Help?

If you encounter issues:
1. Check Jenkins logs: `sudo journalctl -u jenkins`
2. Verify Docker permissions
3. Test webhook manually
4. Check firewall settings
\`\`\`
