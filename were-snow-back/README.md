# SnowBored Game - CI/CD Pipeline Project

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
