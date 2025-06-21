# YelpCamp Go - CI/CD Pipeline Project

A full-stack campground review application built with Go, featuring a complete CI/CD pipeline using Jenkins, MongoDB integration, and AWS S3 for image storage.

## 🏗️ Architecture

- **Backend**: Go with Gin framework
- **Database**: MongoDB
- **File Storage**: AWS S3
- **CI/CD**: Jenkins Pipeline
- **Deployment**: Multi-environment (Dev/Prod)

## 📋 Prerequisites

- Go 1.21 or higher
- MongoDB (local or remote)
- AWS Account with S3 access
- Jenkins with required plugins
- Git

## 🚀 Quick Start

### Local Development

1. **Clone the repository**
   \`\`\`bash
   git clone <your-repo-url>
   cd campgrounds-go
   \`\`\`

2. **Install dependencies**
   \`\`\`bash
   go mod tidy
   \`\`\`

3. **Set up environment variables**
   \`\`\`bash
   cp .env.dev .env
   # Edit .env with your configuration
   \`\`\`

4. **Start MongoDB**
   \`\`\`bash
   # Using Docker
   docker run -d -p 27017:27017 --name mongodb mongo:latest
   
   # Or start your local MongoDB service
   sudo systemctl start mongod
   \`\`\`

5. **Run the application**
   \`\`\`bash
   go run main.go
   \`\`\`

6. **Access the application**
   - Development: http://localhost:3001
   - Health Check: http://localhost:3001/health

## 🧪 Testing

### Run Unit Tests
\`\`\`bash
# Run all tests
go test ./...

# Run tests with coverage
go test -v -race -coverprofile=coverage.out ./...

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html
\`\`\`

### Test Coverage
The project includes comprehensive unit tests covering:
- HTTP endpoints and routing
- Model validation
- Database operations
- S3 integration utilities

## 🔧 CI/CD Pipeline

### Jenkins Pipeline Overview

The Jenkins pipeline (`Jenkinsfile`) includes the following stages:

1. **Checkout** - Source code retrieval
2. **Setup Go Environment** - Go installation and configuration
3. **Build** - Application compilation
4. **Lint & Static Analysis** - Code quality checks
5. **Test** - Unit test execution with coverage
6. **Build Artifacts** - Environment-specific builds
7. **Upload to S3** - Artifact storage
8. **Deploy to Dev** - Automatic deployment to development
9. **Deploy to Production** - Manual deployment to production

### Pipeline Features

- ✅ **Multi-environment support** (dev/prod)
- ✅ **Automated testing** with coverage reports
- ✅ **Code quality checks** (golangci-lint, go vet, go fmt)
- ✅ **Artifact management** via AWS S3
- ✅ **Health checks** post-deployment
- ✅ **Manual approval** for production deployments

### Triggering the Pipeline

#### Automatic Triggers
- **Development**: Push to `develop` branch
- **Production**: Push to `main` branch

#### Manual Triggers
1. Go to Jenkins dashboard
2. Select the project
3. Click "Build with Parameters"
4. Choose deployment environment (dev/prod)
5. Optionally skip tests
6. Click "Build"

### Environment Configuration

#### Development Environment
- **Port**: 3001
- **Mode**: Debug
- **Database**: campgrounds_dev
- **Auto-deploy**: On develop branch push

#### Production Environment
- **Port**: 3000
- **Mode**: Release
- **Database**: campgrounds_prod
- **Deploy**: Manual approval required

## 🏗️ Jenkins Setup

### Required Jenkins Plugins
- Pipeline
- Git
- AWS Steps
- HTML Publisher
- Build Timeout

### Jenkins Credentials Setup

1. **AWS Credentials**
   - ID: `AWSCREDENTIALS`
   - Type: AWS Credentials
   - Access Key ID: Your AWS Access Key
   - Secret Access Key: Your AWS Secret Key
   - Region: ca-central-1

2. **Environment Variables**
   Configure in Jenkins Global Properties:
   - `AWS_DEFAULT_REGION=ca-central-1`
   - `AWS_S3_BUCKET=yelp-camp-project`

### Pipeline Configuration

1. Create new Pipeline job in Jenkins
2. Configure SCM to point to your repository
3. Set Pipeline script path to `Jenkinsfile`
4. Configure branch sources (main, develop)
5. Save and run initial build

## 📦 Deployment

### Artifact Structure
\`\`\`
artifacts/
├── campgrounds-go-{env}     # Compiled binary
├── views/                   # HTML templates
├── public/                  # Static assets
├── .env.{env}              # Environment config
└── deploy.sh               # Deployment script
\`\`\`

### Manual Deployment

1. **Download artifacts from S3**
   \`\`\`bash
   aws s3 cp s3://yelp-camp-project/artifacts/prod/campgrounds-go-latest.tar.gz .
   \`\`\`

2. **Extract and deploy**
   \`\`\`bash
   tar -xzf campgrounds-go-latest.tar.gz
   ./deploy.sh prod
   \`\`\`

### Health Checks

Both environments include health check endpoints:
- Development: http://localhost:3001/health
- Production: http://localhost:3000/health

Expected response:
\`\`\`json
{
  "status": "healthy",
  "service": "campgrounds-go"
}
\`\`\`

## 🗄️ Database Schema

### Collections

#### Campgrounds
\`\`\`go
type Campground struct {
    ID          primitive.ObjectID
    Title       string
    Price       float64
    Description string
    Location    string
    Images      []Image
    Author      primitive.ObjectID
    Reviews     []primitive.ObjectID
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
\`\`\`

#### Users
\`\`\`go
type User struct {
    ID        primitive.ObjectID
    Username  string
    Email     string
    Password  string
    CreatedAt time.Time
    UpdatedAt time.Time
}
\`\`\`

#### Reviews
\`\`\`go
type Review struct {
    ID         primitive.ObjectID
    Body       string
    Rating     int
    Author     primitive.ObjectID
    Campground primitive.ObjectID
    CreatedAt  time.Time
    UpdatedAt  time.Time
}
\`\`\`

## 🔐 AWS S3 Configuration

### Bucket Structure
\`\`\`
yelp-camp-project/
├── campgrounds/           # User uploaded images
│   ├── 1640995200_image1.jpg
│   └── 1640995300_image2.png
└── artifacts/            # CI/CD artifacts
    ├── dev/
    │   ├── campgrounds-go-123.tar.gz
    │   └── campgrounds-go-latest.tar.gz
    └── prod/
        ├── campgrounds-go-124.tar.gz
        └── campgrounds-go-latest.tar.gz
\`\`\`

### Required S3 Permissions
\`\`\`json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "s3:GetObject",
                "s3:PutObject",
                "s3:DeleteObject"
            ],
            "Resource": "arn:aws:s3:::yelp-camp-project/*"
        },
        {
            "Effect": "Allow",
            "Action": [
                "s3:ListBucket"
            ],
            "Resource": "arn:aws:s3:::yelp-camp-project"
        }
    ]
}
\`\`\`

## 🔍 Monitoring & Logging

### Application Logs
- Development: `campgrounds-dev.log`
- Production: `campgrounds-prod.log`

### Jenkins Pipeline Logs
- Build logs available in Jenkins console
- Test coverage reports archived as artifacts
- S3 upload confirmations logged

## 🚨 Troubleshooting

### Common Issues

1. **Build Failures**
   - Check Go version compatibility
   - Verify all dependencies are available
   - Review build logs for specific errors

2. **Test Failures**
   - Ensure MongoDB is running
   - Check test database connectivity
   - Review test logs for specific failures

3. **Deployment Issues**
   - Verify port availability
   - Check environment variables
   - Confirm AWS credentials are valid

4. **S3 Upload Failures**
   - Verify AWS credentials in Jenkins
   - Check S3 bucket permissions
   - Confirm bucket exists and is accessible

### Debug Commands
\`\`\`bash
# Check application status
curl http://localhost:3000/health

# View application logs
tail -f campgrounds-prod.log

# Check running processes
ps aux | grep campgrounds-go

# Test S3 connectivity
aws s3 ls s3://yelp-camp-project/
\`\`\`

## 📈 Performance Considerations

- **Database Indexing**: Ensure proper indexes on frequently queried fields
- **Image Optimization**: Consider image compression before S3 upload
- **Caching**: Implement Redis for session management in production
- **Load Balancing**: Use nginx for production deployments

## 🔒 Security Best Practices

- Environment variables for sensitive data
- HTTPS in production (configure reverse proxy)
- Input validation and sanitization
- Rate limiting for API endpoints
- Regular security updates

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 📞 Support

For issues and questions:
- Create an issue in the repository
- Check existing documentation
- Review Jenkins build logs for CI/CD issues

---

**Note**: This project is designed for educational purposes as part of the PROG8860 CI/CD course midterm assignment.
\`\`\`

Finally, let's create a .gitignore file:
