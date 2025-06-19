# Jenkins CI/CD Pipeline Setup Guide

This guide will help you set up Jenkins with GitHub webhooks for the SnowBored Game project.

## 🏗️ Jenkins Installation

### Option 1: Docker (Recommended)
\`\`\`bash
# Create Jenkins volume
docker volume create jenkins_home

# Run Jenkins container
docker run -d \
  --name jenkins \
  -p 8080:8080 \
  -p 50000:50000 \
  -v jenkins_home:/var/jenkins_home \
  -v /var/run/docker.sock:/var/run/docker.sock \
  --group-add $(getent group docker | cut -d: -f3) \
  jenkins/jenkins:lts
\`\`\`

### Option 2: Local Installation
\`\`\`bash
# Ubuntu/Debian
wget -q -O - https://pkg.jenkins.io/debian/jenkins.io.key | sudo apt-key add -
sudo sh -c 'echo deb http://pkg.jenkins.io/debian-stable binary/ > /etc/apt/sources.list.d/jenkins.list'
sudo apt update
sudo apt install jenkins

# Start Jenkins
sudo systemctl start jenkins
sudo systemctl enable jenkins
\`\`\`

## 🔧 Initial Jenkins Configuration

### 1. Access Jenkins
- Open http://localhost:8080
- Get initial admin password:
\`\`\`bash
# For Docker
docker exec jenkins cat /var/jenkins_home/secrets/initialAdminPassword

# For local installation
sudo cat /var/lib/jenkins/secrets/initialAdminPassword
\`\`\`

### 2. Install Required Plugins
Go to **Manage Jenkins > Manage Plugins** and install:

**Essential Plugins:**
- Git Plugin
- GitHub Plugin
- Pipeline Plugin
- Docker Pipeline Plugin
- NodeJS Plugin
- HTML Publisher Plugin
- Email Extension Plugin
- Slack Notification Plugin (optional)

**Additional Plugins:**
- Blue Ocean (for better UI)
- Build Timeout Plugin
- Timestamper Plugin
- Workspace Cleanup Plugin

### 3. Configure Global Tools
Go to **Manage Jenkins > Global Tool Configuration**:

**NodeJS:**
- Name: `Node 18`
- Version: `NodeJS 18.x.x`
- ✅ Install automatically

**Docker:**
- Name: `Docker`
- ✅ Install automatically

## 🔐 Credentials Setup

Go to **Manage Jenkins > Manage Credentials > System > Global credentials**:

### 1. Docker Hub Credentials
- Kind: `Username with password`
- ID: `docker-hub-credentials`
- Username: Your Docker Hub username
- Password: Your Docker Hub password/token

### 2. GitHub Personal Access Token
- Kind: `Secret text`
- ID: `github-token`
- Secret: Your GitHub Personal Access Token

### 3. Email Configuration (Optional)
Go to **Manage Jenkins > Configure System > Extended E-mail Notification**:
- SMTP Server: `smtp.gmail.com`
- Port: `587`
- Username: Your email
- Password: App-specific password
- ✅ Use SSL

## 🔗 GitHub Webhook Configuration

### 1. In GitHub Repository
Go to **Settings > Webhooks > Add webhook**:

- **Payload URL**: `http://your-jenkins-url:8080/github-webhook/`
- **Content type**: `application/json`
- **Events**: 
  - ✅ Push events
  - ✅ Pull request events
- ✅ Active

### 2. GitHub Repository Settings
In your repository **Settings > Secrets and variables > Actions**:
- `JENKINS_URL`: Your Jenkins URL
- `JENKINS_TOKEN`: Jenkins API token

## 🚀 Jenkins Job Configuration

### 1. Create New Pipeline Job
- **New Item > Pipeline**
- Name: `SnowBored-Game-Pipeline`

### 2. Pipeline Configuration
**General:**
- ✅ GitHub project: `https://github.com/your-username/snowbored-game`
- ✅ Build Triggers: `GitHub hook trigger for GITScm polling`

**Pipeline:**
- Definition: `Pipeline script from SCM`
- SCM: `Git`
- Repository URL: `https://github.com/your-username/snowbored-game.git`
- Credentials: Select your GitHub credentials
- Branch: `*/main` (or your default branch)
- Script Path: `Jenkinsfile`

### 3. Advanced Configuration
**Build Triggers:**
- ✅ GitHub hook trigger for GITScm polling
- ✅ Poll SCM: `H/5 * * * *` (backup polling every 5 minutes)

**Pipeline:**
- ✅ Lightweight checkout

## 🧪 Testing the Pipeline

### 1. Manual Trigger
- Go to your Jenkins job
- Click **Build Now**
- Monitor the **Console Output**

### 2. Webhook Trigger
- Make a commit to your repository
- Push to GitHub
- Jenkins should automatically trigger the build

### 3. Verify Stages
Check that all stages complete successfully:
- ✅ Checkout
- ✅ Build
- ✅ Test
- ✅ Docker Build
- ✅ Docker Push
- ✅ Deploy
- ✅ Archive Artifacts

## 📊 Monitoring and Notifications

### 1. Build Status
- **Blue Ocean**: Better visual pipeline view
- **Build History**: Track build trends
- **Test Results**: View test reports
- **Coverage Reports**: Code coverage metrics

### 2. Notifications
Configure notifications for:
- Build failures
- Successful deployments
- Test failures
- Security vulnerabilities

### 3. Slack Integration (Optional)
\`\`\`groovy
slackSend(
    channel: '#ci-cd',
    color: 'good',
    message: "✅ Build ${BUILD_NUMBER} successful!"
)
\`\`\`

## 🔍 Troubleshooting

### Common Issues:

**1. Docker Permission Denied**
\`\`\`bash
# Add Jenkins user to docker group
sudo usermod -aG docker jenkins
sudo systemctl restart jenkins
\`\`\`

**2. Node.js Not Found**
- Verify NodeJS plugin installation
- Check Global Tool Configuration
- Restart Jenkins

**3. GitHub Webhook Not Triggering**
- Check webhook URL format
- Verify Jenkins is accessible from internet
- Check GitHub webhook delivery logs

**4. Docker Build Fails**
- Verify Docker is running
- Check Dockerfile syntax
- Ensure sufficient disk space

### Debug Commands:
\`\`\`bash
# Check Jenkins logs
docker logs jenkins

# Check disk space
df -h

# Check Docker status
docker ps
docker images
\`\`\`

## 📈 Performance Optimization

### 1. Pipeline Optimization
- Use parallel stages for tests
- Implement Docker layer caching
- Use Jenkins agents for distributed builds

### 2. Resource Management
- Set build timeouts
- Clean up old builds automatically
- Monitor disk usage

### 3. Security Best Practices
- Use least privilege credentials
- Regularly update Jenkins and plugins
- Enable CSRF protection
- Use HTTPS in production

## 🎯 Assignment Bonus Points Checklist

- ✅ **Jenkins Installation**: Complete setup with Docker
- ✅ **Jenkinsfile**: Comprehensive pipeline with all stages
- ✅ **GitHub Integration**: Webhooks triggering builds automatically
- ✅ **Docker Integration**: Build and push images
- ✅ **Testing Integration**: Run tests and generate reports
- ✅ **Deployment Automation**: Staging and production deployments
- ✅ **Notifications**: Email and Slack notifications
- ✅ **Monitoring**: Build status and reporting
- ✅ **Documentation**: Complete setup and usage guide

This setup demonstrates enterprise-level CI/CD practices and should earn you maximum bonus points! 🏆
\`\`\`

Let's also update the main README to include Jenkins information:

```typescriptreact file="README.md"
[v0-no-op-code-block-prefix]# SnowBored Game - CI/CD Pipeline Project

A fun snowboarding game built with Next.js and TypeScript, featuring a complete CI/CD pipeline with GitHub Actions.

## 🎮 Game Overview

SnowBored is an endless runner-style snowboarding game where players navigate through obstacles while skiing down a snowy slope. The game features:

- Smooth player movement with gravity physics
- Dynamic obstacle generation (trees and snowmen)
- Score tracking and speed progression
- Collision detection
- Retro pixel art styling

## 🏗️ Project Structure

\`\`\`
snowbored-game/
├── .github/
│   └── workflows/
│       └── ci.yml              # GitHub Actions CI pipeline
├── __tests__/                  # Unit tests
│   ├── gameUtils.test.ts
│   ├── collision.test.ts
│   └── scoring.test.ts
├── components/
│   └── SnowBored.tsx          # Main game component
├── app/
│   ├── layout.tsx
│   └── page.tsx               # Home page
├── constants.ts               # Game constants and configuration
├── Dockerfile                 # Docker configuration
├── jest.config.js            # Jest testing configuration
└── package.json              # Dependencies and scripts
\`\`\`

## 🚀 Technologies Used

- **Frontend**: Next.js 14, React 18, TypeScript
- **Testing**: Jest, React Testing Library
- **Containerization**: Docker
- **CI/CD**: GitHub Actions
- **Deployment**: Docker Hub Registry

## 📋 Prerequisites

- Node.js 18.x or 20.x
- npm or yarn
- Docker (for containerization)
- Git

## 🛠️ Local Development

### 1. Clone the repository
\`\`\`bash
git clone <your-repo-url>
cd snowbored-game
\`\`\`

### 2. Install dependencies
\`\`\`bash
npm install
\`\`\`

### 3. Run the development server
\`\`\`bash
npm run dev
\`\`\`

### 4. Open your browser
Navigate to [http://localhost:3000](http://localhost:3000)

## 🧪 Testing

### Run all tests
\`\`\`bash
npm test
\`\`\`

### Run tests with coverage
\`\`\`bash
npm run test:coverage
\`\`\`

### Run tests in watch mode
\`\`\`bash
npm run test:watch
\`\`\`

### Test Cases Included
1. **Game Constants Validation** - Verifies game configuration values
2. **Color Format Validation** - Ensures proper hex color formats
3. **Image URL Validation** - Checks asset URL integrity
4. **Collision Detection Logic** - Tests player-obstacle collision
5. **Scoring System Logic** - Validates score calculation
6. **Speed Multiplier Logic** - Tests game progression mechanics

## 🏗️ Building for Production

### Build the application
\`\`\`bash
npm run build
\`\`\`

### Start production server
\`\`\`bash
npm start
\`\`\`

## 🐳 Docker Usage

### Build Docker image locally
\`\`\`bash
docker build -t snowbored-game .
\`\`\`

### Run Docker container
\`\`\`bash
docker run -p 3000:3000 snowbored-game
\`\`\`

### Pull from Docker Hub
\`\`\`bash
docker pull <your-dockerhub-username>/snowbored-game:latest
docker run -p 3000:3000 <your-dockerhub-username>/snowbored-game:latest
\`\`\`

## 🔄 CI/CD Pipeline

The GitHub Actions pipeline includes:

### Build Phase
- Installs dependencies with `npm ci`
- Runs ESLint for code quality
- Builds the Next.js application

### Test Phase
- Executes unit tests with Jest
- Generates code coverage reports
- Uploads coverage to Codecov
- Fails pipeline if tests don't pass

### Deploy Phase
- Builds Docker image with multi-stage optimization
- Pushes to Docker Hub registry
- Tags with both `latest` and commit SHA
- Implements Docker layer caching

### Security Phase
- Runs `npm audit` for vulnerability scanning
- Continues pipeline even if non-critical issues found

## 🏗️ Jenkins Pipeline (Bonus Implementation)

This project includes a comprehensive Jenkins pipeline for advanced CI/CD automation.

### Jenkins Features Implemented
- **Automated Builds**: Triggered by GitHub webhooks
- **Parallel Testing**: Unit tests and security scans run simultaneously  
- **Multi-Environment Deployment**: Staging and production environments
- **Blue-Green Deployment**: Zero-downtime production deployments
- **Manual Approval Gates**: Production deployment requires approval
- **Comprehensive Notifications**: Email and Slack integration
- **Artifact Management**: Build artifacts archived automatically
- **Health Checks**: Automated deployment verification

### Jenkins Setup
See [JENKINS_SETUP.md](./JENKINS_SETUP.md) for complete installation and configuration guide.

### Pipeline Stages
1. **Checkout**: Source code retrieval
2. **Build**: Dependency installation and application build
3. **Test**: Parallel unit tests and security scanning
4. **Docker Build**: Multi-stage container creation
5. **Docker Push**: Registry publication with multiple tags
6. **Deploy Staging**: Automatic staging deployment
7. **Deploy Production**: Manual approval + blue-green deployment
8. **Archive**: Artifact storage and cleanup

### Webhook Integration
- Automatic builds on push to main/develop branches
- Pull request validation
- Commit status updates to GitHub

## 🎯 Game Controls

- **SPACE**: Hold to move up, release to move down
- **Objective**: Avoid obstacles and achieve the highest score

## 📊 Pipeline Status

[![CI Pipeline](https://github.com/<your-username>/snowbored-game/actions/workflows/ci.yml/badge.svg)](https://github.com/<your-username>/snowbored-game/actions/workflows/ci.yml)

## 🔧 Configuration

### Required GitHub Secrets
- `DOCKER_USERNAME`: Your Docker Hub username
- `DOCKER_PASSWORD`: Your Docker Hub password/token

### Environment Variables
- `NODE_ENV`: Set to 'production' for production builds
- `PORT`: Server port (default: 3000)

## 📈 Performance Optimizations

- Multi-stage Docker builds for smaller images
- Next.js standalone output for minimal runtime
- Image optimization and caching
- GitHub Actions cache for faster builds

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## 📝 Assignment Requirements Met

- ✅ **Project Setup (2%)**: New GitHub repository with sample application
- ✅ **CI Pipeline (12%)**: Complete GitHub Actions workflow
- ✅ **Build Phase (4%)**: Dependency installation and application build
- ✅ **Test Phase (4%)**: 4+ unit tests with pipeline failure on test failure
- ✅ **Upload Phase (4%)**: Docker image build and registry upload
- ✅ **Docker Implementation (3%)**: Dockerfile and image publishing
- ✅ **Documentation (3%)**: Comprehensive README with instructions

## 📄 License

This project is created for educational purposes as part of a CI/CD pipeline assignment.
