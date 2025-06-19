# YelpCamp Setup Guide - Ready to Deploy!

## ✅ Your Cloudinary Configuration
From your CLOUDINARY_URL, I've extracted:
- **Cloud Name**: `dl0zvwr9i`
- **API Key**: `719144713242759`
- **API Secret**: `tCIGtTFz8vvEem2Vl1HuOd_7xhE`

## 🚀 Quick Start (You're Ready!)

### 1. Clone and Setup
\`\`\`bash
git clone <your-repo>
cd yelpcamp-go-jenkins

# Create your .env file
cp .env.example .env
\`\`\`

### 2. Your .env File Should Look Like:
\`\`\`env
# Environment Configuration
GIN_MODE=debug

# MongoDB Configuration
MONGO_HOST=mongo
MONGO_PORT=27017
MONGO_USER=yelpcamp_dev
MONGO_PASSWORD=dev_password_123
MONGO_DATABASE=yelpcamp_dev

# JWT Secret
JWT_SECRET=your_super_secret_jwt_key_here

# Cloudinary Configuration (Your actual credentials)
CLOUDINARY_CLOUD_NAME=dl0zvwr9i
CLOUDINARY_API_KEY=719144713242759
CLOUDINARY_API_SECRET=tCIGtTFz8vvEem2Vl1HuOd_7xhE

# Server Configuration
PORT=8080

# AWS Configuration (Already configured)
AWS_REGION=us-east-1
S3_BUCKET=yelpcamp-artifacts-bucket
IAM_ROLE=arn:aws:iam::054037100649:role/yelpcamp-s3-role
\`\`\`

### 3. Test Locally First
\`\`\`bash
# Install dependencies
make deps

# Run tests
make test

# Start development environment
make compose-dev
\`\`\`

Your app will be available at:
- **Application**: http://localhost:8080
- **MongoDB Admin**: http://localhost:8081 (admin/admin123)

### 4. Jenkins Setup
\`\`\`bash
# Start Jenkins (if not already running)
docker run -d -p 8080:8080 -p 50000:50000 \
  -v jenkins_home:/var/jenkins_home \
  -v /var/run/docker.sock:/var/run/docker.sock \
  --name jenkins jenkins/jenkins:lts

# Access Jenkins at http://localhost:8080
\`\`\`

### 5. Jenkins Configuration
1. **Install Plugins**:
   - Pipeline
   - Docker Pipeline
   - AWS Steps
   - S3 Publisher

2. **Add Credentials**:
   - AWS credentials for S3 access
   - Your Cloudinary credentials (already in code)

3. **Create Pipeline Job**:
   - New Item → Pipeline
   - SCM: Git (your repository)
   - Script Path: `Jenkinsfile`

### 6. Run Your First Build
1. Go to your Jenkins job
2. Click "Build with Parameters"
3. Choose environment: `dev`
4. Click "Build"

## 🎯 What Happens During Build:

1. **Checkout** - Gets your code
2. **Environment Setup** - Creates config with your Cloudinary credentials
3. **Lint & Static Analysis** - Checks code quality
4. **Test** - Runs unit tests with coverage
5. **Build Application** - Compiles Go binary
6. **Build Docker Images** - Creates containerized app
7. **Upload Artifacts to S3** - Stores build artifacts
8. **Deploy to Localhost** - Runs with Docker Compose
9. **Post-Deployment Tests** - Verifies everything works

## 📸 Image Upload Testing

Once deployed, you can test image uploads:

\`\`\`bash
# Register a user
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123"}'

# Login to get token
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'

# Create a campground
curl -X POST http://localhost:8080/api/campgrounds \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"title":"Test Camp","description":"A beautiful test campground","location":"Test Valley","price":25.99}'

# Upload images to campground
curl -X POST http://localhost:8080/api/campgrounds/CAMPGROUND_ID/images \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "images=@/path/to/your/image.jpg"
\`\`\`

## 🔍 Monitoring Your Deployment

\`\`\`bash
# Check application logs
make logs

# Check MongoDB logs
make logs-mongo

# Check container status
make status

# Access MongoDB shell
make mongo-shell
\`\`\`

## 🎉 You're All Set!

Your YelpCamp application is now ready with:
- ✅ **Cloudinary integration** configured with your credentials
- ✅ **MongoDB database** with sample data
- ✅ **Jenkins CI/CD pipeline** ready to deploy
- ✅ **S3 artifact storage** configured
- ✅ **Multi-environment support** (dev/prod)
- ✅ **Comprehensive testing** and monitoring

## 🚨 Important Notes:

1. **Change JWT_SECRET** in production
2. **Change MongoDB passwords** for production
3. **Your Cloudinary credentials** are already configured
4. **S3 bucket and IAM role** are already set up
5. **Jenkins will handle** all deployments automatically

## 📞 Need Help?

If you encounter any issues:
1. Check the logs: `make logs`
2. Verify containers: `make status`
3. Test locally first: `make compose-dev`
4. Check Jenkins build logs for detailed error messages

You're ready to deploy! 🚀
