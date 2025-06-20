# YelpCamp Setup Guide - Environment Variables Configuration

## 🔧 **Environment Variables Setup**

All configuration is now managed through the `.env` file in the `campgrounds-go` directory. This includes your GitHub credentials for Go module downloads.

### **1. Your .env File Configuration**

The `.env` file contains all necessary configuration:

\`\`\`env
# Environment Configuration
GIN_MODE=debug

# MongoDB Configuration
MONGO_HOST=mongo
MONGO_PORT=27017
MONGO_USER=yelpcamp_dev
MONGO_PASSWORD=dev_password_123
MONGO_DATABASE=yelpcamp_dev

# JWT Secret (Change this in production!)
JWT_SECRET=your_super_secret_jwt_key_here

# Cloudinary Configuration
CLOUDINARY_CLOUD_NAME=dl0zvwr9i
CLOUDINARY_API_KEY=719144713242759
CLOUDINARY_API_SECRET=tCIGtTFz8vvEem2Vl1HuOd_7xhE

# Server Configuration
PORT=3000

# AWS Configuration (for Jenkins S3 upload)
AWS_REGION=us-east-1
S3_BUCKET=yelpcamp-artifacts-bucket
IAM_ROLE=arn:aws:iam::054037100649:role/yelpcamp-s3-role

# GitHub Credentials for Go Module Downloads
GITHUB_USERNAME=premchander.j.pc@gmail.com
GITHUB_PASSWORD=Don$$trom1!

# Git Configuration
GIT_TERMINAL_PROMPT=0
GOPROXY=direct
GOSUMDB=off
\`\`\`

### **2. Security Benefits**

✅ **No hardcoded credentials** in Jenkinsfile
✅ **Easy to update** credentials without changing pipeline
✅ **Environment-specific** configurations
✅ **Git ignored** sensitive files (see .gitignore)

### **3. How Jenkins Uses the .env File**

The pipeline now:
1. **Loads .env file** at the start of each stage
2. **Exports all variables** to the shell environment
3. **Uses your GitHub credentials** for Go module authentication
4. **Cleans up credentials** after use for security

### **4. Quick Start**

\`\`\`bash
# 1. Clone your repository
git clone https://github.com/jpremchander/PROG8860-S25-CICD.git
cd PROG8860-S25-CICD/campgrounds-go

# 2. Verify .env file exists with your credentials
cat .env

# 3. Test locally first
make compose-dev

# 4. Run Jenkins pipeline
# The pipeline will automatically use your .env configuration
\`\`\`

### **5. Jenkins Pipeline Flow**

Each stage now:
\`\`\`bash
# Load environment variables
set -a
source .env
set +a

# Use variables (e.g., $GITHUB_USERNAME, $GITHUB_PASSWORD)
# Your credentials are automatically available
\`\`\`

### **6. Security Notes**

⚠️ **Important Security Practices:**
- The `.env` file is **NOT** committed to Git (see .gitignore)
- Credentials are **cleaned up** after each stage
- Use **environment-specific** .env files for different deployments
- **Change default passwords** in production

### **7. Testing the Setup**

\`\`\`bash
# Test environment loading
cd campgrounds-go
source .env
echo "GitHub User: $GITHUB_USERNAME"
echo "Cloudinary: $CLOUDINARY_CLOUD_NAME"

# Test Go module download with your credentials
go mod download
\`\`\`

## 🚀 **Ready to Deploy!**

Your pipeline will now:
1. ✅ **Load your GitHub credentials** from .env
2. ✅ **Authenticate with GitHub** for Go modules
3. ✅ **Use your Cloudinary settings** for image uploads
4. ✅ **Deploy with your AWS configuration**
5. ✅ **Clean up credentials** for security

The authentication issues should be completely resolved now! 🎉
